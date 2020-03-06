const express = require('express');
const app = express();
const prometheus = require('prom-client');
const sleep = require('sleep');
const collectDefaultMetrics = prometheus.collectDefaultMetrics;

// define golden signal metrics
// total requests - counter
const nodeRequestsCounter = new prometheus.Counter({
    name: 'node_requests',
    help: 'total requests'
});

// failed requests - counter
const nodeFailedRequestsCounter = new prometheus.Counter({
    name: 'node_failed_requests',
    help: 'failed requests'
});

// latency - histogram
const nodeLatenciesHistogram = new prometheus.Histogram({
    name: 'node_request_latency',
    help: 'request latency by path',
    labelNames: ['route'],
    buckets: [100, 400]
});

// Probe every 5th second.
collectDefaultMetrics({ timeout: 1000 });

app.get('/', (req, res) => {
    // start latency timer
    const requestReceived = new Date().getTime();
    console.log('request made');
    // increment total requests counter
    nodeRequestsCounter.inc();
    // return an error 1% of the time
    if ((Math.floor(Math.random() * 100)) == 100) {
        // increment error counter
        nodeFailedRequestsCounter.inc();
        // return error code
        res.send("error!", 500);
    } 
    else {
        // delay for a bit
        sleep.msleep((Math.floor(Math.random() * 1000)));
        // record response latency
        const responseLatency = new Date().getTime() - requestReceived;
        nodeLatenciesHistogram
            .labels(req.route.path)
            .observe(responseLatency);
        res.send("success in " + responseLatency + " ms");
    }
})

app.get('/metrics', (req, res) => {
    res.set('Content-Type', prometheus.register.contentType)
    res.end(prometheus.register.metrics())
  })

app.listen(8080, () => console.log(`Example app listening on port 8080!`))