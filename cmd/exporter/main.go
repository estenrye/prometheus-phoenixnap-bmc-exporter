package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		promPort         = flag.Int("prometheus.port", 9150, "port to expose prometheus metrics")
		collectGoMetrics = flag.Bool("go.collector.enabled", false, "flag to enable go collector metrics")
	)
	flag.Parse()

	reg := prometheus.NewRegistry()
	if *collectGoMetrics {
		reg.MustRegister(collectors.NewGoCollector())
	}

	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	mux.Handle("/metrics", promHandler)

	port := fmt.Sprintf(":%d", *promPort)
	log.Printf("starting PhoenixNAP BMC Exporter on %q/metrics", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("cannot start PhoenixNAP BMC Exporter: %s", err)
	}
}
