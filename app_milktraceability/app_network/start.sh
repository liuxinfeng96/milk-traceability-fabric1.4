#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -ev

# don't rewrite paths for Windows Git Bash users

docker-compose -f docker-compose-cli.yaml -f docker-compose-ca.yaml -f docker-compose-couch.yaml up -d

docker ps -a

# wait for Hyperledger Fabric to start
# incase of errors when running later commands, issue export FABRIC_START_TIMEOUT=<larger number>
export FABRIC_START_TIMEOUT=5
#echo ${FABRIC_START_TIMEOUT}
sleep ${FABRIC_START_TIMEOUT}

# Create the channel
docker exec cli peer channel create -o orderer.example.com:7050 -c firstchannel -f ./channel-artifacts/channel.tx
docker exec cli peer channel join -b firstchannel.block
docker exec cli peer channel update -o orderer.example.com:7050 -c firstchannel -f ./channel-artifacts/SourceMSPanchors.tx

docker exec -e "CORE_PEER_LOCALMSPID=ProcessMSP" -e "CORE_PEER_ADDRESS=peer0.process.example.com:7051" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/process.example.com/users/Admin@process.example.com/msp" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/process.example.com/peers/peer0.process.example.com/tls/ca.crt" cli peer channel join -b firstchannel.block
docker exec -e "CORE_PEER_LOCALMSPID=ProcessMSP" -e "CORE_PEER_ADDRESS=peer0.process.example.com:7051" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/process.example.com/users/Admin@process.example.com/msp" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/process.example.com/peers/peer0.process.example.com/tls/ca.crt" cli peer channel update -o orderer.example.com:7050 -c firstchannel -f ./channel-artifacts/ProcessMSPanchors.tx

docker exec -e "CORE_PEER_LOCALMSPID=LogisticsMSP" -e "CORE_PEER_ADDRESS=peer0.logistics.example.com:8051" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/logistics.example.com/users/Admin@logistics.example.com/msp" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/tls/ca.crt" cli peer channel join -b firstchannel.block
docker exec -e "CORE_PEER_LOCALMSPID=LogisticsMSP" -e "CORE_PEER_ADDRESS=peer0.logistics.example.com:8051" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/logistics.example.com/users/Admin@logistics.example.com/msp" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/logistics.example.com/peers/peer0.logistics.example.com/tls/ca.crt" cli peer channel update -o orderer.example.com:7050 -c firstchannel -f ./channel-artifacts/LogisticsMSPanchors.tx

docker exec -e "CORE_PEER_LOCALMSPID=SalesMSP" -e "CORE_PEER_ADDRESS=peer0.sales.example.com:9051" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/sales.example.com/users/Admin@sales.example.com/msp" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/sales.example.com/peers/peer0.sales.example.com/tls/ca.crt" cli peer channel join -b firstchannel.block
docker exec -e "CORE_PEER_LOCALMSPID=SalesMSP" -e "CORE_PEER_ADDRESS=peer0.sales.example.com:8051" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/sales.example.com/users/Admin@sales.example.com/msp" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/sales.example.com/peers/peer0.sales.example.com/tls/ca.crt" cli peer channel update -o orderer.example.com:7050 -c firstchannel -f ./channel-artifacts/SalesMSPanchors.tx


