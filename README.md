#Installation process
CHANNEL_NAME=mainchannel TIMEOUT=50000 docker-compose -f docker-compose-cli.yaml up -d --force-recreate

#companies_chaincode
CHAINCODE=lending_chaincode
peer chaincode install -n $CHAINCODE -v 1.0 -p github.com/hyperledger/fabric/examples/chaincode/go/companies_chaincode
sleep 2
peer chaincode instantiate -o orderer.homelend.io:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CHAINCODE -v 1.0 -c '{"Args":["init"]}' -P "OR('NoorBankMSP.member','MashreqBankMSP.member')"
sleep 10

peer chaincode invoke -o orderer.homelend.io:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CHAINCODE -c '{"Args":["create", "{\"Hash\":\"hash_\",\"FlatNumber\":\"1\",\"HouseNumber\":\"1\",\"Street\":\"Main Street\",\"Owner\":\"Kanat Tulbasiyev\", \"Amount\":\"100\", \"Active\":true,\"Timestamp\":111}"]}'

peer chaincode invoke -o orderer.homelend.io:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CHAINCODE -c '{"Args":["create", "{\"Hash\":\"hash2_\",\"FlatNumber\":\"1\",\"HouseNumber\":\"2\",\"Street\":\"Main Street\",\"Owner\":\"Netanel Bitan\", \"Amount\":\"100\", \"Active\":true,\"Timestamp\":111}"]}'
sleep 10

peer chaincode query -C $CHANNEL_NAME -n $CHAINCODE -c '{"Args":["get","hash_"]}'

peer chaincode query -C $CHANNEL_NAME -n $CHAINCODE -c '{"Args":["query","{}"]}'
