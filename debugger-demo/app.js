const express = require('express');
const app = express();

const projectID = 'stack-doctor';
const serviceName = "debugger-demo";
const serviceVersion = "1.0"

require('@google-cloud/debug-agent').start({
    projectId: projectID,
    keyFilename: './key.json',
    serviceContext: {
      service: serviceName,
      version: serviceVersion
    },
    allowExpressions: true
  });

function getRandomInt(max) {
    return Math.floor(Math.random() * Math.floor(max));
  }

app.get('/', (req, res) => {
    console.log('request made');
    randomInt = getRandomInt(10);
    if (randomInt > 5) {
        console.log("an error was encountered")
        return res.status(500).send("there was an error!")
    }
    res.send('request made');
})

app.listen(8080, () => console.log(`Example app listening on port 8080!`))