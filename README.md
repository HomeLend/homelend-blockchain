# Executing process
CHANNEL_NAME=mainchannel TIMEOUT=50000 docker-compose -f docker-compose-cli.yaml up -d --force-recreate

# Installing the chaincode
peer chaincode install -n $DC -v v1 -p $CHAINCODE

peer chaincode instantiate -o orderer.homelend.io:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $DC -v v1 -c '{"Args":["init"]}' -P "OR('POCSellerMSP.admin', 'POCSellerMSP.member')"

# SELLING
peer chaincode invoke -o orderer.homelend.io:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $DC -v v1 -c '{"Args":["sell", "{\"Hash\":\"hash_\",\"FlatNumber\":\"1\",\"HouseNumber\":\"1\",\"Street\":\"Main Street\", \"Amount\":100, \"Active\":true,\"Timestamp\":111}"]}'

peer chaincode invoke -o orderer.homelend.io:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CHAINCODE -c '{"Args":["create", "{\"Hash\":\"hash2_\",\"FlatNumber\":\"1\",\"HouseNumber\":\"2\",\"Street\":\"Main Street\",\"Owner\":\"Netanel Bitan\", \"Amount\":\"100\", \"Active\":true,\"Timestamp\":111}"]}'
sleep 10

peer chaincode query -C $CHANNEL_NAME -n $DC -c '{"Args":["getUserHouses"]}'

peer chaincode query -C $CHANNEL_NAME -n $DC -c '{"Args":["query","{}"]}'
