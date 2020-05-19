const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', '..', 'app_network', 'connection-source.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);
const walletPath = path.join(process.cwd(), 'CA','Sourcewallet');
const wallet = new FileSystemWallet(walletPath);
let controller = {
    addSourceInfo: async function (key, grassState, cowState, milkState) {
        try {
            const userExists = await wallet.exists('user1');
            if (!userExists) {
                console.log('An identity for the user "user1" does not exist in the wallet');
                console.log('Run the registerUser.js application before retrying');
                return;
            }

            const gateway = new Gateway();
            await gateway.connect(ccp, {
                wallet,
                identity: 'user1',
                discovery: {
                    enabled: false
                }
            });

            const network = await gateway.getNetwork('firstchannel');
            const contract = network.getContract('milkchaincode');

            await contract.submitTransaction('addSourceInfo', key, grassState, cowState, milkState);

            await gateway.disconnect();

            return '{ "status" : "1", "message": "添加成功"}';

        } catch (error) {
            return '{ "status" : "0", "message": '+ error + '}';
        }
    }
}
module.exports = controller