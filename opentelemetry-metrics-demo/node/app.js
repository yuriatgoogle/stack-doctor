"use strict";

const { MeterRegistry } = require('@opentelemetry/metrics');
const { PrometheusExporter } = require('@opentelemetry/exporter-prometheus');

const meter = new MeterRegistry().getMeter('example-prometheus');

const exporter = new PrometheusExporter(
  {
    startServer: true,
    port: 8080
  },
  () => {
    console.log("prometheus scrape endpoint: http://localhost:8080/metrics");
  }
);

meter.addExporter(exporter);


const nonMonotonicGauge = meter.createGauge("non_monotonic_gauge", {
  monotonic: false,
  labelKeys: ["metricOrigin"],
  description: "Example of a non-monotonic gauge"
});

let metricValue = 0;
setInterval(() => {
  const labels = meter.labels({ metricOrigin: process.env.ENV});
  metricValue = Math.floor(Math.random() * 100);
  nonMonotonicGauge
    .bind(labels)
    .set(metricValue);
}, 1000);