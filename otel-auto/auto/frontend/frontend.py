from flask import Flask
import requests
import time
import os

from opentelemetry import trace
from opentelemetry.exporter.jaeger.thrift import JaegerExporter
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

trace.set_tracer_provider(
TracerProvider(
        resource=Resource.create({SERVICE_NAME: "my-helloworld-service"})
    )
)
tracer = trace.get_tracer(__name__)

# create a JaegerExporter
jaeger_exporter = JaegerExporter(
    # configure agent
    agent_host_name='localhost',
    agent_port=6831,
)

# Create a BatchSpanProcessor and add the exporter to it
span_processor = BatchSpanProcessor(jaeger_exporter)

# add to the tracer
trace.get_tracer_provider().add_span_processor(span_processor)

backend_addr = os.getenv('BACKEND')

app = Flask(__name__)

@ app.route('/')
def index():
    start = time.time()
    r = requests.get(backend_addr)
    latency = time.time() - start
    return 'Response from backend: ' + str(r.status_code) + ' in ' + str(round(latency, 3) * 1000) + ' ms'

app.run(host='0.0.0.0', port=8080)
