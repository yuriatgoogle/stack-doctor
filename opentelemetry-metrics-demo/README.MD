# Custom metrics in Node.js with OpenTelemetry (and Prometheus)
In my last [post](https://dev.to/yurigrinshteyn/distributed-tracing-with-opentelemetry-in-go-473h), I tackled my first project with [OpenTelemetry](https://opentelemetry.io) and built a basic demo to show how to use distributed tracing and the Stackdriver exporter.  I chose Go for that exercise because at the time it was the only language that had a Stackdriver exporter for tracing available.  This time, I wanted to attempt using OpenTelemetry for metric instrumentation and noticed that opentelemetry-go does not have a Stackdriver exporter for metrics ready yet.  I attempted to use the Prometheus exporter instead, but could not figure out how to make it play nice with the Mux router and switched to Node.js instead.  

So - here's how I built a basic Hello World Node.js app and instrumented it with OpenTelemetry to expose basic "golden signals" metrics (request count, error count, and latency) to Prometheus.

## The app and instrumentation
As always, I've made my code available in a Github [repo](https://github.com/yuriatgoogle/stack-doctor/opentelemetry-metrics/demo).  Let's review the interesting parts.

### Imports and setup
```go
// set up prometheus 
const prometheusPort = 8081;
const app = express();
const meter = new MeterRegistry().getMeter('example-prometheus');
const exporter = new PrometheusExporter(
  {
    startServer: true,
    port: prometheusPort
  },
  () => {
    console.log("prometheus scrape endpoint: http://localhost:"
      + prometheusPort 
      + "/metrics");
  }
);
meter.addExporter(exporter);
```

I'm only using three external packages - Express to handle requests and two OpenTelemetry packages - one to write the metrics, and another to export them to Prometheus.  Requests to `:8081/metrics` will be handled by Prometheus.  
> Note that I've overridden the default port that's set in the Prometheus exporter [code](https://github.com/open-telemetry/opentelemetry-js/blob/master/packages/opentelemetry-exporter-prometheus/src/prometheus.ts#L32-L38).

### Metric definition
```go
// define metrics with description and labels
const requestCount = meter.createCounter("request_count", {
  monotonic: true,
  labelKeys: ["metricOrigin"],
  description: "Counts total number of requests"
});
const errorCount = meter.createCounter("error_count", {
    monotonic: true,
    labelKeys: ["metricOrigin"],
    description: "Counts total number of errors"
});
const responseLatency = meter.createGauge("response_latency", {
    monotonic: false,
    labelKeys: ["metricOrigin"],
    description: "Records latency of response"
});
const labels = meter.labels({ metricOrigin: process.env.ENV});
```
In the next part, I define the three metrics I want to track.  Request count and error count are both monotonic counters, in that they can only be increased.  Response latency is a non-monotonic gauge, since I want to track latency for every request.  I am also using a single label - I could not figure how to get metrics to work properly with no labels (in the next section).

### Request handling
```go
app.get('/', (req, res) => {
    // start latency timer
    const requestReceived = new Date().getTime();
    console.log('request made');
    // increment total requests counter
    requestCount.bind(labels).add(1);
    // return an error 1% of the time
    if ((Math.floor(Math.random() * 100)) > 50) {
        // increment error counter
        errorCount.bind(labels).add(1);
        // return error code
        res.status(500).send("error!")
    } 
    else {
        // delay for a bit
        sleep(Math.floor(Math.random()*10000));
        // record response latency
        const measuredLatency = new Date().getTime() - requestReceived;
        responseLatency.bind(labels).set(measuredLatency)
        res.status(200).send("success in " + measuredLatency + " ms")
    }
})

app.listen(8080, () => console.log(`Example app listening on port 8080!`))
```
Finally, I'm ready to accept incoming requests.  Requests to `:8081/metrics` will be handled by Prometheus, and I'm using Express to handle requests to `:8080/`. For every request, I create a timer (`measuredLatency`) and keep it running until I'm ready to respond.  I generate a random number between 0 and 100 and use it to return an error half the time, in which case latency is not reported.  If there's no error to return, I `sleep` for another random value (between 0 and 10 seconds), take another time snapshot, calculate the elapsed time, and return that value as the latency.

### Metrics endpoint
When I run the app locally, I can then hit the `/metrics` endpoint and see the metrics I've registered there.

![image](https://github.com/yuriatgoogle/stack-doctor/raw/master/opentelemetry-metrics-demo/images/metricsendpoint.png))

# Deploying the app
Now I'm ready to run the app and generate actual traffic to it.  I'm going to 
* create a GKE cluster
* install Prometheus on it
* deploy the app and configure the `/metrics` endpoint to be scraped

Let's dive in!

## Creating the cluster
I took a very straightforward approach and used the Console to create a single-node cluster using an n1-standard node. I then used `gcloud container clusters get-credentials <cluster> --zone=<zone>` to get credentials for it to be able to access it from my local workstaton.

```bash
> kubectl get nodes
NAME                                     STATUS   ROLES    AGE     VERSION
gke-small-cluster-pool-1-b40da891-fc1p   Ready    <none>   7h27m   v1.15.7-gke.23
```
## Installing Prometheus
I took the simple option here:
```bash
brew install helm # install helm, the package manager for Kubernetes
helm repo add stable https://kubernetes-charts.storage.googleapis.com/
helm repo update
helm install prometheus stable/prometheus
```
That was enough to get Prometheus up and running, but I had to increase the size of my cluster to get enough CPU to run everything.
```bash
> kubectl get pods
NAME                                            READY   STATUS    RESTARTS   AGE
prometheus-alertmanager-77995478d8-lh4c5        2/2     Running   0          6m34s
prometheus-kube-state-metrics-f96cbbf97-vc56g   1/1     Running   0          6m34s
prometheus-node-exporter-8xlwp                  1/1     Running   0          114s
prometheus-node-exporter-s2pnt                  1/1     Running   0          6m34s
prometheus-node-exporter-xbtrs                  1/1     Running   0          112s
prometheus-pushgateway-747dc5fb57-rw72p         1/1     Running   0          6m34s
prometheus-server-6f7467bfdd-5vfgx              2/2     Running   0          6m34s
```
I verified that Prometheus was working by using `port-forward`:

![image]()