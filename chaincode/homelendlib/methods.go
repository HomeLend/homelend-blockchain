package lib

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//Helpers - Helps us Getting values from database
type Helpers struct {
}

//GetLastRequest -Get the request of the user by userID
func (t *Helpers) GetLastRequest(stub shim.ChaincodeStubInterface, userID string) (*Request, []*Request, error) {
	valAsBytes, err := stub.GetState(DbRequests + userID)

	if err != nil {
		str := fmt.Sprintf("Failed to get state %+v", err.Error())
		fmt.Println(str)
		return nil, nil, err
	} else if valAsBytes == nil {
		str := fmt.Sprintf("Request record does not exist %s", userID)
		fmt.Println(str)
		return nil, nil, err
	}

	var arrayOfData []*Request
	err = json.Unmarshal(valAsBytes, &arrayOfData)

	if err != nil {
		str := fmt.Sprintf("Failed to unmarshal: %s", err)
		fmt.Println(str)
		return nil, nil, err
	}

	latestRequest := arrayOfData[len(arrayOfData)-1]
	return latestRequest, arrayOfData, nil
}

//UpdateLastRequest -Updates the request of the user by userID
func (t *Helpers) UpdateLastRequest(stub shim.ChaincodeStubInterface, userID string, userRequestsArray []*Request) error {

	dataJSONasBytes, err := json.Marshal(userRequestsArray)
	if err != nil {
		str := fmt.Sprintf("Could not marshal %+v", err.Error())
		fmt.Println(str)
		return err
	}

	err = stub.PutState(DbRequests+userID, dataJSONasBytes)
	return err
}

//PrintAndReturnError - prints and return error
func (t *Helpers) PrintAndReturnError(stub shim.ChaincodeStubInterface, errorStr string) pb.Response {
	str := fmt.Sprintf(errorStr)
	fmt.Println(str)
	return shim.Error(str)
}
