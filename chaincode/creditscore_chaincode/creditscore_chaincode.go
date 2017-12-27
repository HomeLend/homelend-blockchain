package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

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

	if function == "query" {
		return t.query(stub, args[0], args[1])
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// todo: implement credit score system
func (t *HomelendChaincode) query(stub shim.ChaincodeStubInterface, salary string, loanAmount string) pb.Response {
	fmt.Println(fmt.Sprintf("creditscore started %s", salary, loanAmount))

	salaryInt, err := strconv.Atoi(salary)
	if err != nil {
		str := fmt.Sprintf("Salary is invalid %+v it must be an integer", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	switch {
	case salaryInt < 10000:
		return shim.Success([]byte("C"))
	case salaryInt > 10000 && salaryInt < 20000:
		return shim.Success([]byte("B"))
	case salaryInt > 20000:
		return shim.Success([]byte("A"))
	default:
		return shim.Error("Error while calculating credit score")
	}
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
