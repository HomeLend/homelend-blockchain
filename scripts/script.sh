#!/bin/bash

CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="mainchannel"}
: ${TIMEOUT:="60"}
COUNTER=1
MAX_RETRY=5

declare -a CHAINCODES=("creditscore_chaincode" "government_chaincode" "lending_chaincode")

echo "Channel name : "${CHANNEL_NAME}

# verify the result of the end-to-end test
verifyResult () {
	if [ $1 -ne 0 ] ; then
		echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
    echo "========= ERROR !!! FAILED to execute ==========="
		echo
   		exit 1
	fi
}

setGlobals () {
    if [ $1 -eq 0 ] ; then
        CORE_PEER_LOCALMSPID="POCBankMSP"
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pocbank.homelend.io/peers/peer0.pocbank.homelend.io/tls/ca.crt
        CORE_PEER_MSPCONFIGPATH=/var/hyperledger/cli/crypto/peerOrganizations/pocbank.homelend.io/users/Admin@pocbank.homelend.io/msp
        CORE_PEER_ADDRESS=peer0.pocbank.homelend.io:7051
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pocbank.homelend.io/peers/peer0.pocbank.homelend.io/tls/ca.crt
    elif [ $1 -eq 1 ]; then
        CORE_PEER_LOCALMSPID="POCInsuranceMSP"
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pocinsurance.homelend.io/peers/peer0.pocinsurance.homelend.io/tls/ca.crt
        CORE_PEER_MSPCONFIGPATH=/var/hyperledger/cli/crypto/peerOrganizations/pocinsurance.homelend.io/users/Admin@pocinsurance.homelend.io/msp
        CORE_PEER_ADDRESS=peer0.pocinsurance.homelend.io:7051
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pocinsurance.homelend.io/peers/peer0.pocinsurance.homelend.io/tls/ca.crt
    elif [ $1 -eq 2 ]; then
        CORE_PEER_LOCALMSPID="POCAppraiserMSP"
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pocappraiser.homelend.io/peers/peer0.pocappraiser.homelend.io/tls/ca.crt
        CORE_PEER_MSPCONFIGPATH=/var/hyperledger/cli/crypto/peerOrganizations/pocappraiser.homelend.io/users/Admin@pocappraiser.homelend.io/msp
        CORE_PEER_ADDRESS=peer0.pocappraiser.homelend.io:7051
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pocappraiser.homelend.io/peers/peer0.pocappraiser.homelend.io/tls/ca.crt
    elif [ $1 -eq 3 ]; then
        CORE_PEER_LOCALMSPID="POCGovernmentMSP"
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pocgovernment.homelend.io/peers/peer0.pocgovernment.homelend.io/tls/ca.crt
        CORE_PEER_MSPCONFIGPATH=/var/hyperledger/cli/crypto/peerOrganizations/pocgovernment.homelend.io/users/Admin@pocgovernment.homelend.io/msp
        CORE_PEER_ADDRESS=peer0.pocgovernment.homelend.io:7051
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pocgovernment.homelend.io/peers/peer0.pocgovernment.homelend.io/tls/ca.crt
    elif [ $1 -eq 4 ]; then
        CORE_PEER_LOCALMSPID="POCBuyerMSP"
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pocbuyer.homelend.io/peers/peer0.pocbuyer.homelend.io/tls/ca.crt
        CORE_PEER_MSPCONFIGPATH=/var/hyperledger/cli/crypto/peerOrganizations/pocbuyer.homelend.io/users/Admin@pocbuyer.homelend.io/msp
        CORE_PEER_ADDRESS=peer0.pocbuyer.homelend.io:7051
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pocbuyer.homelend.io/peers/peer0.pocbuyer.homelend.io/tls/ca.crt
    elif [ $1 -eq 5 ]; then
        CORE_PEER_LOCALMSPID="POCSellerMSP"
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pocseller.homelend.io/peers/peer0.pocseller.homelend.io/tls/ca.crt
        CORE_PEER_MSPCONFIGPATH=/var/hyperledger/cli/crypto/peerOrganizations/pocseller.homelend.io/users/Admin@pocseller.homelend.io/msp
        CORE_PEER_ADDRESS=peer0.pocseller.homelend.io:7051
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pocseller.homelend.io/peers/peer0.pocseller.homelend.io/tls/ca.crt
    elif [ $1 -eq 6 ]; then
        CORE_PEER_LOCALMSPID="POCHomelendMSP"
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pochomelend.homelend.io/peers/peer0.pochomelend.homelend.io/tls/ca.crt
        CORE_PEER_MSPCONFIGPATH=/var/hyperledger/cli/crypto/peerOrganizations/pochomelend.homelend.io/users/Admin@pochomelend.homelend.io/msp
        CORE_PEER_ADDRESS=peer0.pochomelend.homelend.io:7051
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/pochomelend.homelend.io/peers/peer0.pochomelend.homelend.io/tls/ca.crt
    elif [ $1 -eq 7 ]; then
        CORE_PEER_LOCALMSPID="POCCreditRatingAgencyMSP"
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/poccreditratingagency.homelend.io/peers/peer0.poccreditratingagency.homelend.io/tls/ca.crt
        CORE_PEER_MSPCONFIGPATH=/var/hyperledger/cli/crypto/peerOrganizations/poccreditratingagency.homelend.io/users/Admin@poccreditratingagency.homelend.io/msp
        CORE_PEER_ADDRESS=peer0.poccreditratingagency.homelend.io:7051
        CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/cli/crypto/peerOrganizations/poccreditratingagency.homelend.io/peers/peer0.poccreditratingagency.homelend.io/tls/ca.crt
    fi
    echo " ==================== GLOBALS =================="
    env |grep CORE
}

