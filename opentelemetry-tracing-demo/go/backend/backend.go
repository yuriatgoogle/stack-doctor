package main

import (
	// "fmt"
	"log"
	"net/http"
	"os"
	// "context"
	// "io/ioutil"
    "io"
	// "google.golang.org/grpc/codes"
	//"time"

	// "github.com/gorilla/mux"

	"go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/global"
	// "go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporter/trace/stackdriver"
	"go.opentelemetry.io/otel/plugin/httptrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	projectID   = os.Getenv("PROJECT_ID")
	backendAddr = os.Getenv("BACKEND")
	location    = os.Getenv("LOCATION")
)

func initTracer() {
// Create Stackdriver exporter to be able to retrieve
	// the collected spans.
	exporter, err := stackdriver.NewExporter(
		stackdriver.WithProjectID(projectID),
	)
	if err != nil {
		log.Fatal(err)
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter))
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)
}

func main() {
	initTracer()
	tr := global.TraceProvider().Tracer("backend")

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		attrs, entries, spanCtx := httptrace.Extract(req.Context(), req)

		req = req.WithContext(distributedcontext.WithMap(req.Context(), distributedcontext.NewMap(distributedcontext.MapUpdate{
			MultiKV: entries,
		})))

		ctx, span := tr.Start(
			req.Context(),
			"hello",
			trace.WithAttributes(attrs...),
			trace.ChildOf(spanCtx),
		)
		defer span.End()

		span.AddEvent(ctx, "handling this...")

		_, _ = io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/", helloHandler)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}