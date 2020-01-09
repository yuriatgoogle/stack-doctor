package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	// "context"
	// "io/ioutil"
    // "io"
	// "google.golang.org/grpc/codes"
	//"time"

	"github.com/gorilla/mux"

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
    env			= os.Getenv("ENV")
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

func mainHandler(w http.ResponseWriter, req *http.Request) {
    // start tracer
    tr := global.TraceProvider().Tracer("backend")
    // get context from incoming request
    attrs, entries, spanCtx := httptrace.Extract(req.Context(), req)

    // create request using context
    req = req.WithContext(distributedcontext.WithMap(req.Context(), distributedcontext.NewMap(distributedcontext.MapUpdate{
        MultiKV: entries,
    })))

    // create span
    ctx, span := tr.Start(
        req.Context(),
        "backend call received",
        trace.WithAttributes(attrs...),
        trace.ChildOf(spanCtx),
    )

    span.AddEvent(ctx, "handling backend call")

    // output
    log.Printf("backend call received")
    fmt.Printf("OK")
    // close span
    span.End()
}

func main() {
	initTracer()

	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)
	
	if (env=="LOCAL") {
		http.ListenAndServe("localhost:8081", r)
	} else {
		http.ListenAndServe(":8081", r)
	}
}