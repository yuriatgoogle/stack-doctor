package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return shutdown
}

func main() {
	initTracer()

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		span := trace.SpanFromContext(ctx)
		span.AddEvent("handling incoming request")

		fmt.Printf("OK")
	}

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(helloHandler), "Backend span")

	http.Handle("/", otelHandler)
	if env == "LOCAL" {
		http.ListenAndServe("localhost:8081", nil)
	} else {
		http.ListenAndServe(":8081", nil)
	}
}

/* func mainHandler(w http.ResponseWriter, req *http.Request) {
	/* tr := otel.Tracer("Backend")
	_, childSpan := tr.Start(req.Context(), "Backend Request", trace.WithAttributes(semconv.PeerServiceKey.String("Backend")))
	defer childSpan.End()
	ctx := req.Context()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handling backend")
	// childSpan.AddEvent("handling backend call")

	// output
	log.Printf("backend call received")
	fmt.Printf("OK")
}

func main() {
	shutdown := initTracer()
	defer shutdown()

	r := mux.NewRouter()
	otelHandler := otelhttp.NewHandler(mainHandler, "Hello")
	r.HandleFunc("/", otelHandler)

	if env == "LOCAL" {
		http.ListenAndServe("localhost:8081", r)
	} else {
		http.ListenAndServe(":8081", r)
	}
}
*/
