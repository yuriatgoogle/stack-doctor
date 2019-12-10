const express = require('express');
const app = express();

// opencensus setup
const {globalStats, MeasureUnit, AggregationType} = require('@opencensus/core');
const {StackdriverStatsExporter} = require('@opencensus/exporter-stackdriver');

// Stackdriver export interval is 60 seconds
const EXPORT_INTERVAL = 60;


// define the "golden signals" metrics
// request count
const REQUEST_COUNT = globalStats.createMeasureInt64(
    'request_count',
    MeasureUnit.UNIT,
    'Number of requests to the server'
);

// error count
const ERROR_COUNT = globalStats.createMeasureInt64(
    'error_count',
    MeasureUnit.UNIT,
    'Number of failed requests to the server'
);

// response latency
const RESPONSE_LATENCY = globalStats.createMeasureInt64(
    'response_latency',
    MeasureUnit.MS,
    'The server response latency in milliseconds'
  );

//create and register the view
const view = globalStats.createView(
    'response_latency',
    LATENCY_MS,
    AggregationType.LAST_VALUE,
    [],
    'The distribution of the task latencies.',
    // Latency in buckets:
    // [>=0ms, >=100ms, >=200ms, >=400ms, >=1s, >=2s, >=4s]
    [0, 100, 200, 400, 1000, 2000, 4000]
  );
globalStats.registerView(view);


app.get('/', (req, res) => {
    console.log("request made");

    
    res.status(200).send("success!");
})


app.listen(8080, () => console.log(`Example app listening on port 8080!`))