package exporter

import (
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &quotaCollector{}

type quotaCollector struct {
	PublicIpMaxLimit         *prometheus.Desc
	PublicIpUsedCount        *prometheus.Desc
	PublicIpUsedRate         *prometheus.Desc
	ServersMaxLimit          *prometheus.Desc
	ServersUsedCount         *prometheus.Desc
	ServersUsedRate          *prometheus.Desc
	StorageNetworksMaxLimit  *prometheus.Desc
	StorageNetworksUsedCount *prometheus.Desc
	StorageNetworksUsedRate  *prometheus.Desc

	stats func() ([]QuotaStats, error)
}

func NewQuotaCollector(stats func() ([]QuotaStats, error)) prometheus.Collector {
	q := quotaCollector{
		PublicIpMaxLimit: prometheus.NewDesc(
			"bmc_quota_public_ip_max_limit",
			"Maximum number of Public IP Addresses that can be provisioned on the PhoenixNAP BMC API.",
			[]string{},
			nil,
		),
		PublicIpUsedCount: prometheus.NewDesc(
			"bmc_quota_public_ip_used_count",
			"Number of Public IP Addresses that are currently provisioned on the PhoenixNAP BMC API.",
			[]string{},
			nil,
		),
		PublicIpUsedRate: prometheus.NewDesc(
			"bmc_quota_public_ip_used_rate",
			"Percentage of Public IP Addresses that are currently provisioned vs. the number available to provision.",
			[]string{},
			nil,
		),
		ServersMaxLimit: prometheus.NewDesc(
			"bmc_quota_servers_max_limit",
			"Maximum number of servers that can be provisioned on the PhoenixNAP BMC API.",
			[]string{},
			nil,
		),
		ServersUsedCount: prometheus.NewDesc(
			"bmc_quota_servers_used_count",
			"Number of servers that are currently provisioned on the PhoenixNAP BMC API.",
			[]string{},
			nil,
		),
		ServersUsedRate: prometheus.NewDesc(
			"bmc_quota_servers_used_rate",
			"Percentage of Serverss that are currently provisioned vs. the number available to provision.",
			[]string{},
			nil,
		),
		StorageNetworksMaxLimit: prometheus.NewDesc(
			"bmc_quota_storage_network_max_limit",
			"Maximum amount of network data storage in GB that can be provisioned on the PhoenixNAP BMC API.",
			[]string{},
			nil,
		),
		StorageNetworksUsedCount: prometheus.NewDesc(
			"bmc_quota_storage_network_used_count",
			"Amount of network data storage in GB of that is currently provisioned on the PhoenixNAP BMC API.",
			[]string{},
			nil,
		),
		StorageNetworksUsedRate: prometheus.NewDesc(
			"bmc_quota_storage_network_used_rate",
			"Percentage of Public IP Addresses that are currently provisioned vs. the number available to provision.",
			[]string{},
			nil,
		),
		stats: stats,
	}

	log.WithField("QuotaCollector", q).Info("Created New Quota Collector")
	return &q
}

func (c *quotaCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.PublicIpMaxLimit,
		c.PublicIpMaxLimit,
		c.PublicIpUsedRate,
		c.ServersMaxLimit,
		c.ServersUsedCount,
		c.ServersUsedRate,
		c.StorageNetworksMaxLimit,
		c.StorageNetworksUsedCount,
		c.StorageNetworksUsedRate,
	}

	for _, d := range ds {
		ch <- d
	}
}

func (c *quotaCollector) Collect(ch chan<- prometheus.Metric) {
	stats, err := c.stats()
	if err != nil {
		log.WithError(err).Error("Error encountered when collecting metric.")
		ch <- prometheus.NewInvalidMetric(c.PublicIpMaxLimit, err)
		ch <- prometheus.NewInvalidMetric(c.PublicIpUsedCount, err)
		ch <- prometheus.NewInvalidMetric(c.PublicIpUsedRate, err)
		ch <- prometheus.NewInvalidMetric(c.ServersMaxLimit, err)
		ch <- prometheus.NewInvalidMetric(c.ServersUsedCount, err)
		ch <- prometheus.NewInvalidMetric(c.ServersUsedRate, err)
		ch <- prometheus.NewInvalidMetric(c.StorageNetworksMaxLimit, err)
		ch <- prometheus.NewInvalidMetric(c.StorageNetworksUsedCount, err)
		ch <- prometheus.NewInvalidMetric(c.StorageNetworksUsedRate, err)
		return
	}

	for _, s := range stats {
		ch <- prometheus.MustNewConstMetric(
			c.PublicIpMaxLimit,
			prometheus.GaugeValue,
			s.PublicIpMaxLimit,
		)
		ch <- prometheus.MustNewConstMetric(
			c.PublicIpUsedCount,
			prometheus.GaugeValue,
			s.PublicIpUsedCount,
		)
		ch <- prometheus.MustNewConstMetric(
			c.PublicIpUsedRate,
			prometheus.GaugeValue,
			s.PublicIpUsedRate,
		)
		ch <- prometheus.MustNewConstMetric(
			c.ServersMaxLimit,
			prometheus.GaugeValue,
			s.ServersMaxLimit,
		)
		ch <- prometheus.MustNewConstMetric(
			c.ServersUsedCount,
			prometheus.GaugeValue,
			s.ServersUsedCount,
		)
		ch <- prometheus.MustNewConstMetric(
			c.ServersUsedRate,
			prometheus.GaugeValue,
			s.ServersUsedRate,
		)
		ch <- prometheus.MustNewConstMetric(
			c.StorageNetworksMaxLimit,
			prometheus.GaugeValue,
			s.StorageNetworksMaxLimit,
		)
		ch <- prometheus.MustNewConstMetric(
			c.StorageNetworksUsedCount,
			prometheus.GaugeValue,
			s.StorageNetworksUsedCount,
		)
		ch <- prometheus.MustNewConstMetric(
			c.StorageNetworksUsedRate,
			prometheus.GaugeValue,
			s.StorageNetworksUsedRate,
		)
	}
}
