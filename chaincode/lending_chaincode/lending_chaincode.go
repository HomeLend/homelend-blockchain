package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

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
* 4  - BANK_OFFER_INSTALLED
* 5  - BUYER_SELECTED_BANK_OFFER
* 6  - REQUEST_APPRAISER_CHOSEN
* 7  - APPRAISER_PROVIEDED_AMOUNT
* 8  - INSURANCE_OFFER_PROVIDED
* 9  - INSURANCE_OFFER_SELECTED
* 10  - REQUEST_GOVERNMENT_PROVIDED
* 11  - REQUEST_DECLINED_BY_BANK
* 12  - REQUEST_APPROVED_BY_BANK
* 13  - REQUEST_COMPLETED-ACTIVE-MORTGAGE

 */

// HomelendChaincode basic struct to provide an API
type HomelendChaincode struct {
}

const requests = "requests_"
const bank = "bank_"
const appraiser = "appraiser_"
const insuranceCompany = "insuranceCompany_"
const money = "money_"

const appraiserList = "appraiserList"
const properties4sale = "properties4sale"
const moneyTransferKey = "moneyTransfer"

//states
const creditRankOpenRequests = "creditRankOpenRequests"
const open4bankOffers = "open4bankoffers"
const selectAppraiser = "selectAppraiser"
const appraiserSelected = "appraiserSelected"
const open4InsuranceOffers = "open4InsuranceOffers"
const pending4Government = "pending4Government"
const pending4bankApproval = "pending4bankApproval"
const pending4ChaincodeExecute = "pending4ChaincodeExecute"

//include appraiser hash as suffix
const pendingForAppraiserEstimation = "pendingForAppraiserEstimation_"

// Property describes structure of real estate
type Property struct {
	Hash         string    `json:"Hash"`
	SellerHash   string    `json:"SellerHash"`
	Address      string    `json:"Address"`
	ImageBase64  string    `json:"ImageBase64"`
	SellingPrice float32   `json:"SellingPrice"`
	Timestamp    time.Time `json:"Timestamp"`
}

// InsuranceOffer describes fields of offer
type InsuranceOffer struct {
	Hash            string    `json:"Hash"`
	InsuranceHash   string    `json:"InsuranceHash"`
	InsuranceAmount float32   `json:"InsuranceAmount"`
	Timestamp       time.Time `json:"Timestamp"`
}

// Request defines buy processing and contains
type Request struct {
	Hash                       string             `json:"Hash"`
	PropertyHash               string             `json:"PropertyHash"`
	BuyerHash                  string             `json:"BuyerHash"`
	SellerHash                 string             `json:"SellerHash"`
	AppraiserHash              string             `json:"AppraiserHash"`
	AppraiserAmount            int                `json:"AppraiserAmount"`
	CreditScore                string             `json:"CreditScore"`
	CreditScoreIdentity        string             `json:"CreditScoreIdentity"`
	LoanAmountLeftToRefund     int                `json:"LoanAmountLeftToRefund"`
	GovernmentResultsData      *GovernmentResults `json:"GovernmentResultsData"`
	InsuranceOffers            []InsuranceOffer   `json:"InsuranceOffers"`
	BankOffers                 []BankOffer        `json:"BankOffers"`
	SelectedBankOfferHash      string             `json:"SelectedBankOfferHash"`
	SelectedInsuranceOfferHash string             `json:"SelectedInsuranceOfferHash"`
	Salary                     int                `json:"Salary"`
	SalaryBase64               string             `json:"SalaryBase64"`
	LoanAmount                 int                `json:"LoanAmount"`
	Duration                   int                `json:"Duration"`
	Status                     string             `json:"Status"`
	DeclineInfo                string             `json:"DeclineInfo"`
	Timestamp                  time.Time          `json:"Timestamp"`
}

//GovernmentResults - The results from the government
type GovernmentResults struct {
	CheckLien        bool      `json:"CheckLien"`
	CheckHouseOwner  bool      `json:"CheckHouseOwner"`
	CheckWarningShot bool      `json:"CheckWarningShot"`
	Timestamp        time.Time `json:"Timestamp"`
}

//RequestLink - pointer to request
type RequestLink struct {
	UserHash    string `json:"UserHash"`
	RequestHash string `json:"RequestHash"`
}

//BankOffer Bank offer
type BankOffer struct {
	Hash           string    `json:"Hash"`
	BankHash       string    `json:"BankHash"`
	Interest       float32   `json:"Interest"`
	MonthlyPayment float64   `json:"MonthlyPayment"`
	Timestamp      time.Time `json:"Timestamp"`
}

// Bank describes fields of Bank
type Bank struct {
	SwiftNumber string    `json:"SwiftNumber"`
	Name        string    `json:"Name"`
	TotalSupply int       `json:"TotalSupply"`
	Timestamp   time.Time `json:"Timestamp"`
}

// Seller structure describes the seller fields
type Seller struct {
	FullName  string    `json:"FullName"`
	Email     string    `json:"Email"`
	IDNumber  string    `json:"IDNumber"`
	Timestamp time.Time `json:"Timestamp"`
}

// Buyer describes fields necessary for buyer
type Buyer struct {
	FullName  string    `json:"FullName"`
	Email     string    `json:"Email"`
	IDNumber  string    `json:"IDNumber"`
	IDBase64  string    `json:"IDBase64"`
	Timestamp time.Time `json:"Timestamp"`
}

// Appraiser describes fields necessary for appraiser
type Appraiser struct {
	AppraiserHash string    `json:"AppraiserHash"`
	FirstName     string    `json:"FirstName"`
	LastName      string    `json:"LastName"`
	Email         string    `json:"Email"`
	IDNumber      string    `json:"IDNumber"`
	Timestamp     time.Time `json:"Timestamp"`
}

// InsuranceCompany describes fields necessary for insurance company
type InsuranceCompany struct {
	Hash          string    `json:"Hash"`
	LicenseNumber string    `json:"LicenseNumber"`
	Name          string    `json:"Name"`
	Address       string    `json:"Address"`
	Timestamp     time.Time `json:"Timestamp"`
}

