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

	"cloud.google.com/go/logging"
	logs "github.com/GoogleCloudPlatform/opencensus-spanner-demo/applog"

	trace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporter/trace/stackdriver"
	"go.opentelemetry.io/otel/plugin/httptrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	projectID     = os.Getenv("PROJECT_ID")
	backendAddr   = "https://www.google.com"
	location      = os.Getenv("LOCATION")
	env           = os.Getenv("ENV")
	loggingClient *logging.Client
)

const LOGNAME string = "ot-logs-traces"

func mainHandler(w http.ResponseWriter, r *http.Request) {

	tr := global.TraceProvider().Tracer("OT-tracing-demo")

	client := http.DefaultClient
	ctx := distributedcontext.NewContext(context.Background())

	var body []byte

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
			printWithTrace(ctx, "The process took %d seconds\n", n)

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
			// log query
			printWithTrace(ctx, "Backend response: %d\n", res.StatusCode)
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

// Send to Cloud Logging service including reference to current span
func printWithTrace(ctx context.Context, format string, v ...interface{}) {
	printf(ctx, logging.Info, format, v...)
}

// Send to Cloud Logging service including reference to current span
func printf(ctx context.Context, severity logging.Severity, format string,
	v ...interface{}) {
	span := trace.SpanFromContext(ctx)
	sCtx := span.SpanContext()
	tr := sCtx.TraceIDString()
	lg := loggingClient.Logger(LOGNAME)
	trace := fmt.Sprintf("projects/%s/traces/%s", projectID, tr)
	lg.Log(logging.Entry{
		Severity: severity,
		Payload:  fmt.Sprintf(format, v...),
		Trace:    trace,
		SpanID:   sCtx.SpanIDString(),
	})
}

func main() {
	initTracer()
	logs.Initialize(projectID)
	defer logs.Close()
	initLogger()
	defer closeLogger()

	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)

	if env == "LOCAL" {
		http.ListenAndServe("localhost:8080", r)
	} else {
		http.ListenAndServe(":8080", r)
	}
}

// helper functions down here
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

// Initialize the Cloud Logging client
func initLogger() {
	ctx := context.Background()
	var err error
	loggingClient, err = logging.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("Failed to create logging client: %v", err)
		return
	}
	fmt.Printf("Stackdriver Logging initialized with project id %s, see Cloud "+
		" Console under GCE VM instance > all instance_id\n", projectID)
}

func closeLogger() {
	err := loggingClient.Close()
	if err != nil {
		fmt.Printf("Failed to close logging client: %v", err)
	}
}
