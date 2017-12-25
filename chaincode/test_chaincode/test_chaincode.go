package main

import (
// "bytes"
// "encoding/json"
"fmt"
"github.com/hyperledger/fabric/core/chaincode/shim"
// "github.com/hyperledger/fabric/common/crypto"
// "github.com/hyperledger/fabric/protos/common"
"github.com/hyperledger/fabric/core/chaincode/lib/cid"
pb "github.com/hyperledger/fabric/protos/peer"
// "encoding/binary"
// "crypto/x509"
// "encoding/pem"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// House chaincode
type House struct {
	Hash        string `json:"Hash"`
	FlatNumber  string `json:"BookNumber"`
	HouseNumber string `json:"SerialNumber"`
	Street      string `json:"Street"`
	Owner       string `json:"Owner"`
	Amount      string `json:"Amount"`
	Status      string `json:"Status"`
	Timestamp   int    `json:"Timestamp"`
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke chaincode methods
// ===========================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// function, args := stub.GetFunctionAndParameters()

	creator, _ := stub.GetCreator()
	fmt.Printf("Role is %v\n", string(creator))

	id, _ := cid.GetID(stub)
	mspid, _ := cid.GetMSPID(stub)
	cert, _ := cid.GetX509Certificate(stub)

	fmt.Printf("GETID1 %v\n", id)
	fmt.Printf("GETID2 %v\n", mspid)
	fmt.Printf("GETID3 %v\n", cert.PublicKey)
	/*if function == "create" {
		return t.create(stub, args)
	} else if function == "update" {
		return t.update(stub, args)
	} else if function == "get" {
		return t.get(stub, args)
	} else if function == "query" {
		return t.query(stub, args[0])
	}*/
	// fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
} // Create method creates entity

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
