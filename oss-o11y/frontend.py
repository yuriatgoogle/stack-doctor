from flask import Flask
from prometheus_flask_exporter import PrometheusMetrics

app = Flask(__name__)
metrics = PrometheusMetrics(app)


@app.route('/')
@metrics.counter('frontend_request_count', 'Number of invocations', labels={})
def index():
    return 'Hello world!'


app.run(host='0.0.0.0', port=8080)
