const { FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');
const fs = require('fs');
const ccpPath = path.resolve(__dirname, '..', '..', 'app_network', 'connection-logistics.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);
const walletPath = path.join(process.cwd(), 'CA','Logisticswallet');
const wallet = new FileSystemWallet(walletPath);
let controller = {
    addLogInfo: async function (key, logCopName, logDepartureTm, logArrivalTm, logDeparturePl, logDest, logMOT, tempAvg) {
        try {
            console.log(`Wallet path: ${walletPath}`);

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

            await contract.submitTransaction('addLogInfo', key,logCopName, logDepartureTm, logArrivalTm, logDeparturePl, logDest, logMOT, tempAvg );

            await gateway.disconnect();

            return '{ "status" : "1", "message": "添加成功"}';

        } catch (error) {
            return '{ "status" : "0", "message": '+ error + '}';
        }
    }
}

module.exports = controller