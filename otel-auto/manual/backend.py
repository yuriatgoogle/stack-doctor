from flask import Flask, request
import requests
import time
import os
from random import randint
from time import sleep

from opentelemetry import trace, baggage
from opentelemetry.exporter.cloud_trace import CloudTraceSpanExporter
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

from opentelemetry.propagate import set_global_textmap, extract

from opentelemetry.trace import SpanContext, NonRecordingSpan

from opentelemetry.propagators.cloud_trace_propagator import (
    CloudTraceFormatPropagator,
)
set_global_textmap(CloudTraceFormatPropagator())

tracer_provider = TracerProvider()
cloud_trace_exporter = CloudTraceSpanExporter()
tracer_provider.add_span_processor(
    BatchSpanProcessor(cloud_trace_exporter)
)
trace.set_tracer_provider(tracer_provider)

tracer = trace.get_tracer(__name__)

app = Flask(__name__)

@app.route('/')
def index():
    # get incoming context from headers
    ctx = extract(request.headers)
    print(ctx)
    with tracer.start_as_current_span("backend process", context=ctx) as span:
        start = time.time()
        sleep(randint(1,1000)/1000)
        latency = time.time() - start
        return 'returned in ' + str(round(latency, 3) * 1000) + ' ms'
    
app.run(host='0.0.0.0', port=8081)
