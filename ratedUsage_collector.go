package exporter

import (
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &ratedUsageCollector{}

type ratedUsageCollector struct {
	RatedUsageCost *prometheus.Desc
	stats          func() ([]RatedUsageStats, error)
}

func NewRatedUsageCollector(stats func() ([]RatedUsageStats, error)) prometheus.Collector {
	q := ratedUsageCollector{
		RatedUsageCost: prometheus.NewDesc(
			"bmc_rated_usage_cost_total",
			"Total rated usage cost of actively provisioned resources.",
			[]string{"productCode", "locationCode", "hostname", "priceModel", "productCategory", "yearMonth"},
			nil,
		),
		stats: stats,
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
		ch <- prometheus.MustNewConstMetric(
			c.RatedUsageCost,
			prometheus.CounterValue,
			s.Cost,
			s.ProductCode,
			s.Location,
			s.Hostname,
			s.PriceModel,
			s.ProductCategory,
			s.YearMonth,
		)
	}
}
