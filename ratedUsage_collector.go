package exporter

import (
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &ratedUsageCollector{}

type ratedUsageCollector struct {
	RatedUsageCost *prometheus.Desc
	BillingTags    []string
	stats          func() ([]RatedUsageStats, error)
}

func NewRatedUsageCollector(stats func() ([]RatedUsageStats, error), billingTags []string) prometheus.Collector {
	labels := []string{"productCode", "locationCode", "hostname", "priceModel", "productCategory", "yearMonth"}
	labels = append(labels, billingTags...)
	log.WithField("labels", labels).Debug("NewRatedUsageCollector load labels")
	q := ratedUsageCollector{
		RatedUsageCost: prometheus.NewDesc(
			"bmc_rated_usage_cost_total",
			"Total rated usage cost of actively provisioned resources.",
			labels,
			nil,
		),
		BillingTags: billingTags,
		stats:       stats,
	}

	log.WithField("RatedUsageCollector", q).Info("Created New Rated Usage Collector")
	return &q
}

func (c *ratedUsageCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.RatedUsageCost,
	}

	for _, d := range ds {
		ch <- d
	}
}

func (c *ratedUsageCollector) Collect(ch chan<- prometheus.Metric) {
	stats, err := c.stats()
	if err != nil {
		log.WithError(err).Error("Error encountered when collecting metric.")
		ch <- prometheus.NewInvalidMetric(c.RatedUsageCost, err)
		return
	}

	for _, s := range stats {
		ch <- s.ToPrometheusMetric(c.RatedUsageCost, c.BillingTags)
	}
}
