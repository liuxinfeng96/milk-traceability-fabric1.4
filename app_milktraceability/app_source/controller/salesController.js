const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', '..', 'app_network', 'connection-sales.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);
const walletPath = path.join(process.cwd(), 'CA','Saleswallet');
let controller = {
    getHistoryInfo: async function (key) {
        try {
            console.log(`Wallet path: ${walletPath}`);
            const wallet = new FileSystemWallet(walletPath);
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
            if (key == null) {
                return '请输入正确的产品ID'
            } else {
                const result = await contract.evaluateTransaction('getHistoryInfo', key);
                console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
                return result.toString();
            }
        } catch (error) {
            console.error(`Failed to evaluate transaction: ${error}`);
            process.exit(1);
        }
    },
    queryMilk: async function (key) {
        try {
            console.log(`Wallet path: ${walletPath}`);
            const wallet = new FileSystemWallet(walletPath);
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
            if (key == null) {
                return '请输入正确的产品ID'
            } else {
                const result = await contract.evaluateTransaction('queryMilk', key);
                console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
                return result.toString();
            }
        } catch (error) {
            console.error(`Failed to evaluate transaction: ${error}`);
            process.exit(1);
        }
    },
    queryAllMilks: async function () {
        try {
            console.log(`Wallet path: ${walletPath}`);
            const wallet = new FileSystemWallet(walletPath);
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
            const result = await contract.evaluateTransaction('queryAllMilks');
            console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
            return result.toString();
        } catch (error) {
            console.error(`Failed to evaluate transaction: ${error}`);
            process.exit(1);
        }
    }
}
module.exports = controller