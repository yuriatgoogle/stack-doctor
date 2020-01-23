package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel/api/global"
	// "go.opentelemetry.io/otel/api/key"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporter/metric/prometheus"
)

var (
	env = os.Getenv("ENV")
)

func main() {

	pusher, mainHandler, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}
	if pusher == nil {
		log.Printf("no pusher")
	} else {
		log.Printf("have pusher")
	}
	defer pusher.Stop()

	meter := global.MeterProvider().Meter("custom.googleapis.com/opentelemetry-metrics-demo")
	if meter != nil {
		log.Printf("have meter")
	}

	gaugeMetric := meter.NewFloat64Gauge("gauge",
		metric.WithDescription("A gauge set to 1.0"),
	)

	ctx := context.Background()
	meter.
		meter.RecordBatch(
		ctx,
		gaugeMetric.Measurement(1.0),
	)

	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)

	if env == "LOCAL" {
		http.ListenAndServe("localhost:8080", r)
	} else {
		http.ListenAndServe(":8080", r)
	}

}
