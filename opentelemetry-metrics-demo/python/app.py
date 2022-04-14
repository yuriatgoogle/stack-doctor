from flask import Flask
import time
import os
from random import randint
from time import sleep

from opentelemetry._metrics import get_meter_provider, set_meter_provider
from opentelemetry.exporter.otlp.proto.grpc._metric_exporter import OTLPMetricExporter
from opentelemetry.sdk._metrics import MeterProvider
from opentelemetry.sdk._metrics.export import PeriodicExportingMetricReader

meter = get_meter_provider().get_meter("getting-started")
metric_exporter = OTLPMetricExporter(
    # endpoint:="localhost:4317",
    # credentials=ChannelCredentials(credentials),
    # headers=(("metadata", "metadata")),
)
reader = PeriodicExportingMetricReader(metric_exporter)
provider = MeterProvider(metric_readers=[reader])
set_meter_provider(provider)

requests_counter = meter.create_counter("otel_total_requests")
errors_counter = meter.create_counter("otel_failed_requests")
# request_latency = meter.create_valuerecorder("otel_request_latency")

labels = {}
app = Flask(__name__)

@app.route('/')
def index():
    requests_counter.add(1)
    start = time.time()
    sleep(randint(1,1000)/1000)
    latency = time.time() - start
    if randint(1,100) > 95:
        # fail 5 % of the time
        errors_counter.add(1)
        # request_latency.record(latency)
        return 'Processing failed!', 500
    # request_latency.record(latency, labels=labels)
    return 'returned in ' + str(round(latency, 3) * 1000) + ' ms', 200
    
app.run(host='0.0.0.0', port=8080)