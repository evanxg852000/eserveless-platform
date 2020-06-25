const http = require('http');
const functions = require('./functions') 

//create a server object
http.createServer(functions.API)
    .listen(8000);