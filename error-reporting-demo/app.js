const express = require('express');
const bunyan = require('bunyan');
const {LoggingBunyan} = require('@google-cloud/logging-bunyan');
const winston = require('winston');
const {LoggingWinston} = require('@google-cloud/logging-winston');
const {ErrorReporting} = require('@google-cloud/error-reporting');

const app = express();

const projectID = 'stack-doctor';
const serviceName = "node-error-reporting";

// *************** Bunyan logging setup *************
// Creates a Bunyan Stackdriver Logging client
const loggingBunyan = new LoggingBunyan();
// Create a Bunyan logger that streams to Stackdriver Logging
const bunyanLogger = bunyan.createLogger({
  name: serviceName,
  streams: [
    // Log to the console at 'info' and above
    {stream: process.stdout, level: 'info'},
    // And log to Stackdriver Logging, logging at 'info' and above
    loggingBunyan.stream('info'),
  ],
});

// ************* Winston logging setup *****************
const loggingWinston = new LoggingWinston();
// Create a Winston logger that streams to Stackdriver Logging
const winstonLogger = winston.createLogger({
  level: 'info',
  transports: [
    new winston.transports.Console(),
    // Add Stackdriver Logging
    loggingWinston,
  ],
});

//************** Stackdriver Error Reporting setup ******** */
const errors = new ErrorReporting(
  {
    projectId: projectID,
    reportMode: 'always',
    serviceContext: {
      service: serviceName,
      version: '1'
    }
  }
);

function getRandomInt(max) {
    return Math.floor(Math.random() * Math.floor(max));
  }

app.get('/bunyan-error', (req, res) => {
    bunyanLogger.error('Bunyan error logged');
    res.send('Bunyan error logged!');
}) 

app.get('/winston-error', (req, res) => {
    winstonLogger.error('Winston error logged');
    res.send('Winston error logged!');
}) 

app.get('/report-error', (req, res) => {
  res.send('Stackdriver error reported!');
  errors.report('Stackdriver error reported');
}) 

app.get('/log-exception', (req, res) => {
  res.send('exception');
  bunyanLogger.error(new Error('exception logged'));
})

app.get('/report-exception', (req, res) => {
  res.send('exception');
  errors.report(new Error('exception reported'));
})

app.get('/', (req, res) => {
    console.log('request made');
    res.send('request made');
})

app.use(errors.express);
app.listen(8080, () => console.log(`Example app listening on port 8080!`))