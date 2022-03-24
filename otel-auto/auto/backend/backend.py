from flask import Flask
import requests
import time
import os
from random import randint
from time import sleep

app = Flask(__name__)

@ app.route('/')
def index():
    start = time.time()
    sleep(randint(1,1000)/1000)
    latency = time.time() - start
    return 'returned in ' + str(round(latency, 3) * 1000) + ' ms'

app.run(host='0.0.0.0', port=8081)
