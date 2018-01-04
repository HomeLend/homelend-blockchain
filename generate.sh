#!/bin/bash
export FABRIC_CFG_PATH=$PWD
export CHANNEL_NAME=mainchannel

rm -rf channel-artifacts crypto-config
mkdir channel-artifacts
./new_bin/cryptogen generate --config=./crypto-config.yaml
./new_bin/configtxgen -profile MainOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
./new_bin/configtxgen -profile MainChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
./new_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCBankMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCBankMSP
./new_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCInsuranceMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCInsuranceMSP
./new_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCAppraiserMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCAppraiserMSP
./new_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCGovernmentMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCGovernmentMSP
./new_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCBuyerMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCBuyerMSP
./new_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCSellerMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCSellerMSP
./new_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCCreditRatingAgencyMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCCreditRatingAgencyMSP
./new_bin/configtxgen -profile MainChannel -outputAnchorPeersUpdate ./channel-artifacts/POCHomelendMSPanchors.tx -channelID $CHANNEL_NAME -asOrg POCHomelendMSP
