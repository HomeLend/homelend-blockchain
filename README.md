# Executing process
CHANNEL_NAME=mainchannel TIMEOUT=50000 docker-compose -f docker-compose-cli.yaml up -d --force-recreate

# Installing the chaincode
peer chaincode install -n $DC -v v1 -p $CHAINCODE

# Instiantiating
peer chaincode instantiate -o orderer.homelend.io:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $DC -v v1 -c '{"Args":["init"]}' -P "OR('POCSellerMSP.admin', 'POCSellerMSP.member')"

# SELLING
peer chaincode invoke -o orderer.homelend.io:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $DC -v v1 -c '{"Args":["sell", "{\"Hash\":\"hash_\",\"FlatNumber\":\"1\",\"HouseNumber\":\"1\",\"Street\":\"Main Street\", \"Amount\":100, \"Active\":true,\"Timestamp\":111}"]}'

# REGISTER AS BANK
peer chaincode invoke -o orderer.homelend.io:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $DC -v v1 -c '{"Args":["registerAsBank", "{\"Name\":\"NYBank\",\"LicenseNumber\":\"1\",\"Address\":\"Brooklyn 55\",\"TotalSupply\":10000000, \"Timestamp\":111}"]}'

# REGISTER AS SELLER
peer chaincode invoke -o orderer.homelend.io:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $DC -v v1 -c '{"Args":["registerAsSeller", "{\"Name\":\"Seller#1\",\"Firstname\":\"John\",\"Lastname\":\"Smith\", \"Timestamp\":111}"]}'

# GET USER TOKENS
peer chaincode query -C $CHANNEL_NAME -n $DC -c '{"Args":["getUserHouses"]}'

# GET ALL CHAINCODE RESULTS
peer chaincode query -C $CHANNEL_NAME -n $DC -c '{"Args":["query","{}"]}'
