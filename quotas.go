package exporter

import log "github.com/sirupsen/logrus"

type QuotaStats struct {
	ServersMaxLimit          float64
	ServersUsedCount         float64
	PublicIpMaxLimit         float64
	PublicIpUsedCount        float64
	StorageNetworksMaxLimit  float64
	StorageNetworksUsedCount float64
}

func GetBmcQuotas(config BmcApiConfiguration) ([]QuotaStats, error) {
	apiClient := getBmcApiClient(config.ToClientCredentials())

	var stats []QuotaStats
	var q QuotaStats

	resp, r, err := apiClient.QuotasApi.QuotasGet(getContext()).Execute()
	if err != nil {
		log.WithField("HttpResponse", r).WithError(err).Error("Error when calling `QuotasApi.QuotasGet`.")
		return stats, err
	}

	for _, quota := range resp {
		if quota.Id == "bmc.servers.max_count" {
			q.ServersMaxLimit = float64(quota.Limit)
			q.ServersUsedCount = float64(quota.Used)
		} else if quota.Id == "bmc.public_ips.max_count" {
			q.PublicIpMaxLimit = float64(quota.Limit)
			q.PublicIpUsedCount = float64(quota.Used)
		} else if quota.Id == "bmc.storage_network.max_capacity" {
			q.StorageNetworksMaxLimit = float64(quota.Limit)
			q.StorageNetworksUsedCount = float64(quota.Used)
		}
	}

	stats = append(stats, q)

	log.WithField("stats", stats).Debug("Quota Retrieved from `QuotasApi.QuotasGet`")

	return stats, nil
}
