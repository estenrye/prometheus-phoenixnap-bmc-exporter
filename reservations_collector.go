package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var _ prometheus.Collector = &reservationCollector{}

type reservationCollector struct {
	BmcTotalServerReservationsCount *prometheus.Desc
	BmcTotalServerReservationsCost  *prometheus.Desc
	stats                           func() ([]ReservationStats, error)
}

func NewReservationCollector(stats func() ([]ReservationStats, error)) prometheus.Collector {
	q := reservationCollector{
		BmcTotalServerReservationsCount: prometheus.NewDesc(
			"bmc_total_server_reservations_count",
			"Total number of SERVER reservations by productcode, location, status, autorenewal, and assignment",
			[]string{"productCode", "locationCode", "status", "autorenewal", "assignment"},
			nil,
		),
		BmcTotalServerReservationsCost: prometheus.NewDesc(
			"bmc_total_server_reservations_cost",
			"Total number of SERVER reservations by productcode, location, status, autorenewal, and assignment",
			[]string{"productCode", "locationCode", "status", "autorenewal", "assignment"},
			nil,
		),
		stats: stats,
	}

	log.WithField("QuotaCollector", q).Info("Created New Quota Collector")
	return &q
}

func (c *reservationCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.BmcTotalServerReservationsCount,
		c.BmcTotalServerReservationsCost,
	}

	for _, d := range ds {
		ch <- d
	}
}

func (c *reservationCollector) Collect(ch chan<- prometheus.Metric) {
	stats, err := c.stats()
	if err != nil {
		log.WithError(err).Error("Error encountered when collecting metric.")
		ch <- prometheus.NewInvalidMetric(c.BmcTotalServerReservationsCount, err)
		ch <- prometheus.NewInvalidMetric(c.BmcTotalServerReservationsCost, err)
		return
	}

	for _, s := range stats {
		for _, v := range s.TotalReservations {
			ch <- prometheus.MustNewConstMetric(
				c.BmcTotalServerReservationsCount,
				prometheus.GaugeValue,
				v.Total,
				v.GetProductCodeLabel(),
				v.GetLocationLabel(),
				v.GetStatusLabel(),
				v.GetAutoRenewalLabel(),
				v.GetAssignmentLabel(),
			)
			ch <- prometheus.MustNewConstMetric(
				c.BmcTotalServerReservationsCost,
				prometheus.GaugeValue,
				v.MonthlyCostTotal,
				v.GetProductCodeLabel(),
				v.GetLocationLabel(),
				v.GetStatusLabel(),
				v.GetAutoRenewalLabel(),
				v.GetAssignmentLabel(),
			)
		}
	}
}
