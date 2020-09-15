"use strict";

// imports
const express = require ('express');
const app = express();
const { MeterProvider } = require('@opentelemetry/metrics');
const { MetricExporter } = require('@google-cloud/opentelemetry-cloud-monitoring-exporter');

function sleep (n) {
  Atomics.wait(new Int32Array(new SharedArrayBuffer(4)), 0, 0, n);
}

const exporter = new MetricExporter();
const ERROR_RATE = process.env.ERROR_RATE;

// Register the exporter
const meter = new MeterProvider({
  exporter,
  interval: 60000,
}).getMeter('example-meter');

// define metrics 
const requestCount = meter.createCounter("request_count_sli", {
  monotonic: true,
  labelKeys: ["metricOrigin"],
  description: "Counts total number of requests"
});
const errorCount = meter.createCounter("error_count_sli", {
    monotonic: true,
    labelKeys: ["metricOrigin"],
    description: "Counts total number of errors"
});
// const labels = meter.labels({ metricOrigin: process.env.ENV});

// set metric values on request
app.get('/', (req, res) => {
    // start latency timer
    console.log('request made');
    // increment total requests counter
    requestCount.add(1);
    // return an error based on ERROR_RATE
    if ((Math.floor(Math.random() * 100)) <= ERROR_RATE) {
        // increment error counter
        errorCount.add(1);
        // return error code
        res.status(500).send("error!")
    } 
    // record latency and respond right away
    else {
      sleep(Math.floor(Math.random()*1000));
      res.status(200).send("success!")
    }
})

app.listen(8080, () => console.log(`Example app listening on port 8080!`))