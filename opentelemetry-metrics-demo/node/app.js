"use strict";

// imports
const express = require ('express')
const { MeterRegistry } = require('@opentelemetry/metrics');
const { PrometheusExporter } = require('@opentelemetry/exporter-prometheus');

function sleep (n) {
    Atomics.wait(new Int32Array(new SharedArrayBuffer(4)), 0, 0, n);
}

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


// set metric values on request

app.get('/', (req, res) => {
    // start latency timer
    const requestReceived = new Date().getTime();
    console.log('request made');
    // increment total requests counter
    requestCount.bind(labels).add(1);
    // return an error 2% of the time
    if ((Math.floor(Math.random() * 100)) > 98) {
        // increment error counter
        errorCount.bind(labels).add(1);
        // return error code
        res.status(500).send("error!")
    } 
    // make it slow an additional 2% of the time
    else if ((Math.floor(Math.random() * 100)) > 98) {
      // delay for a bit
      sleep(10000);
      // record response latency
      const measuredLatency = new Date().getTime() - requestReceived;
      responseLatency.bind(labels).set(measuredLatency)
      res.status(200).send("delayed success in " + measuredLatency + " ms")
    }
    // record latency and respond right away
    else {
      sleep(Math.floor(Math.random()*1000));
      const measuredLatency = new Date().getTime() - requestReceived;
      responseLatency.bind(labels).set(measuredLatency)
      res.status(200).send("did not delay - success in " + measuredLatency + " ms")
    }
})

app.listen(8080, () => console.log(`Example app listening on port 8080!`))