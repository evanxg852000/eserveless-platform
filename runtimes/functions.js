//http function
function API(req, res) {
    res.writeHead(200);
    res.end('Hello, World!');
}


//cron function
function Ticker(envs) {
    console.log(envs); 
    console.log(`Ticker ticked at ${Date.now()}`);
}

module.exports = {
    API, 
    Ticker
}
