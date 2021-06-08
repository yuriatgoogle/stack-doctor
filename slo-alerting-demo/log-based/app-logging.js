const express = require('express');
const app = express();
// get error rate


app.get('/', (req, res) => {
    console.log("request made");
    // log error some of the the time
    const randomValue = Math.floor(Math.random() * 100);
    const ERROR_RATE = process.env.ERROR_RATE || 1;
    if (randomValue <= ERROR_RATE) {
        // log error
        console.error("Status log - 500");
        // return error code
        res.status(500).send("error!")
    } else {
      console.log("Status log - 200");
      res.status(200).send("success!");
    }
})


app.listen(8080, () => console.log(`Example app listening on port 8080!`))
