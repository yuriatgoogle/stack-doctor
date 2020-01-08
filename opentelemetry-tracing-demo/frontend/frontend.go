package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"context"
	"io/ioutil"
	"google.golang.org/grpc/codes"

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
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	
	tr := global.TraceProvider().Tracer("OT-tracing-demo")

	client := http.DefaultClient
	ctx := distributedcontext.NewContext(context.Background())
	
	var body []byte

	// send backend request
	
	err := tr.WithSpan(ctx, "incoming call",  // root span here
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
	fmt.Printf("%v\n", 200) //change to status code from backend
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
