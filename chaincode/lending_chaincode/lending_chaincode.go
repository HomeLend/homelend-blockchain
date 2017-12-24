package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"bytes"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// House chaincode
type House struct {
	Hash            string `json:"Hash"`
	FlatNumber      string `json:"BookNumber"`
	HouseNumber     string `json:"SerialNumber"`
	Street          string `json:"Street"`
	Owner     		string `json:"Owner"`
	Amount          string `json:"Amount"`
	Status          string `json:"Status"`
	Timestamp      	int    `json:"Timestamp"`
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

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func Validate(chaincodeID string, hash string, stub shim.ChaincodeStubInterface) bool {
	//return true
	args := [2]string {"get", hash}
	for i := 0; i < 3; i++ {
		invokeArgs := make([][]byte, len(args))
		for i, arg := range args {
			invokeArgs[i] = []byte(arg)
		}

		response := stub.InvokeChaincode(chaincodeID, invokeArgs, "mainchannel")

		if response.Status == shim.OK {
			fmt.Println(fmt.Sprintf("Sucessfully validated %s on %s", chaincodeID, hash))
			return true;
		}

		fmt.Println(fmt.Sprintf("Failed to validate %s on %s", chaincodeID, hash))
	}
	return false;
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
	}
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
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


func (t *SimpleChaincode) OperateHouse(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "put":
		if len(args) < 2 {
			return shim.Error("put operation must include two arguments, a key and value")
		}
		key := args[0]
		value := args[1]

		if err := stub.PutState(key, []byte(value)); err != nil {
			fmt.Printf("Error putting state %s", err)
			return shim.Error(fmt.Sprintf("put operation failed. Error updating state: %s", err))
		}

		indexName := "compositeKeyTest"
		compositeKeyTestIndex, err := stub.CreateCompositeKey(indexName, []string{key})
		if err != nil {
			return shim.Error(err.Error())
		}

		valueByte := []byte{0x00}
		if err := stub.PutState(compositeKeyTestIndex, valueByte); err != nil {
			fmt.Printf("Error putting state with compositeKey %s", err)
			return shim.Error(fmt.Sprintf("put operation failed. Error updating state with compositeKey: %s", err))
		}

		return shim.Success(nil)

	case "remove":
		if len(args) < 1 {
			return shim.Error("remove operation must include one argument, a key")
		}
		key := args[0]

		err := stub.DelState(key)
		if err != nil {
			return shim.Error(fmt.Sprintf("remove operation failed. Error updating state: %s", err))
		}
		return shim.Success(nil)

	case "get":
		if len(args) < 1 {
			return shim.Error("get operation must include one argument, a key")
		}
		key := args[0]
		value, err := stub.GetState(key)
		if err != nil {
			return shim.Error(fmt.Sprintf("get operation failed. Error accessing state: %s", err))
		}
		return shim.Success(value)

	case "keys":
		if len(args) < 2 {
			return shim.Error("put operation must include two arguments, a key and value")
		}
		startKey := args[0]
		endKey := args[1]

		//sleep needed to test peer's timeout behavior when using iterators
		stime := 0
		if len(args) > 2 {
			stime, _ = strconv.Atoi(args[2])
		}

		keysIter, err := stub.GetStateByRange(startKey, endKey)
		if err != nil {
			return shim.Error(fmt.Sprintf("keys operation failed. Error accessing state: %s", err))
		}
		defer keysIter.Close()

		var keys []string
		for keysIter.HasNext() {
			//if sleeptime is specied, take a nap
			if stime > 0 {
				time.Sleep(time.Duration(stime) * time.Millisecond)
			}

			response, iterErr := keysIter.Next()
			if iterErr != nil {
				return shim.Error(fmt.Sprintf("keys operation failed. Error accessing state: %s", err))
			}
			keys = append(keys, response.Key)
		}

		for key, value := range keys {
			fmt.Printf("key %d contains %s\n", key, value)
		}

		jsonKeys, err := json.Marshal(keys)
		if err != nil {
			return shim.Error(fmt.Sprintf("keys operation failed. Error marshaling JSON: %s", err))
		}

		return shim.Success(jsonKeys)
	case "query":
		query := args[0]
		keysIter, err := stub.GetQueryResult(query)
		if err != nil {
			return shim.Error(fmt.Sprintf("query operation failed. Error accessing state: %s", err))
		}
		defer keysIter.Close()

		var keys []string
		for keysIter.HasNext() {
			response, iterErr := keysIter.Next()
			if iterErr != nil {
				return shim.Error(fmt.Sprintf("query operation failed. Error accessing state: %s", err))
			}
			keys = append(keys, response.Key)
		}

		jsonKeys, err := json.Marshal(keys)
		if err != nil {
			return shim.Error(fmt.Sprintf("query operation failed. Error marshaling JSON: %s", err))
		}

		return shim.Success(jsonKeys)
	case "history":
		key := args[0]
		keysIter, err := stub.GetHistoryForKey(key)
		if err != nil {
			return shim.Error(fmt.Sprintf("query operation failed. Error accessing state: %s", err))
		}
		defer keysIter.Close()

		var keys []string
		for keysIter.HasNext() {
			response, iterErr := keysIter.Next()
			if iterErr != nil {
				return shim.Error(fmt.Sprintf("query operation failed. Error accessing state: %s", err))
			}
			keys = append(keys, response.TxId)
		}

		for key, txID := range keys {
			fmt.Printf("key %d contains %s\n", key, txID)
		}

		jsonKeys, err := json.Marshal(keys)
		if err != nil {
			return shim.Error(fmt.Sprintf("query operation failed. Error marshaling JSON: %s", err))
		}

		return shim.Success(jsonKeys)

	default:
		return shim.Success([]byte("Unsupported operation"))
	}
}