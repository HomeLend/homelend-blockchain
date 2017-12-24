#!/bin/bash
export FABRIC_CFG_PATH=$PWD
export CHANNEL_NAME=mainchannel

rm -rf channel-artifacts crypto-config
mkdir channel-artifacts
./bin/cryptogen generate --config=./crypto-config.yaml
./bin/configtxgen -profile MainOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
./bin/configtxgen -profile MainChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
./bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCBankMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCBankMSP
./bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCInsuranceMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCInsuranceMSP
./bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCAppraiserMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCAppraiserMSP
./bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCGovernmentMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCGovernmentMSP
./bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCBuyerMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCBuyerMSP
./bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCSellerMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCSellerMSP
./bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCHomelendMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCHomelendMSP
