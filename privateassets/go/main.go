package main

import (
	"bytes"
	"encoding/json"
	"fmt"
        "os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type AssetsPrivateChaincode struct {
}

type asset struct {
	ObjectType string `json:"docType"` 
	Name       string `json:"name"`    
	Color      string `json:"color"`
	Size       int    `json:"size"`
	Owner      string `json:"owner"`
}

type assetPrivateDetails struct {
	ObjectType string `json:"docType"` 
	Name       string `json:"name"`    
	Price      int    `json:"price"`
}

// initializes chaincode
func (t *AssetsPrivateChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke 
func (t *AssetsPrivateChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	switch function {
	case "initAsset":
		return t.initAsset(stub, args)
	case "readAsset":
		return t.readAsset(stub, args)
	case "readAssetPrivateDetails":
		return t.readAssetPrivateDetails(stub, args)
	case "transferAsset":
		return t.transferAsset(stub, args)
	case "delete":
		return t.delete(stub, args)
	case "getAssetsByRange":
		return t.getAssetsByRange(stub, args)
	case "getAssetHash":
		return t.getAssetHash(stub, args)
	case "getAssetPrivateDetailsHash":
		return t.getAssetPrivateDetailsHash(stub, args)
	default:
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

//initAsset
func (t *AssetsPrivateChaincode) initAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	type assetTransientInput struct {
		Name  string `json:"name"` 
		Color string `json:"color"`
		Size  int    `json:"size"`
		Owner string `json:"owner"`
		Price int    `json:"price"`
	}

	fmt.Println("- start init asset")

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private asset data must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	assetJsonBytes, ok := transMap["asset"]
	if !ok {
		return shim.Error("asset must be a key in the transient map")
	}

	if len(assetJsonBytes) == 0 {
		return shim.Error("asset value in the transient map must be a non-empty JSON string")
	}

	var assetInput assetTransientInput
	err = json.Unmarshal(assetJsonBytes, &assetInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(assetJsonBytes))
	}

	if len(assetInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(assetInput.Color) == 0 {
		return shim.Error("color field must be a non-empty string")
	}
	if assetInput.Size <= 0 {
		return shim.Error("size field must be a positive integer")
	}
	if len(assetInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}
	if assetInput.Price <= 0 {
		return shim.Error("price field must be a positive integer")
	}

	assetAsBytes, err := stub.GetPrivateData("collectionAssets", assetInput.Name)
	if err != nil {
		return shim.Error("Failed to get asset: " + err.Error())
	} else if assetAsBytes != nil {
		fmt.Println("This asset already exists: " + assetInput.Name)
		return shim.Error("This asset already exists: " + assetInput.Name)
	}

	asset := &asset{
		ObjectType: "asset",
		Name:       assetInput.Name,
		Color:      assetInput.Color,
		Size:       assetInput.Size,
		Owner:      assetInput.Owner,
	}
	assetJSONasBytes, err := json.Marshal(asset)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutPrivateData("collectionAssets", assetInput.Name, assetJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	assetPrivateDetails := &assetPrivateDetails{
		ObjectType: "assetPrivateDetails",
		Name:       assetInput.Name,
		Price:      assetInput.Price,
	}
	assetPrivateDetailsBytes, err := json.Marshal(assetPrivateDetails)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutPrivateData("collectionAssetPrivateDetails", assetInput.Name, assetPrivateDetailsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "color~name"
	colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{asset.Color, asset.Name})
	if err != nil {
		return shim.Error(err.Error())
	}

	value := []byte{0x00}
	stub.PutPrivateData("collectionAssets", colorNameIndexKey, value)

	fmt.Println("- end init marble")
	return shim.Success(nil)
}

func (t *AssetsPrivateChaincode) readAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the asset to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionAssets", name) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Asset does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *AssetsPrivateChaincode) readAssetPrivateDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the asset to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionAssetPrivateDetails", name) 
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get private details for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Asset private details does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *AssetsPrivateChaincode) getAssetHash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the Asset to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateDataHash("collectionAssets", name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get asset private data hash for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Asset private asset data hash does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *AssetsPrivateChaincode) getAssetPrivateDetailsHash(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the asset to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateDataHash("collectionAssetPrivateDetails", name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get asset private details hash for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Asset private details hash does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *AssetsPrivateChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start delete asset")

	type assetDeleteTransientInput struct {
		Name string `json:"name"`
	}

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private asset name must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	assetDeleteJsonBytes, ok := transMap["asset_delete"]
	if !ok {
		return shim.Error("asset_delete must be a key in the transient map")
	}

	if len(assetDeleteJsonBytes) == 0 {
		return shim.Error("asset_delete value in the transient map must be a non-empty JSON string")
	}

	var assetDeleteInput assetDeleteTransientInput
	err = json.Unmarshal(assetDeleteJsonBytes, &assetDeleteInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(assetDeleteJsonBytes))
	}

	if len(assetDeleteInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}

	valAsbytes, err := stub.GetPrivateData("collectionAssets", assetDeleteInput.Name) //get the marble from chaincode state
	if err != nil {
		return shim.Error("Failed to get state for " + assetDeleteInput.Name)
	} else if valAsbytes == nil {
		return shim.Error("Asset does not exist: " + assetDeleteInput.Name)
	}

	var assetToDelete asset
	err = json.Unmarshal([]byte(valAsbytes), &assetToDelete)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(valAsbytes))
	}

	err = stub.DelPrivateData("collectionAssets", assetDeleteInput.Name)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	indexName := "color~name"
	colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{assetToDelete.Color, assetToDelete.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.DelPrivateData("collectionAssets", colorNameIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	err = stub.DelPrivateData("collectionAssetPrivateDetails", assetDeleteInput.Name)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *AssetsPrivateChaincode) transferAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("- start transfer asset")

	type assetTransferTransientInput struct {
		Name  string `json:"name"`
		Owner string `json:"owner"`
	}

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private asset data must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	assetOwnerJsonBytes, ok := transMap["asset_owner"]
	if !ok {
		return shim.Error("asset_owner must be a key in the transient map")
	}

	if len(assetOwnerJsonBytes) == 0 {
		return shim.Error("asset_owner value in the transient map must be a non-empty JSON string")
	}

	var assetTransferInput assetTransferTransientInput
	err = json.Unmarshal(assetOwnerJsonBytes, &assetTransferInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(assetOwnerJsonBytes))
	}

	if len(assetTransferInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(assetTransferInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}

	assetAsBytes, err := stub.GetPrivateData("collectionAssets", assetTransferInput.Name)
	if err != nil {
		return shim.Error("Failed to get asset:" + err.Error())
	} else if assetAsBytes == nil {
		return shim.Error("Asset does not exist: " + assetTransferInput.Name)
	}

	assetToTransfer := asset{}
	err = json.Unmarshal(assetAsBytes, &assetToTransfer) 
	if err != nil {
		return shim.Error(err.Error())
	}
	assetToTransfer.Owner = assetTransferInput.Owner 

	assetJSONasBytes, _ := json.Marshal(assetToTransfer)
	err = stub.PutPrivateData("collectionAssets", assetToTransfer.Name, assetJSONasBytes) 
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end transferAsset (success)")
	return shim.Success(nil)
}

func (t *AssetsPrivateChaincode) getAssetsByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetPrivateDataByRange("collectionAssets", startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

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

	fmt.Printf("- getAssetsByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(&AssetsPrivateChaincode{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exiting Simple chaincode: %s", err)
		os.Exit(2)
	}
}
