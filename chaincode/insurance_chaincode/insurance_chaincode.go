package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/homelend-blockchain/chaincode/homelendlib"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//HomelendChaincode is...
type HomelendChaincode struct {
}

// Init initializes chaincode
// ===========================
func (t *HomelendChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke chaincode methods
// ===========================
func (t *HomelendChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	identity, err := cid.GetID(stub)

	if err != nil {
		str := fmt.Sprintf("Identity error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	mspid, err := cid.GetMSPID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println(fmt.Printf("Access log %s %s", identity, mspid))

	if function == "putOffers" {
		return t.putOffer(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *HomelendChaincode) putOffer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("buy executed with args: %+v", args))

	helpers := lib.Helpers{}
	var err error
	if len(args) != 1 {
		str := fmt.Sprintf("Incorrect number of arguments %d.", len(args))
		fmt.Println(str)
		return shim.Error(str)
	}

	if len(args[0]) <= 0 {
		str := fmt.Sprintf("JSON must be non-empty string %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	userID := args[0]

	mspid, err := cid.GetMSPID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := cid.GetID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	if mspid != "POCHomelendMSP" {
		str := fmt.Sprintf("Only POCHomelend Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	offer := lib.InsuranceOffer{Amount: 100 + rand.Intn(200), InsuranceHash: identity, Timestamp: int(time.Now().Unix()), Hash: "someHash"}

	offer.Amount = 12

	currentRequest, requestArr, err := helpers.GetLastRequest(stub, userID)

	currentRequest.InsuranceOffers = append(currentRequest.InsuranceOffers, offer)
	if err != nil {
		return helpers.PrintAndReturnError(stub, "could not get last user request for user "+userID)
	}

	err = helpers.UpdateLastRequest(stub, userID, requestArr)
	if err != nil {
		return helpers.PrintAndReturnError(stub, "could update user request "+userID)
	}
	return shim.Success(nil)
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(HomelendChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
