const http = require('../common/http');

const Client = (baseUrl) => {
    const client = http.Client(baseUrl);
    return {
        listVirtualMashines: () => client.get('/vm_list'),
        connectDisc: (diskId, vmId) => client.post('/connect_disc', { diskId, vmId })
    }

};

module.exports = { Client };
