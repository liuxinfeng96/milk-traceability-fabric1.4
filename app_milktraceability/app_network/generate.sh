#!/bin/sh
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
export PATH=$GOPATH/src/github.com/hyperledger/fabric/build/bin:${PWD}/./bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}
CHANNEL_NAME=firstchannel

# remove previous crypto material and config transactions
rm -fr channel-artifacts/*
rm -fr crypto-config/*

# generate crypto material
cryptogen generate --config=./crypto-config.yaml
if [ "$?" -ne 0 ]; then
  echo "Failed to generate crypto material..."
  exit 1
fi

# generate genesis block for orderer
configtxgen -profile FourOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
if [ "$?" -ne 0 ]; then
  echo "Failed to generate orderer genesis block..."
  exit 1
fi

# generate channel configuration transaction
configtxgen -profile FourOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
if [ "$?" -ne 0 ]; then
  echo "Failed to generate channel configuration transaction..."
  exit 1
fi

# generate anchor peer transaction
configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/SourceMSPanchors.tx -channelID $CHANNEL_NAME -asOrg SourceMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for SourceMSP..."
  exit 1
fi
configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/ProcessMSPanchors.tx -channelID $CHANNEL_NAME -asOrg ProcessMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for ProcessMSP..."
  exit 1
fi
configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/LogisticsMSPanchors.tx -channelID $CHANNEL_NAME -asOrg LogisticsMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for LogisticsMSP..."
  exit 1
fi
configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/SalesMSPanchors.tx -channelID $CHANNEL_NAME -asOrg SalesMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for SalesMSP..."
  exit 1
fi