createChannel() {
	setGlobals 0

    peer channel create -o orderer.homelend.io:7050 -c ${CHANNEL_NAME} --tls $CORE_PEER_TLS_ENABLED  -f ./channel-artifacts/channel.tx --cafile $ORDERER_CA>&log.txt

	res=$?
	cat log.txt
	verifyResult $res "Channel creation failed"
	echo "===================== Channel \"$CHANNEL_NAME\" is created successfully ===================== "
	echo
}

updateAnchorPeers() {
    PEER=$1
    setGlobals $PEER

    peer channel update -o orderer.homelend.io:7050 -c ${CHANNEL_NAME} --tls $CORE_PEER_TLS_ENABLED  --cafile $ORDERER_CA -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx >&log.txt

    res=$?
    cat log.txt
    verifyResult $res "Anchor peer update failed"
    echo "===================== Anchor peers for org \"$CORE_PEER_LOCALMSPID\" on \"$CHANNEL_NAME\" is updated successfully ===================== "
echo
}

## Sometimes Join takes time hence RETRY atleast for 5 times
joinWithRetry () {
	peer channel join -b $CHANNEL_NAME.block >&log.txt
	res=$?
	cat log.txt
	if [ $res -ne 0 -a $COUNTER -lt $MAX_RETRY ]; then
		COUNTER=` expr $COUNTER + 1`
		echo "PEER$1 failed to join the channel, Retry after 2 seconds"
		sleep 5
		joinWithRetry $1
	else
		COUNTER=1
	fi
  verifyResult $res "After $MAX_RETRY attempts, PEER$ch has failed to Join the Channel"
}

joinChannel () {
	for ch in 0 1 2 3 4 5 6 7; do
		setGlobals $ch
		joinWithRetry $ch
		echo "===================== PEER$ch joined on the channel \"$CHANNEL_NAME\" ===================== "
		sleep 2
		echo
	done
}

installChaincode () {
	PEER=$1
	setGlobals $PEER
	for chaincode in ${CHAINCODES[*]}; do
        peer chaincode install -n $chaincode -v v1 --cafile $ORDERER_CA -p $chaincode >&log.txt
        res=$?
        cat log.txt
        verifyResult $res "Chaincode installation on remote peer PEER$PEER has Failed"
        echo "===================== $chaincode is installed on remote peer PEER$PEER ===================== "
        echo
    done
}

instantiateChaincode () {
	PEER=$1
	setGlobals $PEER
	# while 'peer chaincode' command can get the orderer endpoint from the peer (if join was successful),
	# lets supply it directly as we know it using the "-o" option
    
    for chaincode in ${CHAINCODES[*]}; do
        echo "===================== $chaincode Instantiation on PEER$PEER on channel '$CHANNEL_NAME' is in process ===================== "
        peer chaincode instantiate -o orderer.homelend.io:7050 --cafile $ORDERER_CA --tls $CORE_PEER_TLS_ENABLED -C $CHANNEL_NAME -n $chaincode -v v1 -c '{"Args":["init"]}' -P "OR	('POCBankMSP.member','POCSellerMSP.member', 'POCBuyerMSP.member', 'POCAppraiserMSP.member','POCCreditRatingAgencyMSP.member', 'POCInsuranceMSP.member','POCGovernmentMSP.member')"

        res=$?
        cat log.txt
        verifyResult $res "Chaincode instantiation on PEER$PEER on channel '$CHANNEL_NAME' failed"
        echo
    done
}

## Create channel
echo "Creating channel..."
createChannel

## Join all the peers to the channel
echo "Having all peers join the channel..."
joinChannel

## Set the anchor peers for each org in the channel
echo "Updating anchor peers for pocbank..."
updateAnchorPeers 0
sleep 2
echo "Updating anchor peers for pocinsurance..."
updateAnchorPeers 1
sleep 2
echo "Updating anchor peers for pocappraiser..."
updateAnchorPeers 2
sleep 2
echo "Updating anchor peers for pocgovernment..."
updateAnchorPeers 3
sleep 2
echo "Updating anchor peers for pocbuyer..."
updateAnchorPeers 4
sleep 2
echo "Updating anchor peers for pocseller..."
updateAnchorPeers 5
sleep 2
echo "Updating anchor peers for pochomelend..."
updateAnchorPeers 6
sleep 2
echo "Updating anchor peers for poccreditratingagency..."
updateAnchorPeers 7
sleep 2

echo "Installing chaincode on pocbank/peer0..."
installChaincode 0
echo "Installing chaincode on pocinsurance/peer0..."
installChaincode 1
echo "Installing chaincode on pocappraiser/peer0..."
installChaincode 2
echo "Installing chaincode on pocgovernment/peer0..."
installChaincode 3
echo "Installing chaincode on pocbuyer/peer0..."
installChaincode 4
echo "Installing chaincode on pocseller/peer0..."
installChaincode 5
echo "Installing chaincode on pochomelend/peer0..."
installChaincode 6
echo "Installing chaincode on poccreditratingagency/peer0..."
installChaincode 7

echo "Instantiating chaincode on pocbank/peer0..."
instantiateChaincode 0

echo
echo "========= All GOOD, Network execution completed =========== "
echo

exit 0
