const functions = require('./functions') 

const envs = Object
    .keys(process.env)
    .filter(name => name.startsWith('ESERVELESS_'))
    .reduce((envs, name) => {
            envs[name] = process.env[name];
            return envs
        }, {});
    
functions.Ticker(envs)
