# Executing process
CHANNEL_NAME=mainchannel TIMEOUT=50000 docker-compose -f docker-compose-cli.yaml up -d --force-recreate

# Installing the chaincode
peer chaincode install -n $DC -v v1 -p $CHAINCODE

# Installing DC
DC=lending_chaincode

# Instiantiating
peer chaincode instantiate -o orderer.homelend.io:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $DC -v v1 -c '{"Args":["init"]}' -P "OR('POCSellerMSP.admin', 'POCSellerMSP.member')"

# ADVERTISE
peer chaincode invoke -o orderer.homelend.io:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $DC -v v1 -c '{"Args":["advertise", "{\"Hash\":\"hash_\",\"Address\":\"Shahal 5\", \"SellingPrice\":100000, \"Timestamp\":111}"]}'

# BUY
peer chaincode invoke -o orderer.homelend.io:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $DC -v v1 -c '{"Args":["buy", "{\"Hash\":\"hash_\",\"PropertyHash\":\"hash_\",\"Buyer\":\"eDUwOTo6Q049QWRtaW5AcG9jYmFuay5ob21lbGVuZC5pbyxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5wb2NiYW5rLmhvbWVsZW5kLmlvLE89cG9jYmFuay5ob21lbGVuZC5pbyxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT\",\"Salary\":1000, \"LoanAmount\":100, \"Status\": \"PENDING\", \"Active\":true,\"Timestamp\":111}"]}'

# GET USER TOKENS
peer chaincode query -C $CHANNEL_NAME -n $DC -c '{"Args":["getProperties"]}'

# GET ALL CHAINCODE RESULTS
peer chaincode query -C $CHANNEL_NAME -n $DC -c '{"Args":["query","{}"]}'
