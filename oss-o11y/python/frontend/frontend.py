from flask import Flask
import requests
from prometheus_flask_exporter import PrometheusMetrics as prom_metrics

# OpenTelemetry Metrics
from opentelemetry import metrics as otel_metrics
from opentelemetry.exporter.prometheus import PrometheusMetricsExporter
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import ConsoleMetricsExporter
from opentelemetry.sdk.metrics.export.controller import PushController

# OpenTelemetry tracing
from opentelemetry import trace
from opentelemetry.exporter.cloud_trace import CloudTraceSpanExporter
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.trace import Link

app = Flask(__name__)

# Set up metrics
metrics = prom_metrics(app)
processor_mode = "stateful"
otel_metrics.set_meter_provider(MeterProvider())
meter = otel_metrics.get_meter(__name__, processor_mode == "stateful")
exporter = PrometheusMetricsExporter("frontend")
controller = PushController(meter, exporter, 5)

metric_labels = {}

# Request counter metric
otel_requests_counter = meter.create_counter(
    name="frontend_request_count_otel",
    description="Number of requests counted by Otel",
    unit="1",
    value_type=int)

# Set up tracing
tracer_provider = TracerProvider()
cloud_trace_exporter = CloudTraceSpanExporter()
tracer_provider.add_span_processor(
    BatchSpanProcessor(cloud_trace_exporter)
)
trace.set_tracer_provider(tracer_provider)
tracer = trace.get_tracer(__name__)


@ app.route('/')
@ metrics.counter('frontend_request_count_prom', 'Number of requests counted by Prom', labels={})
def index():
    with tracer.start_as_current_span("complete_transaction") as link_target:
        with tracer.start_as_current_span("link_span", links=[Link(link_target.context)]):
            # cound request
            otel_requests_counter.add(1, metric_labels)
            # add event for external request
            with tracer.start_as_current_span(
                "external_request",
                links=[Link(link_target.context, attributes={})],
            ):
                # make external request
                link_target.add_event(name="start external request")
                r = requests.get('https://www.google.com')
                link_target.add_event(name="complete external request")
                return 'Response from backend: ' + str(r.status_code)


app.run(host='0.0.0.0', port=8080)