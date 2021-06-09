package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
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
	_, shutdown, err := texporter.InstallNewPipeline(
		[]texporter.Option{texporter.WithProjectID(projectID)},
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

func mainHandler(w http.ResponseWriter, r *http.Request) {

	tr := otel.Tracer("frontend")

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	ctx := baggage.ContextWithValues(context.Background(),
		attribute.String("username", "donuts"),
	)

	var body []byte

	err := func(ctx context.Context) error { // Root span is created here.

		// create child span
		ctx, span := tr.Start(ctx, "Incoming request", trace.WithAttributes(semconv.PeerServiceKey.String("Frontend")))
		defer span.End()

		// create new request with context
		req, _ := http.NewRequestWithContext(ctx, "GET", "http://"+backendAddr, nil)

		// make the backend request with context.
		fmt.Printf("Sending request...\n")
		span.AddEvent("Making backend call")
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		// Process response.
		body, err = ioutil.ReadAll(res.Body)
		_ = res.Body.Close()
		span.SetStatus(codes.Ok, "")

		// Output result.
		log.Printf("got response: %d\n", res.Status)
		fmt.Printf("%v\n", "OK") //change to status code from backend
		return err
	}(ctx)

	if err != nil {
		panic(err)
	}
}

func main() {
	shutdown := initTracer()
	defer shutdown()

	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)

	if env == "LOCAL" {
		http.ListenAndServe("localhost:8080", r)
	} else {
		http.ListenAndServe(":8080", r)
	}
}
