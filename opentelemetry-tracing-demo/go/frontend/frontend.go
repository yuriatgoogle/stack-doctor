package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	// "io/ioutil"
	// "google.golang.org/grpc/codes"
	"time"

	"github.com/gorilla/mux"

	"go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/global"

	// "go.opentelemetry.io/otel/api/key"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporter/trace/stackdriver"

	//"go.opentelemetry.io/otel/plugin/httptrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	projectID   = os.Getenv("PROJECT_ID")
	backendAddr = os.Getenv("BACKEND")
	location    = os.Getenv("LOCATION")
)

func mainHandler(w http.ResponseWriter, r *http.Request) {

	tr := global.TraceProvider().Tracer("OT-tracing-demo")

	client := http.DefaultClient
	ctx := distributedcontext.NewContext(context.Background())

	// var body []byte

	// create root span - works
	ctx, rootSpan := tr.Start(ctx, "incoming call")
	defer rootSpan.End()

	// how to create child span...?
	ctx, childSpan := tr.Start(ctx, "backend call", apitrace.SpanFromContext(ctx))

	// create request for backend call
	req, err := http.NewRequest("GET", backendAddr, nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
	childCtx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	req = req.WithContext(childCtx)

	// add span context to backend call and make request
	// format := &tracecontext.HTTPFormat{}
	// format.SpanContextToRequest(rootSpan.SpanContext(), req)
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// send backend request

	/* err := tr.WithSpan(ctx, "incoming call",  // root span here
		func(ctx context.Context) error {
			req, _ := http.NewRequest("GET", backendAddr, nil)

			ctx, req = httptrace.W3C(ctx, req)
			httptrace.Inject(ctx, req)

			fmt.Printf("Sending request...\n")
			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			body, err = ioutil.ReadAll(res.Body)
			_ = res.Body.Close()
			trace.SpanFromContext(ctx).SetStatus(codes.OK)

			return err
		})

	if err != nil {
		panic(err)
	}

	*/
	fmt.Printf("%v\n", res.Status) //change to status code from backend
}

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

	// handle root request
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)

	// TODO - add handler with propagation from OT
	//log.Fatal(http.ListenAndServe(":8081", handler))

	http.ListenAndServe(":8080", r)

}
