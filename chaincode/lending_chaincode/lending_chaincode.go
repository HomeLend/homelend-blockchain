package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// todo: add explanation proper comments
// todo: license & terms of use of code

/* *
* STATUSES FOR THE BUYING PROCESS
*
* 1  - REQUEST_INITIALIZED
* 2  - REQUEST_DATA_PROVIDED
* 3  - REQUEST_CREDIT_SCORE_INSTALLED
* 4  - REQUEST_CREDIT_SCORE_CONFIRMED / REQUEST_CREDIT_SCORE_DECLINED
* 5  - REQUEST_DOCUMENT_UPLOADED
* 6  - REQUEST_DOCUMENTS_VERIFIED
* 7  - REQUEST_APPRAISER_CHOSEN
* 8  - REQUEST_APPRAISER_PROVIDED
* 9  - REQUEST_INSURANCE_DATA_PROVIDED
* 10 - REQUEST_INSURANCE_OFFER_CHOSEN
* 11 - REQUEST_GOVERNMENT_CONFIRMED/ REQUEST_GOVERNMENT_DECLINED
* 12 - REQUEST_SMART_CONTRACT_GENERATED
* 13 - REQUEST_TRANSACTION_HAPPENED
* 14 - COMPLETED / REQUEST_CREDIT_SCORE_DECLINED / REQUEST_GOVERNMENT_DECLINED
 */

// HomelendChaincode basic struct to provide an API
type HomelendChaincode struct {
}

// Property describes structure of real estate
type Property struct {
	Hash         string `json:"Hash"`
	SellerHash   string `json:"SellerHash"`
	Address      string `json:"Address"`
	SellingPrice int    `json:"SellingPrice"`
	Timestamp    int    `json:"Timestamp"`
}

// InsuranceOffer describes fields of offer
type InsuranceOffer struct {
	Hash          string `json:"Hash"`
	InsuranceHash string `json:"InsuranceHash"`
	Amount        string `json:"Amount"`
	Timestamp     int    `json:"Timestamp"`
}

// Request defines buy processing and contains
type Request struct {
	Hash              string           `json:"Hash"`
	PropertyHash      string           `json:"PropertyHash"`
	BuyerHash         string           `json:"BuyerHash"`
	SellerHash        string           `json:"SellerHash"`
	AppraiserHash     string           `json:"AppraiserHash"`
	CreditScore       string           `json:"CreditScore"`
	AppraiserPrice    string           `json:"AppraiserPrice"`
	AppraiserAmount   int              `json:"AppraiserAmount"`
	InsuranceHash     string           `json:"InsuranceHash"`
	InsuranceAmount   string           `json:"InsuranceAmount"`
	GovernmentResult1 string           `json:"GovernmentResult1"`
	GovernmentResult2 string           `json:"GovernmentResult2"`
	GovernmentResult3 string           `json:"GovernmentResult3"`
	InsuranceOffers   []InsuranceOffer `json:"InsuranceOffers"`
	Salary            int              `json:"Salary"`
	SalaryBase64      string           `json:"SalaryBase64"`
	LoanAmount        int              `json:"LoanAmount"`
	Duration          int              `json:"LoanAmount"`
	Status            string           `json:"Status"`
	Timestamp         int              `json:"Timestamp"`
}

// Bank describes fields of Bank
type Bank struct {
	Hash        string `json:"Hash"`
	Name        string `json:"Name"`
	TotalSupply int    `json:"TotalSupply"`
	Timestamp   int    `json:"Timestamp"`
}

// Seller structure describes the seller fields
type Seller struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email"`
	IDNumber  string `json:"IDNumber"`
	Timestamp int    `json:"Timestamp"`
}

// Buyer describes fields necessary for buyer
type Buyer struct {
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	Email        string `json:"Email"`
	IDNumber     string `json:"IDNumber"`
	IDBase64     string `json:"IDBase64"`
	Timestamp    int    `json:"Timestamp"`
}

