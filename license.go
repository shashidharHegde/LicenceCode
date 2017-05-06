package main

import (
	"errors"
	"fmt"
	//"strconv"
	//"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	//"regexp"
)

var logger = shim.NewLogger("CLDChaincode")



//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type  SimpleChaincode struct {
}

//==============================================================================================================================
//	Vehicle - Defines the structure for a car object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type LIC struct {

	LICENCE_KEY 	string `json:"licence_key"`
	LICENCE_NAME 	string `json:"licence_name"`
	VENDOR_NAME 	string `json:"vendor_name"`
	Type_of_licence           string `json:"type_of_licence"`
	Allocation_date           string `json:"allocation_date"`
	Expiry_date           string `json:"expiry_date"`
	NUMBER_OF_LICENCE 	string `json:"no_of_licence"`
}

//==============================================================================================================================
//	Lic Holder - Defines the structure that holds all the License Name & Vendorthat have been created.
//				Used as an index when querying all vehicles.
//==============================================================================================================================

type Lic_Holder struct {
	Licholder 	[]string `json:"licholder"`
}
//==============================================================================================================================
//	User_and_eCert - Struct for storing the JSON of a user and their ecert
//==============================================================================================================================

//type User_and_eCert struct {
//	Identity string `json:"identity"`
//	eCert string `json:"ecert"`
//}

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var Lics Lic_Holder
	bytes, err := json.Marshal(Lics)

    if err != nil { return nil, errors.New("Error creating LIC_Holder record") }

	err = stub.PutState("Lics", bytes)

	for i:=0; i < len(args); i=i+2 {
		t.add_licarray(stub, args[i], args[i+1])
	}

	return nil, nil
}

//==============================================================================================================================
//	 General Functions
//==============================================================================================================================
//	 get_ecert - Takes the name passed and calls out to the REST API for HyperLedger to retrieve the ecert
//				 for that user. Returns the ecert as retrived including html encoding.
//==============================================================================================================================
func (t *SimpleChaincode) get_licarray(stub shim.ChaincodeStubInterface, name string) ([]byte, error) {

	holder, err := stub.GetState(name)

	if err != nil { return nil, errors.New("Couldn't retrieve ecert for user " + name) }

	return holder, nil
}

//==============================================================================================================================
//	 add_ecert - Adds a new ecert and user pair to the table of ecerts
//==============================================================================================================================

func (t *SimpleChaincode) add_licarray(stub shim.ChaincodeStubInterface, name string, holder string) ([]byte, error) {


	err := stub.PutState(name, []byte(holder))

	if err == nil {
		return nil, errors.New("Error storing holder for license " + name + " identity: " + holder)
	}

	return nil, nil

}
//=================================================================================================================================
//	Write - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================

func (t *SimpleChaincode) create_entry(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	//var err error
	var Lic LIC

	fmt.Println("running write()")

		if len(args) != 7 {
			return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
		}


	Lic.LICENCE_KEY = args[0] //rename for funsies
	Lic.LICENCE_NAME = args[1]
	Lic.VENDOR_NAME = args[2]
	Lic.Type_of_licence = args[3]
	Lic.Allocation_date = args[4]
	Lic.Expiry_date = args[5]
	Lic.NUMBER_OF_LICENCE = args[6]

	
	LicenseJSONasBytes, err := json.Marshal(Lic)
		if err != nil {
			return nil, err
		}
		
	// === Save employee to state ===
	err = stub.PutState(Lic.LICENCE_KEY, LicenseJSONasBytes)
		if err != nil {
			return nil, err
		}
	
	return nil, nil
}
//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "get_licarray" {
		bytes, err := t.get_licarray(stub, args[0])
		if err != nil {return nil, errors.New("Received unknown function invocation " + function)}
    	return bytes, nil
	}
	if function == "get_entry" {
		bytes, err := t.get_entry(stub, args[0])
		if err != nil {return nil, errors.New("Received unknown function invocation " + function)}
    	return bytes, nil
	}
	return nil, nil
}
//==============================================================================================================================
//   get license details
//==============================================================================================================================
func (t *SimpleChaincode) get_entry(stub shim.ChaincodeStubInterface, key string) ([]byte, error) {

	licentry, err := stub.GetState(key)

	if err != nil { return nil, errors.New("Couldn't retrieve ecert for user " + key) }

	return licentry, nil
	
	
}
//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
   	
	if function == "create_entry" {
         return t.create_entry(stub, args)
	}
	
    return nil, nil
}
//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
func main() {

	err := shim.Start(new(SimpleChaincode))

	if err != nil { fmt.Printf("Error starting Chaincode: %s", err) }
}
