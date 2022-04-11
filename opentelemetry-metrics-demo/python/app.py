from flask import Flask, request
import requests
import time
import os
from random import randint
from time import sleep

from opentelemetry import metrics
from opentelemetry.exporter.cloud_monitoring import (
    CloudMonitoringMetricsExporter,
)
from opentelemetry.sdk.metrics import MeterProvider

metrics.set_meter_provider(MeterProvider())
meter = metrics.get_meter(__name__)
metrics.get_meter_provider().start_pipeline(
    meter, CloudMonitoringMetricsExporter(), 5
)

requests_counter = meter.create_counter(
    name="otel_total_requests",
    description="total requests",
    unit="1",
    value_type=int,
)

errors_counter = meter.create_counter(
    name="otel_failed_requests",
    description="failed requests",
    unit="1",
    value_type=int,
)

request_latency = meter.create_valuerecorder(
    name="otel_request_latency",
    description="request latency",
    unit='ms',
    value_type=float
)

labels = {}
app = Flask(__name__)

@app.route('/')
def index():
    requests_counter.add(1, labels=labels)
    start = time.time()
    sleep(randint(1,1000)/1000)
    latency = time.time() - start
    if randint(1,100) > 95:
        # fail 5 % of the time
        errors_counter.add(1, labels=labels)
        request_latency.record(latency)
        return 'Processing failed!', 500
    request_latency.record(latency, labels=labels)
    return 'returned in ' + str(round(latency, 3) * 1000) + ' ms', 200
    
app.run(host='0.0.0.0', port=8080)