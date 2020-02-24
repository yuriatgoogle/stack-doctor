const express = require('express');
const app = express();
const sleep = require('sleep');

// set up Stackdriver logging
const bunyan = require('bunyan');

// Imports the Google Cloud client library for Bunyan
const {LoggingBunyan} = require('@google-cloud/logging-bunyan');

// Creates a Bunyan Stackdriver Logging client
const loggingBunyan = new LoggingBunyan();

// Create a Bunyan logger that streams to Stackdriver Logging
// Logs will be written to: "projects/YOUR_PROJECT_ID/logs/bunyan_log"
const logger = bunyan.createLogger({
  name: 'node-example',
  streams: [
    // Log to the console at 'info' and above
    {stream: process.stdout, level: 'info'},
    // And log to Stackdriver Logging, logging at 'info' and above
    loggingBunyan.stream('info'),
  ],
});

app.get('/', (req, res) => {
    console.log("request made");
    // add random delay
    randomDelay = Math.floor(Math.random() * 1000);
    sleep.msleep(randomDelay);
    logger.info("slept for " + randomDelay + " ms");
    
    // log error 5% of the time
    if ((Math.floor(Math.random() * 100)) <= 5) {
        // log error
        logger.error("failure!");
        // return error code
        res.status(500).send("error!")
    } 
    logger.info("success!");
    res.status(200).send("success!");
})

app.listen(8080, () => console.log(`Example app listening on port 8080!`))