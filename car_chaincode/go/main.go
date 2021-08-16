package main

import (
	"bytes"
	"encoding/json"
	"fmt"
        "os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// Car Chaincode implementation
type CarPrivateChaincode struct {
}

type car struct {
	ObjectType string `json:"docType"` 
	Name       string `json:"name"`    
	Color      string `json:"color"`
	TireSize       int    `json:"tiresize"`
	Owner      string `json:"owner"`
}

// defnining private_data_details
type carPrivateDetails struct {
	ObjectType string `json:"docType"` 
	Name       string `json:"name"`    
	Price      int    `json:"price"`
}

// Init function
func (t *CarPrivateChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke functions
func (t *CarPrivateChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	switch function {
	case "initCar":
		return t.initCar(stub, args)
	case "readCar":
		return t.readCar(stub, args)
	case "readCarPrivateDetails":
		return t.readCarPrivateDetails(stub, args)
	case "transferCar":
		return t.transferCar(stub, args)
	case "delete":
		return t.delete(stub, args)
	case "getCarByRange":
		return t.getCarByRange(stub, args)
	case "getCarHash":
		return t.getCarHash(stub, args)
	case "getCarPrivateDetailsHash":
		return t.getCarPrivateDetailsHash(stub, args)
	default:
		//error
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}
// initCar - create a new car, store into chaincode state
func (t *CarPrivateChaincode) initCar(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	type carTransientInput struct {
		Name  string `json:"name"` 
		Color string `json:"color"`
		TireSize  int    `json:"tiresize"`
		Owner string `json:"owner"`
		Price int    `json:"price"`
	}

	fmt.Println("- start init car asset")

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	carJsonBytes, ok := transMap["car"]
	if !ok {
		return shim.Error("car must be a key in the transient map")
	}

	if len(carJsonBytes) == 0 {
		return shim.Error("car value in the transient map must be a non-empty JSON string")
	}

	var carInput carTransientInput
	err = json.Unmarshal(carJsonBytes, &carInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(carJsonBytes))
	}

	if len(carInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(carInput.Color) == 0 {
		return shim.Error("color field must be a non-empty string")
	}
	if carInput.TireSize <= 0 {
		return shim.Error("TireSize field must be a positive integer")
	}
	if len(carInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}
	if carInput.Price <= 0 {
		return shim.Error("price field must be a positive integer")
	}

// check if car already exists
	carAsBytes, err := stub.GetPrivateData("collectionCar", carInput.Name)
	if err != nil {
		return shim.Error("Failed to get car: " + err.Error())
	} else if carAsBytes != nil {
		fmt.Println("This car already exists: " + carInput.Name)
		return shim.Error("This car already exists: " + carInput.Name)
	}

	// Create car object marshalling
	car := &car{
		ObjectType: "car",
		Name:       carInput.Name,
		Color:      carInput.Color,
		TireSize:       carInput.TireSize,
		Owner:      carInput.Owner,
	}
	carJSONasBytes, err := json.Marshal(car)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Saving car to state 
	err = stub.PutPrivateData("collectionCar", carInput.Name, carJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Create car private details object for price
	carPrivateDetails := &carPrivateDetails{
		ObjectType: "carPrivateDetails",
		Name:       carInput.Name,
		Price:      carInput.Price,
	}
	carPrivateDetailsBytes, err := json.Marshal(carPrivateDetails)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutPrivateData("collectionCarPrivateDetails", carInput.Name, carPrivateDetailsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

 // Indexing for query optimizations 
	indexName := "color~name"
	colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{car.Color, car.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutPrivateData("collectionCar", colorNameIndexKey, value)

	fmt.Println("- end init car")
	return shim.Success(nil)
}

// Reads a car asset from saved state
func (t *CarPrivateChaincode) readCar(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the car to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionCar", name) 
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Car does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// Reads private details information about the car
func (t *CarPrivateChaincode) readCarPrivateDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the car to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionCarPrivateDetails", name) 
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get private details for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Car private details does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}


// getCarHash - getting carhashvalue for PDC

func (t *CarPrivateChaincode) getCarHash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the car to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateDataHash("collectionCar", name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get car private data hash for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Car private car data hash does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *CarPrivateChaincode) getCarPrivateDetailsHash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the car to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateDataHash("collectionCarPrivateDetails", name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get car private details hash for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Car private details hash does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// Deletes a car

func (t *CarPrivateChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start delete car")

	type carDeleteTransientInput struct {
		Name string `json:"name"`
	}

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private car name must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	carDeleteJsonBytes, ok := transMap["car_delete"]
	if !ok {
		return shim.Error("car_delete must be a key in the transient map")
	}

	if len(carDeleteJsonBytes) == 0 {
		return shim.Error("car_delete value in the transient map must be a non-empty JSON string")
	}

	var carDeleteInput carDeleteTransientInput
	err = json.Unmarshal(carDeleteJsonBytes, &carDeleteInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(carDeleteJsonBytes))
	}

	if len(carDeleteInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}

	// to maintain the color~name index, we need to read the car first and get its color
	valAsbytes, err := stub.GetPrivateData("collectionCar", carDeleteInput.Name) //get the car from chaincode state
	if err != nil {
		return shim.Error("Failed to get state for " + carDeleteInput.Name)
	} else if valAsbytes == nil {
		return shim.Error("Car does not exist: " + carDeleteInput.Name)
	}

	var carToDelete car
	err = json.Unmarshal([]byte(valAsbytes), &carToDelete)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(valAsbytes))
	}

	// delete the car from state
	err = stub.DelPrivateData("collectionCar", carDeleteInput.Name)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	// Also delete the car from the color~name index
	indexName := "color~name"
	colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{carToDelete.Color, carToDelete.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.DelPrivateData("collectionCar", colorNameIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	// Finally, delete private details of car
	err = stub.DelPrivateData("collectionCarPrivateDetails", carDeleteInput.Name)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ===========================================================
// transfer a car by setting a new owner name on the car
// ===========================================================
func (t *CarPrivateChaincode) transferCar(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("- start transfer car")

	type carTransferTransientInput struct {
		Name  string `json:"name"`
		Owner string `json:"owner"`
	}

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private car data must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	carOwnerJsonBytes, ok := transMap["car_owner"]
	if !ok {
		return shim.Error("car_owner must be a key in the transient map")
	}

	if len(carOwnerJsonBytes) == 0 {
		return shim.Error("car_owner value in the transient map must be a non-empty JSON string")
	}

	var carTransferInput carTransferTransientInput
	err = json.Unmarshal(carOwnerJsonBytes, &carTransferInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(carOwnerJsonBytes))
	}

	if len(carTransferInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(carTransferInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}

	carAsBytes, err := stub.GetPrivateData("collectionCar", carTransferInput.Name)
	if err != nil {
		return shim.Error("Failed to get car:" + err.Error())
	} else if carAsBytes == nil {
		return shim.Error("Car does not exist: " + carTransferInput.Name)
	}

	carToTransfer := car{}
	err = json.Unmarshal(carAsBytes, &carToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	carToTransfer.Owner = carTransferInput.Owner //change the owner

	carJSONasBytes, _ := json.Marshal(carToTransfer)
	err = stub.PutPrivateData("collectionCar", carToTransfer.Name, carJSONasBytes) 
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end transferCar (success)")
	return shim.Success(nil)
}

func (t *CarPrivateChaincode) getCarByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetPrivateDataByRange("collectionCar", startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten {
			buffer.WriteString(",")
		}

		buffer.WriteString(
			fmt.Sprintf(
				`{"Key":"%s", "Record":%s}`,
				queryResponse.Key, queryResponse.Value,
			),
		)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getCarByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(&CarPrivateChaincode{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exiting Simple chaincode: %s", err)
		os.Exit(2)
	}
}
