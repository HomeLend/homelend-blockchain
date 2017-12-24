#!/usr/bin/env bash
#/bin/bash
rm -rf /tmp/fabric-client-*/
echo " ==================== Cleaning up containers ==================== "
docker ps -a | grep 'dev-peer0\|hyperledger' | awk '{print $1}' | xargs docker rm -f
echo " ==================== Cleaning up chaincode images ==================== "
docker images | grep 'dev-peer0' | awk '{print $1}' | xargs docker rmi -f
CHANNEL_NAME=mainchannel TIMEOUT=50000 docker-compose -f docker-compose-cli.yaml up -d --force-recreate
docker logs -f cli