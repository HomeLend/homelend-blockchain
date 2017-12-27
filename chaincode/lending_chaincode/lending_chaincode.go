package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
)

// REQUEST STATUSES
const REQUEST_INITIALIZED = 1;
const REQUEST_CREDIT_SCORE = 2;
const REQUEST_CREDIT_SCORE_CONFIRMED = 3;
const REQUEST_CREDIT_SCORE_DECLINED = 4;


type HomelendChaincode struct {
}

type Property struct {
	Hash         string `json:"Hash"`
	Address      string `json:"Address"`
	SellingPrice int    `json:"SellingPrice"`
	Timestamp    int    `json:"Timestamp"`
}

type RequestStatus struct {
	Status    int `json:"Status"`
	Timestamp int `json:"Timestamp"`
}

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

type Bank struct {
	Hash        string `json:"Hash"`
	Name        string `json:"Name"`
	TotalSupply int    `json:"TotalSupply"`
	Timestamp   int    `json:"Timestamp"`
}

type Seller struct {
	ID        string `json:"ID"`
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
	Timestamp int    `json:"Timestamp"`
}

type Buyer struct {
	ID           string `json:"ID"`
	Firstname    string `json:"Firstname"`
	Lastname     string `json:"Lastname"`
	IDNumber     string `json:"IDNumber"`
	IDBase64     string `json:"IDBase64"`
	SalaryBase64 string `json:"SalaryBase64"`
	Timestamp    int    `json:"Timestamp"`
}

type Appraiser struct {
	ID        string `json:"ID"`
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
	Timestamp int    `json:"Timestamp"`
}

type InsuranceCompany struct {
	ID        string `json:"ID"`
	Name      string `json:"Firstname"`
	Address   string `json:"Lastname"`
	Timestamp int    `json:"Timestamp"`
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
		return t.query(stub, args[0])
	} else if function == "advertise" {
		return t.advertise(stub, args)
	} else if function == "getProperties" {
		return t.getProperties(stub)
	} else if function == "registerAsBank" {
		return t.registerAsBank(stub, args)
	} else if function == "registerAsSeller" {
		return t.registerAsSeller(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *HomelendChaincode) advertise(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	// todo: disable
	if mspid != "POCSellerMSP" {
		str := fmt.Sprintf("Only Seller Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	data := &Property{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println(fmt.Printf("Getting state for %+s", identity))
	dataAsBytes, err := stub.GetState(identity)
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		fmt.Println("Does not have houses. Creating first one")
		var arrayOfData []*Property;
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
		fmt.Println("Already has houses. Appending one")
		var arrayOfData []*Property;
		err = json.Unmarshal(dataAsBytes, &arrayOfData)

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

func (t *HomelendChaincode) buy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil);
}

func (t *HomelendChaincode) getProperties(stub shim.ChaincodeStubInterface) pb.Response {
	identity, err := cid.GetID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	valAsBytes, err := stub.GetState(identity)

	if err != nil {
		str := fmt.Sprintf("Failed to get state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	} else if valAsBytes == nil {
		str := fmt.Sprintf("Record does not exist %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("Successfully got")
	return shim.Success(valAsBytes)
}

func (t *HomelendChaincode) query(stub shim.ChaincodeStubInterface, queryString string) pb.Response {
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

func (t *HomelendChaincode) registerAsBank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("registerAsBank executed with args: %+v", args))

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

	if mspid != "POCBankMSP" {
		str := fmt.Sprintf("Only Bank Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	data := &Bank{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println(fmt.Printf("Getting state for %+s", "bank_list"))
	dataAsBytes, err := stub.GetState("bank_list")
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		var bankList []*Bank;
		bankList = append(bankList, data)

		bankListasBytes, err := json.Marshal(bankList)
		if err != nil {
			str := fmt.Sprintf("Could not marshal %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}

		err = stub.PutState("bank_list", bankListasBytes)
	} else {
		var bankList []*Bank;

		// todo: check for existing bank with hash

		err = json.Unmarshal(dataAsBytes, &bankList)
		if err != nil {
			str := fmt.Sprintf("Could not marshal %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}

		bankList = append(bankList, data)

		bankListasBytes, err := json.Marshal(bankList)
		if err != nil {
			str := fmt.Sprintf("Could not marshal %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}

		err = stub.PutState("bank_list", bankListasBytes)
	}

	return shim.Success(nil);
}

func (t *HomelendChaincode) registerAsSeller(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("registerAsSeller executed with args: %+v", args))

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

	if mspid != "POCSellerMSP" {
		str := fmt.Sprintf("Only Bank Node can execute this method error %+v", mspid)
		fmt.Println(str)
		return shim.Error(str)
	}

	data := &Seller{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println(fmt.Printf("Getting state for %+s", "seller_list"))
	dataAsBytes, err := stub.GetState("seller_list")
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		var sellerList []*Seller;
		sellerList = append(sellerList, data)

		listasBytes, err := json.Marshal(sellerList)
		if err != nil {
			str := fmt.Sprintf("Could not marshal %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}

		err = stub.PutState("seller_list", listasBytes)
	} else {
		var sellerList []*Seller;
		err = json.Unmarshal(dataAsBytes, &sellerList)
		if err != nil {
			str := fmt.Sprintf("Could not marshal %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}

		sellerList = append(sellerList, data)
		listAsBytes, err := json.Marshal(sellerList)
		if err != nil {
			str := fmt.Sprintf("Could not marshal %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}

		err = stub.PutState("seller_list", listAsBytes)
	}

	return shim.Success(nil);
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
