const express = require('express');
const app = express();

// opencensus setup
const {globalStats, MeasureUnit, AggregationType} = require('@opencensus/core');
const {StackdriverStatsExporter} = require('@opencensus/exporter-stackdriver');

// Stackdriver export interval is 60 seconds
const EXPORT_INTERVAL = 60;

// define the "golden signals" metrics and views
// request count measure
const REQUEST_COUNT = globalStats.createMeasureInt64(
    'request_count',
    MeasureUnit.UNIT,
    'Number of requests to the server'
);
// request count view
const request_count_view = globalStats.createView(
    'request_count_view',
    REQUEST_COUNT,
    AggregationType.LAST_VALUE
);
globalStats.registerView(request_count_view);

// error count measure
const ERROR_COUNT = globalStats.createMeasureInt64(
    'error_count',
    MeasureUnit.UNIT,
    'Number of failed requests to the server'
);
// error count view
const error_count_view = globalStats.createView(
    'error_count_view',
    ERROR_COUNT,
    AggregationType.LAST_VALUE
);
globalStats.registerView(error_count_view);

// response latency measure
const RESPONSE_LATENCY = globalStats.createMeasureInt64(
    'response_latency',
    MeasureUnit.MS,
    'The server response latency in milliseconds'
  );
// response latency view
const latency_view = globalStats.createView(
    'response_latency_view',
    RESPONSE_LATENCY,
    AggregationType.DISTRIBUTION,
    [],
    'Server response latency distribution',
    // Latency in buckets:
    [0, 100, 500, 1000]
  );
globalStats.registerView(latency_view);

// set up the Stackdriver exporter
const projectId = 'stack-doctor';

// GOOGLE_APPLICATION_CREDENTIALS are expected by a dependency of this code
// Not this code itself. Checking for existence here but not retaining (as not needed)
if (!projectId || !process.env.GOOGLE_APPLICATION_CREDENTIALS) {
  throw Error('Unable to proceed without a Project ID');
}
const exporter = new StackdriverStatsExporter({
  projectId: projectId,
  period: EXPORT_INTERVAL * 1000,
});
globalStats.registerExporter(exporter);


app.get('/', (req, res) => {
    console.log("request made");

    // record metric values
    globalStats.record([
        {
          measure: REQUEST_COUNT,
          value: 1,
        },
      ]);
    res.status(200).send("success!");
})


app.listen(8080, () => console.log(`Example app listening on port 8080!`))