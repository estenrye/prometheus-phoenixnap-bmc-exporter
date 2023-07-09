package exporter

import (
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

func GetRatedUsageStats(config BmcApiConfiguration) ([]RatedUsageStats, error) {
	apiClient := getBillingApiClient(config.ToClientCredentials())

	var stats []RatedUsageStats

	resp, r, err := apiClient.RatedUsageApi.RatedUsageMonthToDateGet(getContext()).ProductCategory(billingapi.SERVER).Execute()
	if err != nil {
		log.WithField("HttpResponse", r).WithError(err).Error("Error when calling `RatedUsageApi.RatedUsageMonthToDateGet`.")
		return stats, err
	}

	for _, ratedUsage := range resp {
		var ratedUsageStat RatedUsageStats

		if ratedUsage.ServerRecord.GetActive() {
			ratedUsageStat.Cost = float64(ratedUsage.ServerRecord.GetCost()) / float64(100)
			ratedUsageStat.Hostname = ratedUsage.ServerRecord.GetMetadata().Hostname
			ratedUsageStat.Location = string(ratedUsage.ServerRecord.GetLocation())
			ratedUsageStat.PriceModel = ratedUsage.ServerRecord.GetPriceModel()
			ratedUsageStat.ProductCategory = ratedUsage.ServerRecord.GetProductCategory()
			ratedUsageStat.ProductCode = ratedUsage.ServerRecord.GetProductCode()
			ratedUsageStat.YearMonth = ratedUsage.ServerRecord.GetYearMonth()
			ratedUsageStat.BillingTags = make([]RatedUsageTagInfo, 0)

			stats = append(stats, ratedUsageStat)
		}
	}
	return stats, nil
}
