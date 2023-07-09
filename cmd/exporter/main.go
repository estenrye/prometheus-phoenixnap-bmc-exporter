package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	exporter "github.com/estenrye/prometheus-phoenix-nap-exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		bmcClientCredentialsFile = flag.String("configFile", "", "Location of the PhoenixNAP BMC API credentials file")
		bmcOauth2ClientId        = flag.String("clientId", "", "Client Id used to authenticate to the PhoenixNAP BMC API")
		bmcOauth2ClientSecret    = flag.String("clientSecret", "", "Client Secret used to authenticate to the PhoenixNAP BMC API")
		bmcApiTokenUrl           = flag.String("tokenUrl", "", "Base Url for the PhoenixNAP BMC API")
		promPort                 = flag.Int("prometheus.port", 9150, "port to expose prometheus metrics")
		collectGoMetrics         = flag.Bool("go.collector.enabled", false, "flag to enable go collector metrics")
		logFormat                = flag.String("log.format", "", "Selects the log format. Expects: 'json' or 'text'")
		logLevel                 = flag.String("log.level", "", "Sets the minimum level of logs displayed. Expects: 'panic', 'fatal', 'warning', 'info', 'debug', or 'trace'")
	)
	flag.Parse()

	bmc_exporter_configuration := exporter.NewBmcApiConfiguration(
		*bmcClientCredentialsFile,
		*bmcOauth2ClientId,
		*bmcOauth2ClientSecret,
		*bmcApiTokenUrl).SetLogFormat(*logFormat).SetLogLevel(*logLevel)

	log.SetFormatter(bmc_exporter_configuration.GetLogFormatter())
	log.SetLevel(bmc_exporter_configuration.GetLogLevel())

	log.WithField("bmc_exporter_configuration", bmc_exporter_configuration).Trace("BMC API configuration", bmc_exporter_configuration)

	quotaStats := func() ([]exporter.QuotaStats, error) {
		return exporter.GetBmcQuotas(*bmc_exporter_configuration)
	}
	reservationStats := func() ([]exporter.ReservationStats, error) {
		return exporter.GetBmcReservations(*bmc_exporter_configuration)
	}
	ratedUsageStats := func() ([]exporter.RatedUsageStats, error) {
		return exporter.GetRatedUsageStats(*bmc_exporter_configuration)
	}

	qc := exporter.NewQuotaCollector(quotaStats)
	rc := exporter.NewReservationCollector(reservationStats)
	ru := exporter.NewRatedUsageCollector(ratedUsageStats)

	reg := prometheus.NewRegistry()
	if *collectGoMetrics {
		reg.MustRegister(collectors.NewGoCollector())
	}
	reg.MustRegister(qc)
	reg.MustRegister(rc)
	reg.MustRegister(ru)

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	mux.Handle("/metrics", promHandler)

	port := fmt.Sprintf(":%d", *promPort)
	log.Printf("starting PhoenixNAP BMC Exporter on %q/metrics", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("cannot start PhoenixNAP BMC Exporter: %s", err)
	}
}
