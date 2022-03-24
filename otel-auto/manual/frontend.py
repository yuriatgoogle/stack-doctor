from flask import Flask
import requests
import time
import os

from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import (
    BatchSpanProcessor,
    ConsoleSpanExporter,
)

provider = TracerProvider()
processor = BatchSpanProcessor(ConsoleSpanExporter())
provider.add_span_processor(processor)
trace.set_tracer_provider(provider)

tracer = trace.get_tracer(__name__)

backend_addr = os.getenv('BACKEND')

app = Flask(__name__)

@ app.route('/')
def index():
    with tracer.start_as_current_span("Root span") as parent:
        start = time.time()
        with tracer.start_as_current_span("Backend request") as child:
            r = requests.get(backend_addr, timeout=3)
            latency = time.time() - start
            return 'Response from backend: ' + str(r.status_code) + ' in ' + str(round(latency, 3) * 1000) + ' ms'

app.run(host='0.0.0.0', port=8080)
