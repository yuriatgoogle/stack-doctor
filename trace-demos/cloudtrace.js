require('@google-cloud/trace-agent').start();

const express = require('express');
const got = require('got');

const app = express();
const BACKEND_URL = 'https://www.google.com/';

// This incoming HTTP request should be captured by Trace
app.get('/', async (req, res) => {
  // This outgoing HTTP request should be captured by Trace
  try {
    const backend_response = await got(BACKEND_URL);
    res.status(200).send(backend_response.statusCode).end();
  } catch (err) {
    console.error(err);
    res.status(500).end();
  }
});

// Start the server
const PORT = process.env.PORT || 8080;
app.listen(PORT, () => {
  console.log(`App listening on port ${PORT}`);
  console.log('Press Ctrl+C to quit.');
});
