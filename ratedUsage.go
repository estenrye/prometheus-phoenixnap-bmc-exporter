package exporter

import (
	"time"

	"github.com/phoenixnap/go-sdk-bmc/billingapi"
	"github.com/phoenixnap/go-sdk-bmc/bmcapi"
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

func GetRatedUsageStats(config BmcApiConfiguration) ([]RatedUsageStats, error) {
	apiClient := getBillingApiClient(config.ToClientCredentials())
	bmcClient := getBmcApiClient(config.ToClientCredentials())

	var stats []RatedUsageStats

	servers, r, err := bmcClient.ServersApi.ServersGet(getContext()).Execute()
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
			if ratedUsageStat := ConvertRatedUsageServerRecordToStats(ratedUsage.ServerRecord, servers); ratedUsageStat != nil {
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
			if ratedUsageStat := ConvertRatedUsageServerRecordToStats(ratedUsage.ServerRecord, servers); ratedUsageStat != nil {
				stats = append(stats, *ratedUsageStat)
			}
		}
	}

	return stats, nil
}

func ConvertRatedUsageServerRecordToStats(record *billingapi.ServerRecord, servers []bmcapi.Server) *RatedUsageStats {
	var ratedUsageStat RatedUsageStats

	if record.GetActive() {
		ratedUsageStat.Cost = float64(record.GetCost()) / float64(100)
		ratedUsageStat.Hostname = record.GetMetadata().Hostname
		ratedUsageStat.Location = string(record.GetLocation())
		ratedUsageStat.PriceModel = record.GetPriceModel()
		ratedUsageStat.ProductCategory = record.GetProductCategory()
		ratedUsageStat.ProductCode = record.GetProductCode()
		ratedUsageStat.YearMonth = record.GetYearMonth()

		s := GetServerFromList(record.GetId(), servers)
		ratedUsageStat.BillingTags = GetBillingTagsFromServer(s)

		return &ratedUsageStat
	}

	return nil
}

func GetServerFromList(serverId string, servers []bmcapi.Server) *bmcapi.Server {
	for _, s := range servers {
		if s.GetId() == serverId {
			return &s
		}
	}
	return nil
}

func GetBillingTagsFromServer(server *bmcapi.Server) []RatedUsageTagInfo {
	tags := make([]RatedUsageTagInfo, 0)

	if server != nil {
		for _, t := range server.GetTags() {
			if t.GetIsBillingTag() {
				tags = append(tags, RatedUsageTagInfo{
					Key:   t.GetName(),
					Value: t.GetValue(),
				})
			}
		}
	}

	return tags
}
