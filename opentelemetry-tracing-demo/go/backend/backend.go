package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
)

var (
	projectID   = os.Getenv("PROJECT_ID")
	backendAddr = os.Getenv("BACKEND")
	location    = os.Getenv("LOCATION")
	env         = os.Getenv("ENV")
)

func initTracer() func() {
	projectID := os.Getenv("PROJECT_ID")

	// Create Google Cloud Trace exporter to be able to retrieve
	// the collected spans.
	_, shutdown, err := cloudtrace.InstallNewPipeline(
		[]cloudtrace.Option{cloudtrace.WithProjectID(projectID)},
		// For this example code we use sdktrace.AlwaysSample sampler to sample all traces.
		// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	if err != nil {
		log.Fatal(err)
	}
	return shutdown
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("server", "handling this..."))

	_, _ = io.WriteString(w, "Hello, world!\n")
}

func main() {
	shutdown := initTracer()
	defer shutdown()

	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)

	if env == "LOCAL" {
		http.ListenAndServe("localhost:8081", r)
	} else {
		http.ListenAndServe(":8081", r)
	}
}
