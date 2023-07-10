package exporter

import (
	"time"

	"github.com/phoenixnap/go-sdk-bmc/billingapi"
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

	var stats []RatedUsageStats

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
			if ratedUsageStat := ConvertRatedUsageServerRecordToStats(ratedUsage.ServerRecord, config); ratedUsageStat != nil {
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
			if ratedUsageStat := ConvertRatedUsageServerRecordToStats(ratedUsage.ServerRecord, config); ratedUsageStat != nil {
				stats = append(stats, *ratedUsageStat)
			}
		}
	}

	return stats, nil
}

func ConvertRatedUsageServerRecordToStats(record *billingapi.ServerRecord, config BmcApiConfiguration) *RatedUsageStats {
	var ratedUsageStat RatedUsageStats

	if record.GetActive() {
		ratedUsageStat.Cost = float64(record.GetCost()) / float64(100)
		ratedUsageStat.Hostname = record.GetMetadata().Hostname
		ratedUsageStat.Location = string(record.GetLocation())
		ratedUsageStat.PriceModel = record.GetPriceModel()
		ratedUsageStat.ProductCategory = record.GetProductCategory()
		ratedUsageStat.ProductCode = record.GetProductCode()
		ratedUsageStat.YearMonth = record.GetYearMonth()
		ratedUsageStat.BillingTags = GetBillingTagsFromInstance(record.GetMetadata().Id, config)

		log.WithField("ratedUsageStat", ratedUsageStat).Trace("Rated Usage Stat")

		return &ratedUsageStat
	}

	return nil
}

func GetBillingTagsFromInstance(serverId string, config BmcApiConfiguration) []RatedUsageTagInfo {
	var tags []RatedUsageTagInfo

	apiClient := getBmcApiClient(config.ToClientCredentials())
	resp, r, err := apiClient.ServersApi.ServersServerIdGet(getContext(), serverId).Execute()
	if err != nil {
		log.WithField("HttpResponse", r).WithError(err).Error("Error when calling `ServersApi.ServersServerIdGet`.")
		return tags
	}

	for _, tag := range resp.GetTags() {
		var ratedUsageTag RatedUsageTagInfo

		if tag.GetIsBillingTag() {
			ratedUsageTag.Key = tag.GetName()
			ratedUsageTag.Value = tag.GetValue()
			tags = append(tags, ratedUsageTag)
		}
	}

	return tags
}
