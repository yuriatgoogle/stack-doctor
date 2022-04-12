from contextvars import Context
from multiprocessing import context
from flask import Flask
import requests
import time
import os

from opentelemetry import trace, baggage
from opentelemetry.trace import SpanContext, NonRecordingSpan
from opentelemetry.exporter.cloud_trace import CloudTraceSpanExporter
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

from opentelemetry.propagate import set_global_textmap, inject
from opentelemetry.propagators.cloud_trace_propagator import (
    CloudTraceFormatPropagator,
)
set_global_textmap(CloudTraceFormatPropagator())

def set_span_context_headers():
  headers = {}
  inject(headers)
  return headers

tracer_provider = TracerProvider()
cloud_trace_exporter = CloudTraceSpanExporter()
tracer_provider.add_span_processor(
    BatchSpanProcessor(cloud_trace_exporter)
)
trace.set_tracer_provider(tracer_provider)
tracer = trace.get_tracer(__name__)

app = Flask(__name__)
backend_addr = 'http://localhost:8081'

@app.route('/')
def index():
    with tracer.start_as_current_span("Root span") as parent:
        start = time.time()
        with tracer.start_as_current_span("Backend request") as child:
            headers = set_span_context_headers()
            r = requests.get(backend_addr, timeout=3, headers=headers)
            latency = time.time() - start
            return 'Response from backend: ' + str(r.status_code) + ' in ' + str(round(latency, 3) * 1000) + ' ms'

app.run(host='0.0.0.0', port=8080)
