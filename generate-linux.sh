#!/bin/bash
export FABRIC_CFG_PATH=$PWD
export CHANNEL_NAME=mainchannel

rm -rf channel-artifacts crypto-config
mkdir channel-artifacts
./linux_bin/cryptogen generate --config=./crypto-config.yaml
./linux_bin/configtxgen -profile MainOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
./linux_bin/configtxgen -profile MainChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
./linux_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCBankMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCBankMSP
./linux_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCInsuranceMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCInsuranceMSP
./linux_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCAppraiserMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCAppraiserMSP
./linux_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCGovernmentMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCGovernmentMSP
./linux_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCBuyerMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCBuyerMSP
./linux_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCSellerMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCSellerMSP
./linux_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCHomelendMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCHomelendMSP
