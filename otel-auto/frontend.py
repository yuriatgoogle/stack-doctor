from flask import Flask
import requests
import time

app = Flask(__name__)

@ app.route('/')
def index():
    start = time.time()
    r = requests.get('https://www.google.com')
    latency = time.time() - start
    return 'Response from backend: ' + str(r.status_code) + ' in ' + str(round(latency, 3) * 1000) + ' ms'

app.run(host='0.0.0.0', port=8080)
