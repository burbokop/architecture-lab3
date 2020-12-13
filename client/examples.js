// This file contains examples of scenarios implementation using

const channels = require('./interface/client');

const client = channels.Client('http://localhost:8080');

// Scenario 1: Display available virtual mashines.
client.listVirtualMashines()
    .then((list) => {
        console.log('=== Scenario 1 ===');
        console.log('Available virtual mashines:');
        list.forEach((c) => console.log(c));
    })
    .catch((e) => {
        console.log(`Problem listing available virtual mashines: ${e.message}`);
    });

// Scenario 2: Connect disc.
client.connectDisc(0, 0)
    .then((resp) => {
        console.log('=== Scenario 2 ===');
        console.log('Disc connection responce:', resp);
        //return client.listChannels()
        //    .then((list) => list.map((c) => c.name).join(', '))
        //    .then((str) => {
        //        console.log(`vm list: ${str}`);
        //    })
    })
    .catch((e) => {
        console.log(`Problem creating a new channel: ${e.message}`);
    });
