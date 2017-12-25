package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
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
	Amount      int    `json:"Amount"`
	Timestamp   int    `json:"Timestamp"`
}

type Person struct {
	Hash        string `json:"Hash"`
	Firstname   string `json:"BookNumber"`
	Lastname    string `json:"SerialNumber"`
	Citizenship string `json:"Street"`
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
	function, args := stub.GetFunctionAndParameters()
	//TODO: Code for checking the role of the user
	if function == "create" {
		return t.create(stub, args)
	} else if function == "update" {
		return t.update(stub, args)
	} else if function == "get" {
		return t.get(stub, args)
	} else if function == "query" {
		return t.query(stub, args[0])
	} else if function == "sell" {
		return t.sell(stub, args)
	} else if function == "buy" {
		return t.query(stub, args[0])
	} else if function == "registerBank" {
		return t.query(stub, args[0])
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleChaincode) registerHomelendUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil);
}

func (t *SimpleChaincode) registerBank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil);
}

func (t *SimpleChaincode) registerAppraiser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil);
}

func (t *SimpleChaincode) registerInsurance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil);
}

func (t *SimpleChaincode) registerSeller(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil);
}

func (t *SimpleChaincode) registerBuyer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil);
}

func (t *SimpleChaincode) sell(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("sell executed with args: %+v", args))

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

	mspid, err := cid.GetMSPID(stub)

	if err!= nil {
		str := fmt.Sprintf("MSPID error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := cid.GetID(stub)

	if err!= nil {
		str := fmt.Sprintf("MSPID error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	if mspid != "POCSellerMSP" {
		str := fmt.Sprintf("Only Seller Node can execute this method error %+v", mspid)
		fmt.Println(str)
		return shim.Error(str)
	}

	data := &House{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	dataAsBytes, err := stub.GetState(identity)
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", data.Hash)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		var arrayOfData []*House;
		arrayOfData = append(arrayOfData, data)

		dataJSONasBytes, err := json.Marshal(arrayOfData)
		if err != nil {
			str := fmt.Sprintf("Could not marshal %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}

		err = stub.PutState(identity, dataJSONasBytes)
		if err != nil {
			str := fmt.Sprintf("Could not put state %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}

		fmt.Println("Sucessfully executed");
	} else {
		var arrayOfData []*House;
		err = json.Unmarshal(dataAsBytes, arrayOfData)

		if err != nil {
			str := fmt.Sprintf("Failed to unmarshal: %s", err)
			fmt.Println(str)
			return shim.Error(str)
		}

		arrayOfData = append(arrayOfData, data)
		arrayOfDataAsBytes, err := json.Marshal(arrayOfData)

		err = stub.PutState(identity, arrayOfDataAsBytes)
		if err != nil {
			str := fmt.Sprintf("Could not put state %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) buy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil);
}

// Create method creates entity
func (t *SimpleChaincode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("create executed with args: %+v", args))

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

	data := &House{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	dataAsBytes, err := stub.GetState(data.Hash)
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", data.Hash)
		fmt.Println(str)
		return shim.Error(str)
	} else if dataAsBytes != nil {
		str := fmt.Sprintf("Record already exists: %s", data.Hash)
		fmt.Println(str)
		return shim.Error(str)
	}
	dataJSONasBytes, err := json.Marshal(data)
	if err != nil {
		str := fmt.Sprintf("Could not marshal %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	err = stub.PutState(data.Hash, dataJSONasBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("Sucessfully executed");
	return shim.Success(nil)
}

func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	data := &House{}
	err = json.Unmarshal([]byte(args[0]), data)

	dataAsBytes, err := stub.GetState(data.Hash)

	if err != nil {
		str := fmt.Sprintf("Failed to get: %+v", err.Error());
		fmt.Println(str)
		return shim.Error(str)
	} else if dataAsBytes == nil {
		str := fmt.Sprintf("Record does not exists: %s", data.Hash)
		fmt.Println(str)
		return shim.Error(str)
	}

	dataJSONasBytes, err := json.Marshal(data)

	if err != nil {
		str := fmt.Sprintf("Can not marshal %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	err = stub.PutState(data.Hash, dataJSONasBytes)

	if err != nil {
		str := fmt.Sprintf("Can not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("Successfully updated")
	return shim.Success(nil)
}

func (t *SimpleChaincode) get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		str := fmt.Sprintf("Incorrect number(%d) of arguments", len(args))
		fmt.Println(str)
		return shim.Error(str)
	}

	hash := args[0]
	valAsBytes, err := stub.GetState(hash)

	if err != nil {
		str := fmt.Sprintf("Failed to get state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	} else if valAsBytes == nil {
		str := fmt.Sprintf("Record does not exist %s", hash)
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("Successfully got")
	return shim.Success(valAsBytes)
}

func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, queryString string) pb.Response {
	fmt.Println(fmt.Sprintf("query started %s", queryString))
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		fmt.Println(fmt.Sprintf("incorrect query: %s", queryString))
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Println("Sucessfully queried")
	return shim.Success(buffer.Bytes())
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
