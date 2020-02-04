package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/global"

	// "go.opentelemetry.io/otel/api/key"
	logs "github.com/GoogleCloudPlatform/opencensus-spanner-demo/applog"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporter/trace/stackdriver"
	"go.opentelemetry.io/otel/plugin/httptrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	projectID   = os.Getenv("PROJECT_ID")
	backendAddr = "https://www.google.com"
	location    = os.Getenv("LOCATION")
	env         = os.Getenv("ENV")
)

func mainHandler(w http.ResponseWriter, r *http.Request) {

	tr := global.TraceProvider().Tracer("OT-tracing-demo")

	client := http.DefaultClient
	ctx := distributedcontext.NewContext(context.Background())

	var body []byte // try deleting?

	err := tr.WithSpan(ctx, "incoming call", // root span here
		func(ctx context.Context) error {

			// create span for internal processing and backend
			ctx, processSpan := tr.Start(ctx, "process and query")
			processSpan.AddEvent(ctx, "process and query")
			// do some delay
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(10) // n will be between 0 and 10
			log.Printf("sleeping for: %d\n", n)
			time.Sleep(time.Duration(n) * time.Second)
			logs.Printf(ctx, "The process took %d seconds\n", n)
			// create another child span for the query
			log.Printf("making backend request")
			ctx, querySpan := tr.Start(ctx, "query")
			querySpan.AddEvent(ctx, "query")

			// create backend request
			req, _ := http.NewRequest("GET", backendAddr, nil)
			// inject context
			ctx, req = httptrace.W3C(ctx, req)
			httptrace.Inject(ctx, req)
			// do request
			log.Printf("Sending request...\n")
			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			body, err = ioutil.ReadAll(res.Body)
			_ = res.Body.Close()

			// finish tracing
			trace.SpanFromContext(ctx).SetStatus(codes.OK)
			querySpan.End()
			processSpan.End()

			// output
			log.Printf("got response: %d\n", res.Status)
			fmt.Printf("%v\n", "OK") //change to status code from backend
			return err
		})

	if err != nil {
		panic(err)
	}
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

	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)

	if env == "LOCAL" {
		http.ListenAndServe("localhost:8080", r)
	} else {
		http.ListenAndServe(":8080", r)
	}
}
