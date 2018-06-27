#!/bin/bash
set -e

echo
echo "#################################################################"
echo "#######        Generating cryptographic material       ##########"
echo "#################################################################"


PROJPATH=${PWD}
ARTIFACTSPATH=$PROJPATH/artifacts
CRYPTOPATH=$ARTIFACTSPATH/crypto
ORDERERS=$CRYPTOPATH/ordererOrganizations
PEERS=$CRYPTOPATH/peerOrganizations
CHANNELPATH=$PROJPATH/artifacts/channel
CHANNEL_NAME=salamatapp

# Remove any old artifacts
rm -rf $ARTIFACTSPATH

# Generate Crypto Material
$PROJPATH/cryptogen generate --config=$PROJPATH/crypto-config.yaml --output=$CRYPTOPATH
if [ "$?" -ne 0 ]; then
  echo "Failed to generate crypto material..."
  exit 1
fi

# Create Channel Config Folder
mkdir $CHANNELPATH

echo
echo "##########################################################"
echo "#########  Generating Orderer Genesis block ##############"
echo "##########################################################"
$PROJPATH/configtxgen -profile OneOrgGenesis -outputBlock $CHANNELPATH/genesis.block

echo
echo "#################################################################"
echo "### Generating channel configuration transaction 'channel.tx' ###"
echo "#################################################################"
$PROJPATH/configtxgen -profile OneOrgChannel -outputCreateChannelTx $CHANNELPATH/channel.tx -channelID $CHANNEL_NAME
cp $CHANNELPATH/channel.tx $PROJPATH/web

echo
echo "#################################################################"
echo "####### Generating anchor peer update for ThanksOrg ##########"
echo "#################################################################"
$PROJPATH/configtxgen -profile OneOrgChannel -outputAnchorPeersUpdate $CHANNELPATH/Org1MSPAnchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP