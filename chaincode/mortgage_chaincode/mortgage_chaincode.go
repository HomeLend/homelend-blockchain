package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//Bank is...
type Bank struct {
	Hash        string `json:"Hash"`
	Name        string `json:"Name"`
	TotalSupply int    `json:"TotalSupply"`
	Timestamp   int    `json:"Timestamp"`
}

//Seller is...
type Seller struct {
	ID        string `json:"ID"`
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
	Timestamp int    `json:"Timestamp"`
}

//Buyer is...
type Buyer struct {
	ID           string `json:"ID"`
	Firstname    string `json:"Firstname"`
	Lastname     string `json:"Lastname"`
	IDNumber     string `json:"IDNumber"`
	IDBase64     string `json:"IDBase64"`
	SalaryBase64 string `json:"SalaryBase64"`
	Timestamp    int    `json:"Timestamp"`
}

//Property is...
type Property struct {
	Hash         string `json:"Hash"`
	Address      string `json:"Address"`
	SellingPrice int    `json:"SellingPrice"`
	Timestamp    int    `json:"Timestamp"`
}

//Request is...
type Request struct {
	Hash         string `json:"Hash"`
	PropertyHash string `json:"Name"`
	BuyerHash    string `json:"BuyerHash"`
	SellerHash   string `json:"SellerHash"`
	CreditScore  string `json:"CreditScore"`
	Salary       int    `json:"TotalSupply"`
	LoanAmount   int    `json:"LoanAmount"`
	Status       string `json:"Status,omitempty"`
	Timestamp    int    `json:"Timestamp"`
}

//HomelendChaincode is...
type HomelendChaincode struct {
}

//MoneyTransferKey ...
const MoneyTransferKey = "moneyTransfer"

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

	if function == "buy" {
		return t.buy(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
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

	// identity, err := cid.GetID(stub)

	// if err != nil {
	// 	str := fmt.Sprintf("MSPID error %+v", args)
	// 	fmt.Println(str)
	// 	return shim.Error(str)
	// }

	if mspid != "POCHomelendMSP" {
		str := fmt.Sprintf("Only POCHomelend Node can execute this method error %+v", mspid)
		fmt.Println(str)
		// return shim.Error(str)
	}

	requestData := &Request{}
	err = json.Unmarshal([]byte(args[0]), requestData)
	if err != nil {
		str := fmt.Sprintf("Failed to parse JSON: %+v", err)
		fmt.Println(str)
		return shim.Error(str)
	}

	sellerProperty, sellerPropertyIndex, sellerPropertyArray, err := t.findHouse(stub, requestData.SellerHash, requestData.PropertyHash)
	if err != nil {
		str := fmt.Sprintf("Failed getting seller property")
		fmt.Println(str)
		return shim.Error(str)
	}

	buyerPropertyArray, err := t.getHouseList(stub, requestData.BuyerHash)
	if err != nil {
		str := fmt.Sprintf("Failed getting buyer property")
		fmt.Println(str)
		return shim.Error(str)
	}

	//add house to buyer
	if buyerPropertyArray == nil {
		buyerPropertyArray = append(buyerPropertyArray, sellerProperty)
	}

	//remove house from seller
	sellerPropertyArray = append(sellerPropertyArray[:sellerPropertyIndex], sellerPropertyArray[sellerPropertyIndex+1:]...)

	buyerPropertyArrayBytes, err := json.Marshal(buyerPropertyArray)
	if err != nil {
		str := fmt.Sprintf("Failed Marshal: buyerPropertyArray")
		fmt.Println(str)
		return shim.Error(str)
	}
	sellerPropertyArrayBytes, err := json.Marshal(buyerPropertyArray)
	if err != nil {
		str := fmt.Sprintf("Failed Marshal: buyerPropertyArray")
		fmt.Println(str)
		return shim.Error(str)
	}

	err = stub.PutState(requestData.BuyerHash, buyerPropertyArrayBytes)
	if err != nil {
		str := fmt.Sprintf("Failed PutState: BuyerHash->buyerPropertyArrayBytes")
		fmt.Println(str)
		return shim.Error(str)
	}
	err = stub.PutState(requestData.SellerHash, sellerPropertyArrayBytes)
	if err != nil {
		str := fmt.Sprintf("Failed PutState: SellerHash->sellerPropertyArrayBytes")
		fmt.Println(str)
		return shim.Error(str)
	}

	err = t.moveMoney(stub, "escrow_"+requestData.BuyerHash, requestData.SellerHash, sellerProperty.SellingPrice)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *HomelendChaincode) findHouse(stub shim.ChaincodeStubInterface, sellerID string, propertyID string) (*Property, int, []*Property, error) {

	list, err := t.getHouseList(stub, sellerID)

	if err != nil {
		return nil, 0, nil, err
	}

	if list != nil {
		for i, v := range list {
			if v.Hash == propertyID {
				return v, i, list, nil
			}
		}
		return nil, 0, list, errors.New("Could not found property")
	}

	return nil, 0, nil, errors.New("Could not found property")
}

func (t *HomelendChaincode) getMoney(stub shim.ChaincodeStubInterface, userID string) int {

	dataAsBytes, err := stub.GetState(userID)
	if err != nil {
		str := fmt.Sprintf("Could not getMoney userID=" + userID)
		fmt.Println(str)
		return -1
	}

	escrowMoney, err := strconv.Atoi(string(dataAsBytes))
	if err != nil {
		str := fmt.Sprintf("Could not getMoney - Atoi userID=" + userID)
		fmt.Println(str)
		return -1
	}

	return escrowMoney
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

	stub.PutState(srcUserID, []byte(srcMoneyStr))
	stub.PutState(destUserID, []byte(destMoneyStr))

	moneyTransferBytes, err := stub.GetState(MoneyTransferKey)

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
	stub.PutState(MoneyTransferKey, moneyTransferBytes)

	return nil
}

func (t *HomelendChaincode) getHouseList(stub shim.ChaincodeStubInterface, userID string) ([]*Property, error) {

	dataAsBytes, err := stub.GetState(userID)
	if err != nil {
		str := fmt.Sprintf("Failed to get Property from user: %s", userID)
		fmt.Println(str)
		return nil, err
	}

	if dataAsBytes != nil && len(dataAsBytes) > 0 {
		var arrayOfData []*Property
		err = json.Unmarshal(dataAsBytes, &arrayOfData)

		return arrayOfData, nil
	}

	return nil, nil
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