// CreditRatingAgency describes fields necessary for credit rating agency/company
type CreditRatingAgency struct {
	LicenseNumber string    `json:"LicenseNumber"`
	Name          string    `json:"Name"`
	Address       string    `json:"Address"`
	Timestamp     time.Time `json:"Timestamp"`
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
	} else if function == "putSellerPersonalInfo" {
		return t.putSellerPersonalInfo(stub, args)
	} else if function == "appraiserputPersonalInfo" {
		return t.appraiserPutPersonalInfo(stub, args)
	} else if function == "appraiserProvideAmount" {
		return t.appraiserProvideAmount(stub, args)
	} else if function == "putBankInfo" {
		return t.putBankInfo(stub, args)
	} else if function == "creditRatingPull" {
		return t.getArray(stub, "POCCreditRatingAgencyMSP", creditRankOpenRequests, false)
	} else if function == "putCreditRatingAgencyInfo" {
		return t.putCreditRatingAgencyInfo(stub, args)
	} else if function == "putInsuranceCompanyInfo" {
		return t.putInsuranceCompanyInfo(stub, args)
	} else if function == "insurancePutOffer" {
		return t.insurancePutOffer(stub, args)
	} else if function == "buyerSelectAppraiser" {
		return t.buyerSelectAppraiser(stub, args)
	} else if function == "creditScore" {
		return t.calcCreditScore(stub, args)
	} else if function == "getProperties" {
		return t.getArray(stub, "", "", true)
	} else if function == "getRequestInfo" {
		return t.getRequestForSpecificPlayer(stub, args[0], args[1])
	} else if function == "appraiserPullPendingRequests" {
		return t.getArray(stub, "POCAppraiserMSP", pendingForAppraiserEstimation, true)
	} else if function == "buyerGetMyRequests" {
		return t.getArray(stub, "POCBuyerMSP", requests, true)
	} else if function == "governmentPullPending" {
		return t.getArray(stub, "POCGovernmentMSP", pending4Government, false)
	} else if function == "governmentPutData" {
		return t.governmentPutData(stub, args)
	} else if function == "bankApprove" {
		return t.bankApprove(stub, args)
	} else if function == "bankRunChaincode" {
		return t.bankRunChaincode(stub, args)
	} else if function == "bankPullOpen4bankOffers" {
		return t.getArray(stub, "POCBankMSP", open4bankOffers, false)
	} else if function == "bankPutOffer" {
		return t.bankPutOffer(stub, args)
	} else if function == "buyerGetAllAppraisers" {
		return t.getArray(stub, "POCBuyerMSP", appraiserList, false)
	} else if function == "buyerSelectBankOffer" {
		return t.buyerSelectBankOffer(stub, args)
	} else if function == "buyerSelectInsuranceOffer" {
		return t.buyerSelectInsuranceOffer(stub, args)
	} else if function == "getProperties4Sale" {
		return t.getArray(stub, "POCBuyerMSP", properties4sale, false)
	} else if function == "insuranceGetOpenRequests" {
		return t.getArray(stub, "POCInsuranceMSP", open4InsuranceOffers, false)
	} else if function == "query" {
		return t.query(stub, args[0])
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

//seller
func (t *HomelendChaincode) putSellerPersonalInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	identity, err := t.getIdentity(stub, "POCSellerMSP")
	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", args)
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

	data.Timestamp = time.Now()
	dataJSONasBytes, err := json.Marshal(data)
	if err != nil {
		str := fmt.Sprintf("Could not marshal %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	err = stub.PutState("seller-"+identity, dataJSONasBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("Sucessfully executed")

	return shim.Success(nil)
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

	identity, err := t.getIdentity(stub, "POCSellerMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	data := &Property{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}
	data.SellerHash = identity
	data.Timestamp = time.Now()
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

	dataAsBytes, err = stub.GetState(properties4sale)

	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", properties4sale)
		fmt.Println(str)
		return shim.Error(str)
	}

	if dataAsBytes == nil {
		fmt.Println("properties4sale Does not have houses. Creating first one")
		var arrayOfData []*Property
		arrayOfData = append(arrayOfData, data)

		dataJSONasBytes, err := json.Marshal(arrayOfData)
		if err != nil {
			str := fmt.Sprintf("properties4sale Could not marshal %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}

		err = stub.PutState(properties4sale, dataJSONasBytes)
		if err != nil {
			str := fmt.Sprintf("properties4sale Could not put state %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}

		fmt.Println("properties4sale Sucessfully executed")
	} else {
		fmt.Println("Already has houses. Appending one")
		var arrayOfData []*Property
		err = json.Unmarshal(dataAsBytes, &arrayOfData)

		if err != nil {
			str := fmt.Sprintf("properties4sale Failed to unmarshal: %s", err)
			fmt.Println(str)
			return shim.Error(str)
		}

		for i := 0; i < len(arrayOfData); i++ {
			if arrayOfData[i].Hash == data.Hash {
				str := fmt.Sprintf("proprty already exists in properties4sale hash: %s : %s", data.Hash, err)
				fmt.Println(str)
				return shim.Error(str)
			}
		}

		arrayOfData = append(arrayOfData, data)
		arrayOfDataAsBytes, err := json.Marshal(arrayOfData)

		err = stub.PutState(properties4sale, arrayOfDataAsBytes)
		if err != nil {
			str := fmt.Sprintf("Could not put state %+v", err.Error())
			fmt.Println(str)
			return shim.Error(str)
		}
	}

	return shim.Success(nil)
}

//government
func (t *HomelendChaincode) governmentPutData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("updateGovernmentResult executed with args: %+v", args))

	var err error
	if len(args) != 5 {
		str := fmt.Sprintf("Incorrect number of arguments %d.", len(args))
		fmt.Println(str)
		return shim.Error(str)
	}

	buyerHash := args[0]
	requestHash := args[1]

	checkHouseOwner := args[2] == "true"
	checkLien := args[3] == "true"
	checkWarningShot := args[4] == "true"

	_, err = t.getIdentity(stub, "POCGovernmentMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	request, err := t.getRequest(stub, buyerHash, requestHash)
	if err != nil {
		str := fmt.Sprintf("getRequest error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	request.GovernmentResultsData = t.govResultsGetter(checkHouseOwner, checkLien, checkWarningShot)
	request.Status = "REQUEST_GOVERNMENT_PROVIDED"

	err = t.addOrUpdateRequest(stub, request)
	if err != nil {
		str := fmt.Sprintf("Could not addOrUpdateRequest %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	rl := &RequestLink{UserHash: request.BuyerHash, RequestHash: request.Hash}
	t.removeFromRequestArray(stub, pending4Government, rl)
	if err != nil {
		str := fmt.Sprintf("Could not removeFromRequestArray %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}
	err = t.addRequestToArray(stub, pending4bankApproval, rl)
	if err != nil {
		str := fmt.Sprintf("Could not addRequestToArray %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("Successfully updated")
	return shim.Success(nil)
}

//appraiser
func (t *HomelendChaincode) appraiserPutPersonalInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("putAppraiserPersonalInfo executed with args: %+v", args))

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

	identity, err := t.getIdentity(stub, "POCAppraiserMSP")
	if err != nil {
		str := fmt.Sprintf("error getIdentity %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	appraiser := &Appraiser{}
	err = json.Unmarshal([]byte(args[0]), appraiser)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	appraiser.AppraiserHash = identity
	appraiser.Timestamp = time.Now()
	dataJSONasBytes, err := json.Marshal(appraiser)
	if err != nil {
		str := fmt.Sprintf("Could not marshal %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	//I think it should be only in array for now

	// err = stub.PutState(appraiser+identity, dataJSONasBytes)
	// if err != nil {
	// 	str := fmt.Sprintf("Could not put state %+v", err.Error())
	// 	fmt.Println(str)
	// 	return shim.Error(str)
	// }

	var aprList []*Appraiser
	valAsBytes, err := stub.GetState(appraiserList)
	if err != nil {
		str := fmt.Sprintf("Failed to get state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	if len(valAsBytes) > 0 {
		err = json.Unmarshal(valAsBytes, &aprList)
		if err != nil {
			str := fmt.Sprintf("Failed to unmarshal: %s", err)
			return shim.Error(str)
		}
	}

	found := false
	for i := 0; i < len(aprList); i++ {
		if aprList[i].AppraiserHash == appraiser.AppraiserHash {
			aprList[i] = appraiser
			found = true
		}
	}

	if !found {
		aprList = append(aprList, appraiser)
	}

	dataJSONasBytes, err = json.Marshal(aprList)
	if err != nil {
		str := fmt.Sprintf("Could not marshal %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}
	err = stub.PutState(appraiserList, dataJSONasBytes)
	if err != nil {
		str := fmt.Sprintf("Could PutState->appraiserList %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("putAppraiserPersonalInfo Sucessfully executed")
	return shim.Success(nil)
}

func (t *HomelendChaincode) appraiserProvideAmount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("appraiserProvideAmount executed with args"))

	if len(args) != 3 {
		str := fmt.Sprintf("Incorrect number of arguments %d.", len(args))
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := t.getIdentity(stub, "POCAppraiserMSP")
	if err != nil {
		str := fmt.Sprintf("error getIdentity %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	buyerHash := args[0]
	requestHash := args[1]
	amount := args[2]
	appraiserAmount, err := strconv.Atoi(amount)
	if err != nil {
		str := fmt.Sprintf("Could not parse amount %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	request, err := t.getRequest(stub, buyerHash, requestHash)
	if err != nil {
		str := fmt.Sprintf("Could not getRequest %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	if request.AppraiserHash != identity {
		str := "This appraiser has no access to this request"
		fmt.Println(str)
		return shim.Error(str)
	}

	request.Status = "APPRAISER_PROVIEDED_AMOUNT"
	request.AppraiserAmount = appraiserAmount
	t.addOrUpdateRequest(stub, request)
	if err != nil {
		str := fmt.Sprintf("Could not addOrUpdateRequest %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	rl := &RequestLink{UserHash: request.BuyerHash, RequestHash: request.Hash}
	err = t.removeFromRequestArray(stub, pendingForAppraiserEstimation+identity, rl)
	if err != nil {
		str := fmt.Sprintf("Could not removeFromRequestArray %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	err = t.addRequestToArray(stub, open4InsuranceOffers, rl)
	if err != nil {
		str := fmt.Sprintf("Could not addRequestToArray %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	return shim.Success(nil)
}

//insurance
func (t *HomelendChaincode) putInsuranceCompanyInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("putInsuranceCompanyInfo executed with args: %+v", args))

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

	identity, err := t.getIdentity(stub, "POCInsuranceMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	data := &InsuranceCompany{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	data.Hash = identity
	data.Timestamp = time.Now()
	dataJSONasBytes, err := json.Marshal(data)
	if err != nil {
		str := fmt.Sprintf("Could not marshal %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	err = stub.PutState(insuranceCompany+identity, dataJSONasBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("Sucessfully executed")

	return shim.Success(nil)
}
func (t *HomelendChaincode) insurancePutOffer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("updateInsuranceOffers executed with args: %+v", args))

	var err error
	if len(args) != 4 {
		str := fmt.Sprintf("Incorrect number of arguments %d.", len(args))
		fmt.Println(str)
		return shim.Error(str)
	}

	if len(args[0]) <= 0 {
		str := fmt.Sprintf("Provide hash for the request %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	if len(args[1]) <= 0 {
		str := fmt.Sprintf("Provide Insurance Amount for the request %+v", args[0])
		fmt.Println(str)
		return shim.Error(str)
	}

	if len(args[2]) <= 0 {
		str := fmt.Sprintf("Provide Insurance Offer Hash for the request %+v", args[0])
		fmt.Println(str)
		return shim.Error(str)
	}

	userHash := args[0]
	requestHash := args[1]
	amountStr := args[2]
	newHash := args[3]

	identity, err := t.getIdentity(stub, "POCInsuranceMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	request, err := t.getRequest(stub, userHash, requestHash)
	if err != nil {
		str := fmt.Sprintf("getRequest error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	amount, err := strconv.ParseFloat(amountStr, 32)
	if err != nil {
		str := fmt.Sprintf("Amount value is wrong %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	//TODO: check if insurance is registered
	request.Status = "INSURANCE_OFFER_PROVIDED"
	request.InsuranceOffers = append(request.InsuranceOffers, InsuranceOffer{Hash: newHash, InsuranceHash: identity, InsuranceAmount: float32(amount), Timestamp: time.Now()})

	err = t.addOrUpdateRequest(stub, request)
	if err != nil {
		str := fmt.Sprintf("addOrUpdateRequest %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	rl := &RequestLink{UserHash: request.BuyerHash, RequestHash: request.Hash}
	err = t.removeFromRequestArray(stub, open4InsuranceOffers, rl)
	fmt.Println("Successfully updated")
	return shim.Success(nil)
}

//Bank
func (t *HomelendChaincode) putBankInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("putBankInfo executed with args: %+v", args))

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

	identity, err := t.getIdentity(stub, "POCBankMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	data := &Bank{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	data.Timestamp = time.Now()
	dataJSONasBytes, err := json.Marshal(data)
	if err != nil {
		str := fmt.Sprintf("Could not marshal %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	err = stub.PutState(bank+identity, dataJSONasBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("Sucessfully executed")

	return shim.Success(nil)
}

func (t *HomelendChaincode) bankPutOffer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("bankPutOffer executed with args: %+v", args))

	var err error
	if len(args) != 3 {
		str := fmt.Sprintf("Incorrect number of arguments %d.", len(args))
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := t.getIdentity(stub, "POCBankMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	//input validations
	if len(args[0]) <= 0 {
		str := fmt.Sprintf("Provide RequestLink for the request")
		fmt.Println(str)
		return shim.Error(str)
	}
	if len(args[1]) <= 0 {
		str := fmt.Sprintf("Provide Bank Offer Hash for the request %+v", args[0])
		fmt.Println(str)
		return shim.Error(str)
	}
	if len(args[2]) <= 0 {
		str := fmt.Sprintf("Provide Bank Interest for the request  %+v", args[0])
		fmt.Println(str)
		return shim.Error(str)
	}

	requestLinkStr := args[0]
	requestLink := &RequestLink{}
	err = json.Unmarshal([]byte(requestLinkStr), requestLink)
	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal requestLinkStr: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	hash := args[1]
	interest, err := strconv.ParseFloat(args[2], 32)
	if err != nil {
		str := fmt.Sprintf("interest args[2] - invalid input %+v", args[0])
		fmt.Println(str)
		return shim.Error(str)
	}

	offer := &BankOffer{BankHash: identity, Hash: hash, Interest: float32(interest), Timestamp: time.Now()}

	request, err := t.getRequest(stub, requestLink.UserHash, requestLink.RequestHash)
	if err != nil {
		str := fmt.Sprintf("Failed to getRequest: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	offer.MonthlyPayment, err = t.calcPmt(float64(interest), request.Duration, float64(request.LoanAmount))
	if err != nil {
		str := fmt.Sprintf("Failed calcPmt: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	request.Status = "BANK_OFFER_INSTALLED"
	request.BankOffers = append(request.BankOffers, *offer)
	err = t.addOrUpdateRequest(stub, request)
	if err != nil {
		str := fmt.Sprintf("Failed to addOrUpdateRequest: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	err = t.removeFromRequestArray(stub, open4bankOffers, requestLink)
	if err != nil {
		str := fmt.Sprintf("Failed to removeFromRequestArray: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("BankOffer -> Successfully updated")
	return shim.Success(nil)
}

func (t *HomelendChaincode) bankApprove(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("bankApprove executed with args: %+v", args))

	var err error
	if len(args) != 1 {
		str := fmt.Sprintf("Incorrect number of arguments %d.", len(args))
		fmt.Println(str)
		return shim.Error(str)
	}

	bankIdentity, err := t.getIdentity(stub, "POCBankMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	requestLinkStr := args[0]
	requestLink := &RequestLink{}
	err = json.Unmarshal([]byte(requestLinkStr), requestLink)
	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal requestLinkStr: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	request, err := t.getRequest(stub, requestLink.UserHash, requestLink.RequestHash)
	if err != nil {
		str := fmt.Sprintf("Failed to getRequest: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	validations, err := t.bankValidateBeforeApprove(request, bankIdentity)
	if err != nil {
		str := fmt.Sprintf("error in bankValidateBeforeApprove: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	if len(validations) != 0 {
		request.Status = "REQUEST_DECLINED_BY_BANK"
		request.DeclineInfo = validations
		err = t.addOrUpdateRequest(stub, request)
		if err != nil {
			if err != nil {
				str := fmt.Sprintf("saddOrUpdateRequest - unapproved %s", err)
				fmt.Println(str)
				return shim.Error(str)
			}
		}

		return shim.Success(nil)
	}

	escrowAccountKey := money + "escrow_" + request.Hash
	err = stub.PutState(escrowAccountKey, []byte(strconv.Itoa(request.LoanAmount)))
	if err != nil {
		str := fmt.Sprintf("PutState %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	request.LoanAmountLeftToRefund = request.LoanAmount
	request.Status = "REQUEST_APPROVED_BY_BANK"

	err = t.addOrUpdateRequest(stub, request)
	if err != nil {
		str := fmt.Sprintf("saddOrUpdateRequest %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	err = t.removeFromRequestArray(stub, pending4bankApproval, requestLink)
	if err != nil {
		str := fmt.Sprintf("removeFromRequestArray %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	err = t.addRequestToArray(stub, pending4ChaincodeExecute, requestLink)
	if err != nil {
		str := fmt.Sprintf("addRequestToArray %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("bankApprove -> Successfully updated")
	return shim.Success(nil)
}

func (t *HomelendChaincode) bankRunChaincode(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("bankRunChaincode executed with args: %+v", args))

	var err error
	if len(args) != 1 {
		str := fmt.Sprintf("Incorrect number of arguments %d.", len(args))
		fmt.Println(str)
		return shim.Error(str)
	}

	bankIdentity, err := t.getIdentity(stub, "POCBankMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	requestLinkStr := args[0]
	requestLink := &RequestLink{}
	err = json.Unmarshal([]byte(requestLinkStr), &requestLink)
	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal requestLinkStr: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	request, err := t.getRequest(stub, requestLink.UserHash, requestLink.RequestHash)
	if err != nil {
		str := fmt.Sprintf("Failed to getRequest: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	err = t.validateBankOwner(request, bankIdentity)
	if err != nil {
		str := fmt.Sprintf("validateBankOwner Failed: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	property, err := t.getPropertyAndRemove(stub, request.SellerHash, request.PropertyHash)
	if err != nil {
		str := fmt.Sprintf("Failed to getPropertyAndRemove: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	dataAsBytes, err := stub.GetState(request.BuyerHash)
	if err != nil {
		str := fmt.Sprintf("Failed to get Property from user: %s", request.BuyerHash)
		fmt.Println(str)
		return shim.Error(str)
	}

	var buyerPropertylist []*Property
	if len(dataAsBytes) > 0 {
		err = json.Unmarshal(dataAsBytes, &buyerPropertylist)
		if err != nil {
			str := fmt.Sprintf("Could not Unmarshal %+v", err.Error())
			return shim.Error(str)
		}
	}

	buyerPropertylist = append(buyerPropertylist, property)
	buyerPropertylistBytes, err := json.Marshal(buyerPropertylist)
	if err != nil {
		str := fmt.Sprintf("Could not Marshal %+v", err.Error())
		return shim.Error(str)
	}

	err = stub.PutState(request.BuyerHash, buyerPropertylistBytes)
	if err != nil {
		str := fmt.Sprintf("Could not PutState %+v", err.Error())
		return shim.Error(str)
	}

	err = t.moveMoney(stub, "escrow_"+request.Hash, request.SellerHash, request.LoanAmount)
	if err != nil {
		str := fmt.Sprintf("Could not moveMoney %+v", err.Error())
		return shim.Error(str)
	}

	request.Status = "REQUEST_COMPLETED-ACTIVE-MORTGAGE"
	err = t.addOrUpdateRequest(stub, request)
	if err != nil {
		str := fmt.Sprintf("Could not addOrUpdateRequest %+v", err.Error())
		return shim.Error(str)
	}

	err = t.removeFromRequestArray(stub, pending4ChaincodeExecute, requestLink)
	if err != nil {
		str := fmt.Sprintf("removeFromRequestArray %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("bankRunChaincode -> Successfully updated")
	return shim.Success(nil)
}

//creditScore
func (t *HomelendChaincode) putCreditRatingAgencyInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("putCreditRatingAgencyInfo executed with args: %+v", args))

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

	identity, err := t.getIdentity(stub, "POCCreditRatingAgencyMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	data := &CreditRatingAgency{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	data.Timestamp = time.Now()
	dataJSONasBytes, err := json.Marshal(data)
	if err != nil {
		str := fmt.Sprintf("Could not marshal %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	err = stub.PutState("credit-rating-agency-"+identity, dataJSONasBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	fmt.Println("Sucessfully executed")

	return shim.Success(nil)
}

func (t *HomelendChaincode) calcCreditScore(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("calcCreditScore executed with args"))

	var err error
	if len(args) != 1 {
		str := fmt.Sprintf("Incorrect number of arguments %d.", len(args))
		fmt.Println(str)
		return shim.Error(str)
	}

	if len(args[0]) <= 0 {
		str := fmt.Sprintf("Argument must be non-empty string %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	requestLinkStr := args[0]

	identity, err := t.getIdentity(stub, "POCCreditRatingAgencyMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	requestLink := &RequestLink{}
	err = json.Unmarshal([]byte(requestLinkStr), requestLink)
	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	if err != nil {
		str := fmt.Sprintf("Could not marshal array %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	request, err := t.getRequest(stub, requestLink.UserHash, requestLink.RequestHash)
	if err != nil {
		str := fmt.Sprintf("Could not getRequest %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	// bargs := make([][]byte, 2)
	// bargs[0] = []byte(strconv.Itoa(request.Salary))
	// bargs[1] = []byte(strconv.Itoa(request.LoanAmount))

	// if err != nil {
	// 	str := fmt.Sprintf("Could not marshal array %+v", err.Error())
	// 	fmt.Println(str)
	// 	return shim.Error(str)
	// }
	// resp := stub.InvokeChaincode("creditscore_chaincode", bargs, "mainchannel")
	// if resp.Status != 200 {
	// 	fmt.Println("Could InvokeChaincode creditscore_chaincode")
	// 	return resp
	// }

	// strResult := string(resp.Payload)
	strResult, _ := t.getCreditRankScore(stub, request.Salary, request.LoanAmount)

	request.CreditScore = strResult
	request.CreditScoreIdentity = identity
	request.Status = "REQUEST_CREDIT_SCORE_INSTALLED"
	err = t.addOrUpdateRequest(stub, request)
	if err != nil {
		str := fmt.Sprintf("Could not updateRequest %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	rl := &RequestLink{UserHash: request.BuyerHash, RequestHash: request.Hash}
	err = t.removeFromRequestArray(stub, creditRankOpenRequests, rl)
	if err != nil {
		str := fmt.Sprintf("Could not removeFromRequestArray %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	err = t.addRequestToArray(stub, open4bankOffers, rl)
	if err != nil {
		str := fmt.Sprintf("Could not addRequestToArray open4bankOffers %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	return shim.Success(nil)
}

//buyer
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

	identity, err := t.getIdentity(stub, "POCBuyerMSP")
	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", args)
		fmt.Println(str)
		return shim.Error(str)
	}

	data := &Buyer{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	data.Timestamp = time.Now()
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

	identity, err := t.getIdentity(stub, "POCBuyerMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity error %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	data := &Request{}
	err = json.Unmarshal([]byte(args[0]), data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	data.Status = "REQUEST_INITIALIZED"
	data.BuyerHash = identity
	data.Timestamp = time.Now()

	fmt.Println(fmt.Printf("Getting state for %+s", identity))
	err = t.addOrUpdateRequest(stub, data)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}
	rl := &RequestLink{UserHash: identity, RequestHash: data.Hash}
	t.addRequestToArray(stub, creditRankOpenRequests, rl)

	//check if the property exists and remove from array
	dataAsBytes, err := stub.GetState(properties4sale)
	var properties4saleArray []*Property

	if err != nil {
		str := fmt.Sprintf("Failed to get: %s", properties4sale)
		fmt.Println(str)
		return shim.Error(str)
	}

	err = json.Unmarshal(dataAsBytes, &properties4saleArray)
	if err != nil {
		str := fmt.Sprintf("properties4sale Failed to unmarshal: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	proptyIndex := -1
	for i := 0; i < len(properties4saleArray); i++ {
		if properties4saleArray[i].Hash == data.PropertyHash && properties4saleArray[i].SellerHash == data.SellerHash {
			proptyIndex = i
		}
	}

	if proptyIndex < 0 {
		str := fmt.Sprintf("propery does not exists in properties4sale array hash: %s sellerHash: %s", data.PropertyHash, data.SellerHash)
		fmt.Println(str)
		return shim.Error(str)
	}

	properties4saleArray = append(properties4saleArray[:proptyIndex], properties4saleArray[proptyIndex+1:]...)

	dataAsBytes, err = json.Marshal(properties4saleArray)
	if err != nil {
		str := fmt.Sprintf("properties4sale Failed to Marshal again: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}
	stub.PutState(properties4sale, dataAsBytes)

	return shim.Success([]byte(identity))
}

func (t *HomelendChaincode) getProperties(stub shim.ChaincodeStubInterface) pb.Response {
	identity, err := t.getIdentity(stub, "POCBuyerMSP")

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

func (t *HomelendChaincode) buyerSelectBankOffer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	identity, err := t.getIdentity(stub, "POCBuyerMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	if len(args) != 2 {
		str := fmt.Sprintf("Incorrect number of arguments %d.", len(args))
		fmt.Println(str)
		return shim.Error(str)
	}

	requestHash := args[0]
	selectedBankOfferHash := args[1]

	request, err := t.getRequest(stub, identity, requestHash)
	if err != nil {
		str := fmt.Sprintf("getRequest error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	if request == nil || len(request.BankOffers) < 1 {
		str := fmt.Sprintf("No BankOffers %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	found := false

	for i := 0; i < len(request.BankOffers); i++ {
		if request.BankOffers[i].Hash == selectedBankOfferHash {
			found = true
		}
	}

	if !found {
		str := fmt.Sprintf("Bank Offer was not found %+v", selectedBankOfferHash)
		fmt.Println(str)
		return shim.Error(str)
	}

	request.Status = "BUYER_SELECTED_BANK_OFFER"
	request.SelectedBankOfferHash = selectedBankOfferHash
	t.addOrUpdateRequest(stub, request)

	rl := &RequestLink{UserHash: request.BuyerHash, RequestHash: request.Hash}
	err = t.removeFromRequestArray(stub, open4bankOffers, rl)
	if err != nil {
		str := fmt.Sprintf("Could not removeFromRequestArray %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	err = t.addRequestToArray(stub, selectAppraiser, rl)
	if err != nil {
		str := fmt.Sprintf("Could not addRequestToArray open4bankOffers %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}
	return shim.Success(nil)
}

func (t *HomelendChaincode) buyerSelectAppraiser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("buyerSelectAppraiser executed with args %+v", args))

	if len(args) != 2 {
		str := fmt.Sprintf("Incorrect number of arguments %d.", len(args))
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := t.getIdentity(stub, "POCBuyerMSP")
	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	requestHash := args[0]
	appraiserHash := args[1]

	request, err := t.getRequest(stub, identity, requestHash)
	if err != nil {
		str := fmt.Sprintf("Failed getRequest: %s", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	//TODO: check if exists in appraiser list
	request.AppraiserHash = appraiserHash
	request.Status = "REQUEST_APPRAISER_CHOSEN"
	t.addOrUpdateRequest(stub, request)

	rl := &RequestLink{UserHash: request.BuyerHash, RequestHash: request.Hash}
	err = t.removeFromRequestArray(stub, selectAppraiser, rl)
	if err != nil {
		str := fmt.Sprintf("Could not removeFromRequestArray %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	err = t.addRequestToArray(stub, pendingForAppraiserEstimation+appraiserHash, rl)
	if err != nil {
		str := fmt.Sprintf("Could not addRequestToArray open4bankOffers %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	return shim.Success(nil)
}

func (t *HomelendChaincode) buyerSelectInsuranceOffer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println(fmt.Sprintf("buyerSelectInsuranceOffer executed with args"))

	var err error
	if len(args) != 2 {
		str := fmt.Sprintf("Incorrect number of arguments %d.", len(args))
		fmt.Println(str)
		return shim.Error(str)
	}

	identity, err := t.getIdentity(stub, "POCBuyerMSP")
	if err != nil {
		str := fmt.Sprintf("getIdentity error %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	requestHash := args[0]
	offerHash := args[1]

	request, err := t.getRequest(stub, identity, requestHash)
	if err != nil {
		str := fmt.Sprintf("getRequest:  %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	found := false
	for i := 0; i < len(request.InsuranceOffers); i++ {
		if request.InsuranceOffers[i].Hash == offerHash {
			found = true
		}
	}

	if !found {
		return shim.Error("offer was not found")
	}

	request.Status = "INSURANCE_OFFER_SELECTED"
	request.SelectedInsuranceOfferHash = offerHash
	err = t.addOrUpdateRequest(stub, request)
	if err != nil {
		str := fmt.Sprintf("addOrUpdateRequest:  %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	rl := &RequestLink{UserHash: request.BuyerHash, RequestHash: request.Hash}
	t.addRequestToArray(stub, pending4Government, rl)
	return shim.Success(nil)
}

//fix
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

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(HomelendChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

//helper
func (t *HomelendChaincode) addRequestToArray(stub shim.ChaincodeStubInterface, key string, data *RequestLink) error {
	dataAsBytes, err := stub.GetState(key)
	if err != nil {
		return err
	}

	if dataAsBytes == nil {
		var arrayOfData []*RequestLink
		arrayOfData = append(arrayOfData, data)

		dataJSONasBytes, err := json.Marshal(arrayOfData)
		if err != nil {
			str := fmt.Sprintf("Could not marshal %+v", err.Error())
			return errors.New(str)
		}

		err = stub.PutState(key, dataJSONasBytes)
		if err != nil {
			str := fmt.Sprintf("Could not put state %+v", err.Error())
			return errors.New(str)
		}

	} else {
		fmt.Println("Already has requests. Appending one")
		var arrayOfData []*RequestLink
		err = json.Unmarshal(dataAsBytes, &arrayOfData)

		if err != nil {
			str := fmt.Sprintf("Failed to unmarshal: %s", err)
			return errors.New(str)
		}

		for i := 0; i < len(arrayOfData); i++ {
			if arrayOfData[i].UserHash == data.UserHash && arrayOfData[i].RequestHash == data.RequestHash {
				str := fmt.Sprintf("item already exists in array: %s : %s", key, err)
				fmt.Println(str)
				return errors.New(str)
			}
		}

		arrayOfData = append(arrayOfData, data)
		arrayOfDataAsBytes, err := json.Marshal(arrayOfData)
		if err != nil {
			str := fmt.Sprintf("Could not Marshal %+v", err.Error())
			return errors.New(str)
		}

		err = stub.PutState(key, arrayOfDataAsBytes)
		if err != nil {
			str := fmt.Sprintf("Could not put state %+v", err.Error())
			return errors.New(str)
		}
	}

	return nil
}

func (t *HomelendChaincode) removeFromRequestArray(stub shim.ChaincodeStubInterface, arrayName string, data *RequestLink) error {

	dataAsBytes, err := stub.GetState(arrayName)
	if err != nil {
		return err
	}
	if dataAsBytes == nil {
		return errors.New("removeFromRequestArray -> src is empty")
	}

	var arrayOfData []*RequestLink
	err = json.Unmarshal(dataAsBytes, &arrayOfData)

	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal: %s", err)
		return errors.New(str)
	}

	found := false
	for i := 0; len(arrayOfData) > i; i++ {
		if data.RequestHash == arrayOfData[i].RequestHash {
			arrayOfData = append(arrayOfData[:i], arrayOfData[i+1:]...)
			found = true
		}
	}

	if !found {
		return errors.New("RequestLink was not found in array " + arrayName)
	}

	arrayOfDataAsBytes, err := json.Marshal(arrayOfData)
	if err != nil {
		str := fmt.Sprintf("Failed to marshal: %s", err)
		return errors.New(str)
	}

	err = stub.PutState(arrayName, arrayOfDataAsBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		return errors.New(str)
	}
	return nil
}

func (t *HomelendChaincode) getIdentity(stub shim.ChaincodeStubInterface, mspidValue string) (string, error) {

	mspid, err := cid.GetMSPID(stub)

	if err != nil {
		str := fmt.Sprintf("MSPID error %+v", err)
		fmt.Println(str)
		return "", errors.New(str)
	}

	identity, err := cid.GetID(stub)

	if err != nil {
		str := fmt.Sprintf("GetID error %+v", err)
		fmt.Println(str)
		return "", errors.New(str)
	}

	if mspidValue != "" && mspid != mspidValue {
		str := fmt.Sprintf("Only "+mspidValue+" Node can execute this method not "+mspid+" error %+v", mspid)
		fmt.Println(str)
		return "", errors.New(str)
	}

	return identity, nil
}

func (t *HomelendChaincode) getRequest(stub shim.ChaincodeStubInterface, userHash string, requestHash string) (*Request, error) {

	key := requests + userHash
	dataAsBytes, err := stub.GetState(key)

	if err != nil {
		str := fmt.Sprintf("Failed to get state %+v", err.Error())
		fmt.Println(str)
		return nil, errors.New(str)
	} else if dataAsBytes == nil {
		str := fmt.Sprintf("Record does not exist - empty array %s", key)
		fmt.Println(str)
		return nil, errors.New(str)
	}

	var arrayOfData []*Request
	err = json.Unmarshal(dataAsBytes, &arrayOfData)

	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal: %s", err)
		return nil, errors.New(str)
	}

	for i := 0; len(arrayOfData) > i; i++ {
		if arrayOfData[i].Hash == requestHash {
			return arrayOfData[i], nil
		}
	}

	return nil, errors.New("Not found")
}

func (t *HomelendChaincode) addOrUpdateRequest(stub shim.ChaincodeStubInterface, request *Request) error {

	key := requests + request.BuyerHash
	dataAsBytes, err := stub.GetState(key)

	if err != nil {
		str := fmt.Sprintf("Failed to get state %+v", err.Error())
		fmt.Println(str)
		return errors.New(str)
	}

	var arrayOfData []*Request
	if dataAsBytes == nil {
		arrayOfData = make([]*Request, 0)
	} else {
		err = json.Unmarshal(dataAsBytes, &arrayOfData)
		if err != nil {
			str := fmt.Sprintf("Failed to unmarshal: %s", err)
			return errors.New(str)
		}
	}

	found := false
	for i := 0; len(arrayOfData) > i; i++ {
		if arrayOfData[i].Hash == request.Hash {
			arrayOfData[i] = request
			found = true
		}
	}

	if !found {
		arrayOfData = append(arrayOfData, request)
	}

	arrayOfDataAsBytes, err := json.Marshal(arrayOfData)
	if err != nil {
		str := fmt.Sprintf("Could not Marshal %+v", err.Error())
		return errors.New(str)
	}

	err = stub.PutState(key, arrayOfDataAsBytes)
	if err != nil {
		str := fmt.Sprintf("Could not put state %+v", err.Error())
		return errors.New(str)
	}

	return nil
}

func (t *HomelendChaincode) getArray(stub shim.ChaincodeStubInterface, msp string, arrayName string, addIdentityasSuffix bool) pb.Response {
	str := fmt.Sprintf("getArray= %s MSP= %s", arrayName, msp)
	fmt.Println(str)

	identity, err := t.getIdentity(stub, msp)
	if err != nil {
		return shim.Error(err.Error())
	}

	if addIdentityasSuffix {
		arrayName += identity
	}

	valAsBytes, err := stub.GetState(arrayName)
	if err != nil {
		str := fmt.Sprintf("Failed to get state %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	} else if valAsBytes == nil {
		str := fmt.Sprintf("Record does not exist %s", arrayName)
		fmt.Println(str)
		return shim.Success(nil)
	}

	fmt.Println("Successfully got " + creditRankOpenRequests)
	return shim.Success(valAsBytes)
}

func (t *HomelendChaincode) getRequestForSpecificPlayer(stub shim.ChaincodeStubInterface, userHash string, requestHash string) pb.Response {
	str := fmt.Sprintf("getRequestForSpecificPlayer= userHash %s requestHash= %s", userHash, requestHash)
	fmt.Println(str)

	mspid, err := cid.GetMSPID(stub)

	request, err := t.getRequest(stub, userHash, requestHash)
	if err != nil {
		str := fmt.Sprintf("Failed to getRequest %+v", err.Error())
		fmt.Println(str)
		return shim.Error(str)
	}

	//TODO: limit the rquestInfo for each player
	if mspid == "POCBankMSP1" {
		return shim.Success(nil)
	}

	byteArr, err := json.Marshal(request)
	if err != nil {
		str := fmt.Sprintf("Could not Marshal request %+v", err.Error())
		return shim.Error(str)
	}

	return shim.Success(byteArr)
}

func (t *HomelendChaincode) getCreditRankScore(stub shim.ChaincodeStubInterface, salary int, loanAmount int) (string, error) {
	switch {
	case salary < 10000:
		return "C", nil
	case salary > 10000 && salary < 20000:
		return "B", nil
	case salary > 20000:
		return "A", nil
	default:
		return "", errors.New("invalid salary?")
	}
}

func (t *HomelendChaincode) govResultsGetter(checkHouseOwner bool, checkLien bool, checkWarningShot bool) *GovernmentResults {

	result := &GovernmentResults{}

	result.CheckHouseOwner = checkHouseOwner
	result.CheckLien = checkLien
	result.CheckWarningShot = checkWarningShot
	result.Timestamp = time.Now()

	return result
}

func (t *HomelendChaincode) bankValidateBeforeApprove(request *Request, bankIdentity string) (string, error) {

	if request.GovernmentResultsData.CheckHouseOwner {
		return "", nil
	}
	if request.GovernmentResultsData.CheckLien {
		return "", nil
	}
	if request.GovernmentResultsData.CheckWarningShot {
		return "", nil
	}

	delta := float32(request.AppraiserAmount) * 0.1
	if float32(request.AppraiserAmount) < (float32(request.LoanAmount) + delta) {
		return "appraiser amount is too low for this loan", nil
	}

	if len(request.SelectedInsuranceOfferHash) == 0 {
		return "No insurance offer was selected", nil
	}

	err := t.validateBankOwner(request, bankIdentity)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (t *HomelendChaincode) validateBankOwner(request *Request, bankIdentity string) error {

	for i := 0; i < len(request.BankOffers); i++ {
		if request.BankOffers[i].Hash == request.SelectedBankOfferHash && request.BankOffers[i].BankHash != bankIdentity {
			return errors.New("this bank cannot approve offer that other bank created")
		}
	}

	return nil
}

func (t *HomelendChaincode) moveMoney(stub shim.ChaincodeStubInterface, srcUserID string, destUserID string, sum int) error {

	srcMoney := t.getMoney(stub, srcUserID)
	if srcMoney < 0 {
		str := fmt.Sprintf("Could not get money from srcUserID" + srcUserID)
		fmt.Println(str)
		return errors.New(str)
	}

	destMoney := t.getMoney(stub, destUserID)
	if srcMoney < 0 {
		str := fmt.Sprintf("Could not get money from srcUserID" + destUserID)
		fmt.Println(str)
		return errors.New(str)
	}

	if srcMoney < sum {
		str := fmt.Sprintf("not enough money in srcMoney %d and price is %d", srcMoney, sum)
		fmt.Println(str)
		return errors.New(str)
	}

	srcMoney -= sum
	destMoney += sum

	srcMoneyStr := strconv.Itoa(srcMoney)
	destMoneyStr := strconv.Itoa(destMoney)

	stub.PutState(money+srcUserID, []byte(srcMoneyStr))
	stub.PutState(money+destUserID, []byte(destMoneyStr))

	return nil

	moneyTransferBytes, err := stub.GetState(moneyTransferKey)

	if err != nil {
		return err
	}

	var moneyTransferArr []string
	err = json.Unmarshal(moneyTransferBytes, &moneyTransferArr)
	if err != nil {
		return err
	}

	tx2Append := srcUserID + "_" + destUserID + "_" + strconv.Itoa(sum)
	moneyTransferArr = append(moneyTransferArr, tx2Append)

	moneyTransferBytes, err = json.Marshal(moneyTransferArr)
	if err != nil {
		return err
	}
	stub.PutState(moneyTransferKey, moneyTransferBytes)
	return nil
}

func (t *HomelendChaincode) getMoney(stub shim.ChaincodeStubInterface, userID string) int {

	dataAsBytes, err := stub.GetState(money + userID)
	if err != nil {
		str := fmt.Sprintf("Could not getMoney userID=" + userID)
		fmt.Println(str)
		return -1
	}

	if len(dataAsBytes) == 0 {
		return 0
	}

	escrowMoney, err := strconv.Atoi(string(dataAsBytes))
	if err != nil {
		str := fmt.Sprintf("Could not getMoney - Atoi userID=" + userID)
		fmt.Println(str)
		return -1
	}

	return escrowMoney
}

func (t *HomelendChaincode) getProperty(stub shim.ChaincodeStubInterface, userHash string, propertyHash string) (*Property, int, []*Property, error) {

	dataAsBytes, err := stub.GetState(userHash)
	if err != nil {
		str := fmt.Sprintf("Failed to get Property from user: %s", userHash)
		fmt.Println(str)
		return nil, 0, nil, err
	}

	if len(dataAsBytes) <= 0 {
		str := fmt.Sprintf("Empty properties for user: %s", userHash)
		fmt.Println(str)
		return nil, 0, nil, err
	}

	var list []*Property
	err = json.Unmarshal(dataAsBytes, &list)
	if err != nil {
		return nil, 0, nil, err
	}

	for i := 0; i < len(list); i++ {
		if list[i].Hash == propertyHash {
			return list[i], i, list, nil
		}
	}

	return nil, 0, nil, errors.New("Could not found property: " + propertyHash + "in property array of user" + userHash)
}

func (t *HomelendChaincode) getPropertyAndRemove(stub shim.ChaincodeStubInterface, userHash string, propertyHash string) (*Property, error) {

	property, index, properties, err := t.getProperty(stub, userHash, propertyHash)
	if err != nil {
		str := fmt.Sprintf("Failed to getProperty %+v", err.Error())
		fmt.Println(str)
		return nil, errors.New(str)
	}

	properties = append(properties[:index], properties[index+1:]...)
	propertiesArrayBytes, err := json.Marshal(properties)
	if err != nil {
		return nil, errors.New("Failed Marshal: properties")
	}

	err = stub.PutState(userHash, propertiesArrayBytes)
	if err != nil {
		return nil, errors.New("Failed PutState: propertiesArrayBytes")
	}
	return property, nil
}

func (t *HomelendChaincode) calcPmt(yearlyInterestRate float64, totalNumberOfMonths int, loanAmount float64) (float64, error) {
	fmt.Println("calcPmt", yearlyInterestRate, totalNumberOfMonths, loanAmount)
	if yearlyInterestRate > 100 || yearlyInterestRate < 0 {
		return 0, errors.New("invalid: interest")
	}

	if totalNumberOfMonths < 1 || totalNumberOfMonths > 500 {
		return 0, errors.New("invalid: duration")
	}

	if loanAmount < 1 || loanAmount > 100000000 {
		return 0, errors.New("invalid: loanAmount")
	}

	rate := yearlyInterestRate / 100 / 12
	denominator := math.Pow((1+rate), float64(totalNumberOfMonths)) - 1
	result := (rate + (rate / denominator)) * loanAmount
	return result, nil
}
