package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Request struct {
	Hash         string          `json:"Hash"`
	PropertyHash string          `json:"Name"`
	Buyer        string          `json:"Buyer"`
	Salary       int             `json:"TotalSupply"`
	LoanAmount   int             `json:"LoanAmount"`
	Status       int             `json:"Status"`
	Statuses     []RequestStatus `json:"RequestStatuses"`
	Timestamp    int             `json:"Timestamp"`
}

type RequestStatus struct {
	Status    int `json:"Status"`
	Timestamp int `json:"Timestamp"`
}

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
		numOfArgsResult := t.validateNumOfArgs(stub, args, 1)
		if len(numOfArgsResult) > 0 {
			return shim.Error(numOfArgsResult)
		}

		return t.query(stub, args[0], args[1])
	} else if function == "checkHouseOwner" {
		numOfArgsResult := t.validateNumOfArgs(stub, args, 1)
		if len(numOfArgsResult) > 0 {
			return shim.Error(numOfArgsResult)
		}

		request := getRequest(args[0])
		return t.checkHouseOwner(stub, request)
	} else if function == "checkLien" {
		numOfArgsResult := t.validateNumOfArgs(stub, args, 1)
		if len(numOfArgsResult) > 0 {
			return shim.Error(numOfArgsResult)
		}

		request := getRequest(args[0])
		return t.checkLien(stub, request)
	} else if function == "checkWarningShot" {
		numOfArgsResult := t.validateNumOfArgs(stub, args, 1)
		if len(numOfArgsResult) > 0 {
			return shim.Error(numOfArgsResult)
		}

		request := getRequest(args[0])
		return t.checkWarningShot(stub, request)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *HomelendChaincode) validateNumOfArgs(stub shim.ChaincodeStubInterface, args []string, count int) string {
	if len(args) != count {
		str := fmt.Sprintf("Incorrect number of arguments %d.", len(args))
		fmt.Println(str)
		return str
	}
	return ""
}

func (t *HomelendChaincode) query(stub shim.ChaincodeStubInterface, arg1 string, arg2 string) pb.Response {
	fmt.Println(fmt.Sprintf("government started"))

	return shim.Success([]byte("OK"))
}

func getRequest(requestStr string) *Request {
	data := &Request{}
	err := json.Unmarshal([]byte(requestStr), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return nil
	}
	return data
}

func (t *HomelendChaincode) checkLien(stub shim.ChaincodeStubInterface, request *Request) pb.Response {
	fmt.Println(fmt.Sprintf("government checkLien"))

	if request != nil {
		val := rand.Intn(100)
		result := val > 95
		if result {
			s := make([]byte, 1)
			return shim.Success(s)
		}
		s := make([]byte, 0)
		return shim.Success(s)
	}
	return shim.Error("Request is null")
}

func (t *HomelendChaincode) checkHouseOwner(stub shim.ChaincodeStubInterface, request *Request) pb.Response {
	fmt.Println(fmt.Sprintf("government checkHouseOwner"))

	if request != nil {
		val := rand.Intn(100)
		result := val > 95
		if result {
			s := make([]byte, 1)
			return shim.Success(s)
		}
		s := make([]byte, 0)
		return shim.Success(s)
	}
	return shim.Error("Request is null")
}

func (t *HomelendChaincode) checkWarningShot(stub shim.ChaincodeStubInterface, request *Request) pb.Response {
	fmt.Println(fmt.Sprintf("government checkWarningShot"))

	if request != nil {
		val := rand.Intn(100)
		result := val > 95
		if result {
			s := make([]byte, 1)
			return shim.Success(s)
		}
		s := make([]byte, 0)
		return shim.Success(s)
	}
	return shim.Error("Request is null")
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
