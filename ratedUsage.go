package exporter

import (
	"net/http"
	"time"

	"github.com/phoenixnap/go-sdk-bmc/billingapi"
	"github.com/phoenixnap/go-sdk-bmc/bmcapi"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type RatedUsageTagInfo struct {
	Key   string
	Value string
}

type RatedUsageStats struct {
	Cost            float64
	Hostname        string
	Location        string
	PriceModel      string
	ProductCategory string
	ProductCode     string
	YearMonth       string
	BillingTags     []RatedUsageTagInfo
}

var priorMonthsRatedUsageLoaded bool = false
var servers []bmcapi.Server

func GetRatedUsageStats(config BmcApiConfiguration) ([]RatedUsageStats, error) {
	apiClient := getBillingApiClient(config.ToClientCredentials())
	bmcClient := getBmcApiClient(config.ToClientCredentials())

	var stats []RatedUsageStats
	var r *http.Response
	var err error
	servers, r, err = bmcClient.ServersApi.ServersGet(getContext()).Execute()
	if err != nil {
		log.WithField("HttpResponse", r).WithError(err).Error("Error when calling `ServersApi.ServersGet`.")
		return stats, err
	}

	log.WithField("priorMonthsRatedUsageLoaded", priorMonthsRatedUsageLoaded).Trace("Prior Month Load State")

	if !priorMonthsRatedUsageLoaded && config.HistoricalRatedUsage.Enable {
		toDate := time.Now().Format("2006-01")
		fromDate := time.Now().AddDate(0, -config.HistoricalRatedUsage.NumberOfPriorMonths, 0).Format("2006-01")
		log.WithField("ToYearMonth", toDate).WithField("FromYearMonth", fromDate).WithField("ProductCategory", billingapi.SERVER).Trace("Call RatedUsageGet")

		resp, r, err := apiClient.RatedUsageApi.RatedUsageGet(getContext()).FromYearMonth(fromDate).ToYearMonth(toDate).ProductCategory(billingapi.SERVER).Execute()
		if err != nil {
			log.WithField("HttpResponse", r).WithError(err).Error("Error when calling `RatedUsageApi.RatedUsageGet`.")
			return stats, err
		}

		for _, ratedUsage := range resp {
			if ratedUsageStat := ConvertRatedUsageServerRecordToStats(ratedUsage.ServerRecord); ratedUsageStat != nil {
				stats = append(stats, *ratedUsageStat)
			}
		}

		priorMonthsRatedUsageLoaded = true
	} else {
		resp, r, err := apiClient.RatedUsageApi.RatedUsageMonthToDateGet(getContext()).ProductCategory(billingapi.SERVER).Execute()
		if err != nil {
			log.WithField("HttpResponse", r).WithError(err).Error("Error when calling `RatedUsageApi.RatedUsageMonthToDateGet`.")
			return stats, err
		}

		for _, ratedUsage := range resp {
			if ratedUsageStat := ConvertRatedUsageServerRecordToStats(ratedUsage.ServerRecord); ratedUsageStat != nil {
				stats = append(stats, *ratedUsageStat)
			}
		}
	}

	return stats, nil
}

func ConvertRatedUsageServerRecordToStats(record *billingapi.ServerRecord) *RatedUsageStats {
	var ratedUsageStat RatedUsageStats

	if record.GetActive() {
		ratedUsageStat.Cost = float64(record.GetCost()) / float64(100)
		ratedUsageStat.Hostname = record.GetMetadata().Hostname
		ratedUsageStat.Location = string(record.GetLocation())
		ratedUsageStat.PriceModel = record.GetPriceModel()
		ratedUsageStat.ProductCategory = record.GetProductCategory()
		ratedUsageStat.ProductCode = record.GetProductCode()
		ratedUsageStat.YearMonth = record.GetYearMonth()
		ratedUsageStat.BillingTags = GetBillingTagsFromServers(record.GetId())

		return &ratedUsageStat
	}

	return nil
}

func GetBillingTagsFromServers(serverId string) []RatedUsageTagInfo {
	tags := make([]RatedUsageTagInfo, 0)

	for _, s := range servers {
		if s.GetId() == serverId {
			for _, t := range s.GetTags() {
				if t.GetIsBillingTag() {
					tags = append(tags, RatedUsageTagInfo{
						Key:   t.GetName(),
						Value: t.GetValue(),
					})
					return tags
				}
			}
		}
	}

	return tags
}

func (s *RatedUsageStats) ToPrometheusMetric(desc *prometheus.Desc, billingLabels []string) prometheus.Metric {
	labelValues := []string{
		s.ProductCode,
		s.Location,
		s.Hostname,
		s.PriceModel,
		s.ProductCategory,
		s.YearMonth,
	}

	log.WithField("s", s).Trace("ToPrometheusMetrix")
	for _, label := range billingLabels {
		var value string = "untagged"
		for _, t := range s.BillingTags {
			if t.Key == label {
				value = t.Value
				break
			}
		}
		labelValues = append(labelValues, value)
	}

	return prometheus.MustNewConstMetric(
		desc,
		prometheus.CounterValue,
		s.Cost,
		labelValues...,
	)
}
