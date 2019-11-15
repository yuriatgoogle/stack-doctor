# What's the best way to log errors (in Node.js)?

I wanted to address another one in the mostly-in-my-head series of questions with the running title of "things people often ask me".  Today's episode in the series is all about logging errors to Stackdriver.  Specifically, I've found that folks are somewhat confused about the multiple options they have for error logging and even more so when they want to understand how to log and track exceptions.  My opinion is that this is in part caused by Stackdriver providing multiple features that enable this - Error Reporting and Logging.  This is further confusing because Error Reporting is in a way a subset of Logging.  As such, I set out to explore exactly what happens when I tried to log both errors and exceptions using Logging and Error Reporting in a sample Node.js [app](https://github.com/yuriatgoogle/stack-doctor/tree/master/error-reporting-demo).  Let's see what I found!

# Logging Errors

I think that the confusion folks face starts with the fact that Stackdriver actually supports three different [options](https://cloud.google.com/logging/docs/setup/nodejs) for logging in Node.js - Bunyan, Winston, and the API client library.  I wanted to see how the first two treat error logs. At this point, I do not believe we recommend using the client library directly (in the same way that we recommend using OpenCensus for metric telemetry, rather than calling the Monitoring API directly).  

## Logging with Bunyan

The [documentation](https://cloud.google.com/logging/docs/setup/nodejs) is pretty straightforward - setting up Bunyan logging in my app was very easy.

```javascript
// *************** Bunyan logging setup *************
// Creates a Bunyan Stackdriver Logging client
const loggingBunyan = new LoggingBunyan();
// Create a Bunyan logger that streams to Stackdriver Logging
const bunyanLogger = bunyan.createLogger({
  name: serviceName, // this is set by an env var or as a parameter
  streams: [
    // Log to the console at 'info' and above
    {stream: process.stdout, level: 'info'},
    // And log to Stackdriver Logging, logging at 'info' and above
    loggingBunyan.stream('info'),
  ],
});
```

From there, logging an error message is as simple as:

```javascript
app.get('/bunyan-error', (req, res) => {
    bunyanLogger.error('Bunyan error logged');
    res.send('Bunyan error logged!');
})
```

When I ran my app, I saw this  logging output in the console:

    {"name":"node-error-reporting","hostname":"ygrinshteyn-macbookpro1.roam.corp.google.com","pid":5539,"level":50,"msg":"Bunyan error logged","time":"2019-11-15T17:19:58.001Z","v":0}

And this in Stackdriver Logging:

![image](https://github.com/yuriatgoogle/stack-doctor/blob/master/error-reporting-demo/images/1.png?raw=true)

Note that the log entry is created against the "global" resource because the log entry is being sent from my local machine not running on GCP, and the logName is bunyan_log.  The output is nicely structured, and the severity is set to ERROR.  

## Logging with Winston

I again followed the documentation to set up the Winston client:

```javascript
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
```

Then I logged an error:

```javascript
app.get('/winston-error', (req, res) => {
    winstonLogger.error('Winston error logged');
    res.send('Winston error logged!');
}) 
```

This time, the console output was much more concise:

    {"message":"Winston error logged","level":"error"}

Here's what I saw in the Logs Viewer:

![image](https://github.com/yuriatgoogle/stack-doctor/blob/master/error-reporting-demo/images/2.png?raw=true)

The severity was again set properly, but there's a lot less information in this entry.  For example, my hostname is not logged.  This may be a good choice for folks looking to reduce the amount of data that is logged while still retaining enough information to be useful.  

## Error Reporting

At this point, I had a good understanding of how logging errors works.  I next wanted to investigate whether using Error Reporting for this purpose would provide additional value.  First, I set up Error Reporting in the app:

```javascript
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
```

I then sent an error using the client:

```javascript
app.get('/report-error', (req, res) => {
  res.send('Stackdriver error reported!');
  errors.report('Stackdriver error reported');
}) 
```

This time, there was no output in the console AND nothing was logged to Stackdriver Logging.  I went to Error Reporting to find my error:

![image](https://github.com/yuriatgoogle/stack-doctor/blob/master/error-reporting-demo/images/4.png?raw=true)

When I clicked on the error, I was able to get a lot of detail:

![image](https://github.com/yuriatgoogle/stack-doctor/blob/master/error-reporting-demo/images/5.png?raw=true)

This is great because I can see when the error started happening, I get a histogram if and when it continues to happen, and I get a full stack trace showing me exactly where in my code the error is generated - this is all incredibly valuable information that I don't get from simply logging with the ERROR severity.  

The tradeoff here is that this message never makes it to Stackdriver Logging.  This means that I can't use errors reported through Error Reporting to, for example, create [log based metrics](https://dev.to/yurigrinshteyn/can-you-alert-on-logs-in-stackdriver-1lp8), which may make for a great SLI and/or alerting policy condition.

# Logging Exceptions

Next, I wanted to investigate what would happen if my app were to throw an exception and log it - how would it show up?  I used Bunyan to log an exception:

```javascript
app.get('/log-exception', (req, res) => {
  res.send('exception');
  bunyanLogger.error(new Error('exception logged'));
})
```

The console output contained the entire exception:

    {"name":"node-error-reporting","hostname":"<hostname>","pid":5539,"level":50,"err":{"message":"exception logged","name":"Error","stack":"Error: exception logged\n    at app.get (/Users/ygrinshteyn/src/error-reporting-demo/app.js:72:22)\n    at Layer.handle [as handle_request] (/Users/ygrinshteyn/src/error-reporting-demo/node_modules/express/lib/router/layer.js:95:5)\n    at next (/Users/ygrinshteyn/src/error-reporting-demo/node_modules/express/lib/router/route.js:137:13)\n    at Route.dispatch (/Users/ygrinshteyn/src/error-reporting-demo/node_modules/express/lib/router/route.js:112:3)\n    at Layer.handle [as handle_request] (/Users/ygrinshteyn/src/error-reporting-demo/node_modules/express/lib/router/layer.js:95:5)\n    at /Users/ygrinshteyn/src/error-reporting-demo/node_modules/express/lib/router/index.js:281:22\n    at Function.process_params (/Users/ygrinshteyn/src/error-reporting-demo/node_modules/express/lib/router/index.js:335:12)\n    at next (/Users/ygrinshteyn/src/error-reporting-demo/node_modules/express/lib/router/index.js:275:10)\n    at expressInit (/Users/ygrinshteyn/src/error-reporting-demo/node_modules/express/lib/middleware/init.js:40:5)\n    at Layer.handle [as handle_request] (/Users/ygrinshteyn/src/error-reporting-demo/node_modules/express/lib/router/layer.js:95:5)"},"msg":"exception logged","time":"2019-11-15T17:47:50.981Z","v":0}

The logging entry looked like this:

![image](https://github.com/yuriatgoogle/stack-doctor/blob/master/error-reporting-demo/images/6.png?raw=true)

And the jsonPayload contained the exception:

![image](https://github.com/yuriatgoogle/stack-doctor/blob/master/error-reporting-demo/images/7.png?raw=true)

This is definitely a lot of useful data.  I next wanted to see if Error Reporting would work as [advertised](https://cloud.google.com/error-reporting/docs/setup/nodejs) and identify this exception in the log as an error.  After carefully reviewing the documentation, I realized that this functionality works specifically on GCE, GKE, App Engine, and Cloud Functions, whereas I was just running my code on my local desktop. I tried running the code in Cloud Shell and immediately got a new entry in Error Reporting:

![image](https://github.com/yuriatgoogle/stack-doctor/blob/master/error-reporting-demo/images/8.png?raw=true)

The full stack trace of the exception is available in the detail view:

![image](https://github.com/yuriatgoogle/stack-doctor/blob/master/error-reporting-demo/images/9.png?raw=true)

So, logging an exception gives me the best of **both** worlds - I get a logging entry that I can use for things like log based metrics, and I get an entry in Error Reporting that I can use for analysis and tracking.

# Reporting Exceptions

I next wanted to see what would happen if I used Error Reporting to report the same exception.  

```javascript
app.get('/report-exception', (req, res) => {
  res.send('exception');
  errors.report(new Error('exception reported'));
})
```

Once again, there was no console output.  My error was immediately visible in Error Reporting:

![image](https://github.com/yuriatgoogle/stack-doctor/blob/master/error-reporting-demo/images/10.png?raw=true)

And somewhat to my surprise, I was able to see an entry in Logging, as well:

![image](https://github.com/yuriatgoogle/stack-doctor/blob/master/error-reporting-demo/images/11.png?raw=true)

As it turns out, exceptions are recorded in both Error Reporting AND Logging - no matter which of the two you use to send them.

# So, what now?

Here's what I've learned from this exercise:

1.  Bunyan logging is more verbose than Winston, which could be a consideration if cost is an issue.
1.  **Exceptions** can be sent to Stackdriver through Logging or Error Reporting - they will then be available in both.
1.  Using Error Reporting to report** non-exception** errors adds a lot of value for developers, but gives up value for SREs or ops folks who need to use logs for metrics or SLIs.

Thanks for joining me - come back soon for more!