// Appraiser describes fields necessary for appraiser
type Appraiser struct {
	ID        string `json:"ID"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Timestamp int    `json:"Timestamp"`
}

// InsuranceCompany describes fields necessary for insurance company
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

	// steps
	if function == "advertise" {
		return t.advertise(stub, args)
	} else if function == "buy" {
		return t.buy(stub, args)
	} else if function == "putBuyerPersonalInfo" {
		return t.putBuyerPersonalInfo(stub, args)
	} else if function == "creditScore" {
		return t.getCreditScore(stub)
	}

	// additional getters
	if function == "getProperties" {
		return t.getProperties(stub)
	} else if function == "pullBankOffers" {
     		return t.pullBankOffers(stub)
    } else if function == "query" {
		return t.query(stub, args[0])
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *HomelendChaincode) advertise(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("advertise executed with args: %+v", args))

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
		var arrayOfData []*Property
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

		fmt.Println("Sucessfully executed")
	} else {
		fmt.Println("Already has houses. Appending one")
		var arrayOfData []*Property
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

func (t *HomelendChaincode) putBuyerPersonalInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("putBuyerPersonalInfo executed with args: %+v", args))

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
	if mspid != "POCBuyerMSP" {
		str := fmt.Sprintf("Only Seller Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	data := &Buyer{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	dataJSONasBytes, err := json.Marshal(data)
	if err != nil {
		str := fmt.Sprintf("Could not marshal %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	err = stub.PutState("buyer-"+identity, dataJSONasBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	dataAsBytes, err := stub.GetState("requests_" + identity)
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	// -----
	var arrayOfData []*Request
	err = json.Unmarshal(dataAsBytes, &arrayOfData)

	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	currentRequest := arrayOfData[len(arrayOfData)-1]
	currentRequest.Status = "REQUEST_DATA_PROVIDED"

	arrayOfData[len(arrayOfData)-1] = currentRequest

	arrayOfDataAsBytes, err := json.Marshal(arrayOfData)

	err = stub.PutState("requests_"+identity, arrayOfDataAsBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("Sucessfully executed")

	return shim.Success(nil)
}

func (t *HomelendChaincode) buy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("buy executed with args: %+v", args))

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
	if mspid != "POCBuyerMSP" {
		str := fmt.Sprintf("Only Buyer Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	data := &Request{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	data.Status = "REQUEST_INITIALIZED"

	fmt.Println(fmt.Printf("Getting state for %+s", identity))
	dataAsBytes, err := stub.GetState("requests_" + identity)
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		fmt.Println("Does not have houses. Creating first one")
		var arrayOfData []*Request
		arrayOfData = append(arrayOfData, data)

		dataJSONasBytes, err := json.Marshal(arrayOfData)
		if err != nil {
			str := fmt.Sprintf("Could not marshal %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}

		err = stub.PutState("requests_"+identity, dataJSONasBytes)
		if err != nil {
			str := fmt.Sprintf("Could not put state %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}

		fmt.Println("Sucessfully executed")
	} else {
		fmt.Println("Already has requests. Appending one")
		var arrayOfData []*Request
		err = json.Unmarshal(dataAsBytes, &arrayOfData)

		if err != nil {
			str := fmt.Sprintf("Failed to unmarshal: %s", err)
			fmt.Println(str)
			return shim.Error(str)
		}

		arrayOfData = append(arrayOfData, data)
		arrayOfDataAsBytes, err := json.Marshal(arrayOfData)

		err = stub.PutState("requests_"+identity, arrayOfDataAsBytes)
		if err != nil {
			str := fmt.Sprintf("Could not put state %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}
	}

	return shim.Success(nil)
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

func (t *HomelendChaincode) pullBankOffers(stub shim.ChaincodeStubInterface) pb.Response {
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

func (t *HomelendChaincode) getCreditScore(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println(fmt.Sprintf("getCreditScore executed with args"))

	mspid, err := cid.GetMSPID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := cid.GetID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	// todo: disable
	if mspid != "POCBuyerMSP" {
		str := fmt.Sprintf("Only Buyer Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	dataAsBytes, err := stub.GetState("requests_" + identity)
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		str := "Does not have requests. Can't execute credit score"
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("Has requests. Executing credit score")
	var arrayOfData []*Request
	err = json.Unmarshal(dataAsBytes, &arrayOfData)

	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	latestRequest := arrayOfData[len(arrayOfData)-1]

	bargs := make([][]byte, 1)
	bargs[0], err = json.Marshal(latestRequest)

	if err != nil {
		str := fmt.Sprintf("Could not marshal array %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	resp := stub.InvokeChaincode("creditscore_chaincode", bargs, "mainchannel")

	strResult := string(resp.Payload)
	latestRequest.CreditScore = strResult
	latestRequest.Status = "REQUEST_CREDIT_SCORE_INSTALLED"
	arrayOfData[len(arrayOfData)-1] = latestRequest

	arrayOfDataAsBytes, err := json.Marshal(arrayOfData)

	err = stub.PutState("requests_"+identity, arrayOfDataAsBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	return shim.Success(nil)
}

func (t *HomelendChaincode) getGovernmentResult(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println(fmt.Sprintf("getGovernmentResult executed with args"))

	mspid, err := cid.GetMSPID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := cid.GetID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	// todo: disable
	if mspid != "POCBuyerMSP" {
		str := fmt.Sprintf("Only Buyer Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	dataAsBytes, err := stub.GetState("requests_" + identity)
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		str := "Does not have requests. Can't execute credit score"
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("Has requests. Executing credit score")
	var arrayOfData []*Request
	err = json.Unmarshal(dataAsBytes, &arrayOfData)

	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	latestRequest := arrayOfData[len(arrayOfData)-1]

	bargs := make([][]byte, 1)
	bargs[0], err = json.Marshal(latestRequest)

	if err != nil {
		str := fmt.Sprintf("Could not marshal array %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	resp := stub.InvokeChaincode("government_chaincode", bargs, "mainchannel")

	strResult := string(resp.Payload)
	latestRequest.CreditScore = strResult
	latestRequest.Status = "REQUEST_GOVERNMENT_RESP_INSTALLED"
	arrayOfData[len(arrayOfData)-1] = latestRequest

	arrayOfDataAsBytes, err := json.Marshal(arrayOfData)

	err = stub.PutState("requests_"+identity, arrayOfDataAsBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	return shim.Success(nil)
}

func (t *HomelendChaincode) getRequestsForBuyer(stub shim.ChaincodeStubInterface) pb.Response {
	identity, err := cid.GetID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	mspid, err := cid.GetMSPID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	// todo: disable§
	if mspid != "POCBuyerMSP" {
		str := fmt.Sprintf("Only Buyer Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	valAsBytes, err := stub.GetState("requests_" + identity)

	if err != nil {
		str := fmt.Sprintf("Failed to get state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	} else if valAsBytes == nil {
		str := fmt.Sprintf("Record does not exist %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	return shim.Success(valAsBytes)
}

func (t *HomelendChaincode) confirmCreditScore(stub shim.ChaincodeStubInterface, response string) pb.Response {
	fmt.Println(fmt.Sprintf("confirmCreditScore executed with args"))

	mspid, err := cid.GetMSPID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := cid.GetID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	// todo: disable
	if mspid != "POCBuyerMSP" {
		str := fmt.Sprintf("Only Buyer Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	dataAsBytes, err := stub.GetState("requests_" + identity)
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		str := "Does not have requests. Can't execute credit score"
		fmt.Println(str)
		return shim.Error(str)
	}

	var arrayOfData []*Request
	err = json.Unmarshal(dataAsBytes, &arrayOfData)

	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	latestRequest := arrayOfData[len(arrayOfData)-1]

	bargs := make([][]byte, 1)
	bargs[0], err = json.Marshal(latestRequest)

	if err != nil {
		str := fmt.Sprintf("Could not marshal array %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	if response == "confirm" {
		latestRequest.Status = "REQUEST_CREDIT_SCORE_CONFIRMED"
	} else {
		latestRequest.Status = "REQUEST_CREDIT_SCORE_DECLINED"
	}
	arrayOfData[len(arrayOfData)-1] = latestRequest

	arrayOfDataAsBytes, err := json.Marshal(arrayOfData)

	err = stub.PutState("requests_"+identity, arrayOfDataAsBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	return shim.Success(nil)
}

func (t *HomelendChaincode) buyerUploadDocuments(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("confirmCreditScore executed with args"))

	mspid, err := cid.GetMSPID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := cid.GetID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	// todo: disable
	if mspid != "POCBuyerMSP" {
		str := fmt.Sprintf("Only Buyer Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	data := &Buyer{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	dataAsBytes, err := stub.GetState("buyers_" + identity)
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		str := "Does not have requests. Can't execute credit score"
		fmt.Println(str)
		return shim.Error(str)
	}

	var savedData Buyer
	err = json.Unmarshal(dataAsBytes, &savedData)

	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	savedData.IDBase64 = data.IDBase64

	savedDataAsBytes, err := json.Marshal(savedData)

	err = stub.PutState("buyers_"+identity, savedDataAsBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	return shim.Success(nil)
}

func (t *HomelendChaincode) chooseAppraiser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("chooseAppraiser executed with args"))

	mspid, err := cid.GetMSPID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := cid.GetID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	// todo: disable
	if mspid != "POCBuyerMSP" {
		str := fmt.Sprintf("Only Buyer Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	appraiserHash := args[0]

	dataAsBytes, err := stub.GetState("requests_" + identity)
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		str := "Does not have any data. Can't execute"
		fmt.Println(str)
		return shim.Error(str)
	}

	var requests []*Request
	err = json.Unmarshal(dataAsBytes, &requests)

	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	latestRequest := requests[len(requests)-1]
	latestRequest.AppraiserHash = appraiserHash

	savedDataAsBytes, err := json.Marshal(dataAsBytes)

	err = stub.PutState("requests_"+identity, savedDataAsBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	return shim.Success(nil)
}

func (t *HomelendChaincode) chooseInsuranceOffer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("chooseInsuranceCompany executed with args"))

	mspid, err := cid.GetMSPID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := cid.GetID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	// todo: disable
	if mspid != "POCBuyerMSP" {
		str := fmt.Sprintf("Only Buyer Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	insuranceHash := args[0]
	insuranceAmount := args[1]

	dataAsBytes, err := stub.GetState("requests_" + identity)
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", identity)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		str := "Does not have any data. Can't execute"
		fmt.Println(str)
		return shim.Error(str)
	}

	var requests []*Request
	err = json.Unmarshal(dataAsBytes, &requests)

	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	latestRequest := requests[len(requests)-1]
	latestRequest.InsuranceHash = insuranceHash
	latestRequest.InsuranceAmount = insuranceAmount

	savedDataAsBytes, err := json.Marshal(dataAsBytes)

	err = stub.PutState("requests_"+identity, savedDataAsBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	return shim.Success(nil)
}

func (t *HomelendChaincode) appraiserProvideAmount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("appraiserProvideAmount executed with args"))

	mspid, err := cid.GetMSPID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := cid.GetID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	// todo: disable
	if mspid != "POCAppraiserMSP" {
		str := fmt.Sprintf("Only Buyer Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	buyerHash := args[0]
	appraiserAmount, _ := strconv.Atoi(args[1])

	dataAsBytes, err := stub.GetState("requests_" + buyerHash)
	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", buyerHash)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		str := "Does not have any data. Can't execute"
		fmt.Println(str)
		return shim.Error(str)
	}

	var requests []*Request
	err = json.Unmarshal(dataAsBytes, &requests)

	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	latestRequest := requests[len(requests)-1]
	latestRequest.AppraiserAmount = appraiserAmount

	savedDataAsBytes, err := json.Marshal(dataAsBytes)

	err = stub.PutState("requests_"+identity, savedDataAsBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	return shim.Success(nil)
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

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
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
