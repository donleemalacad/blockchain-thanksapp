package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ThanksCcChaincode implementation of Chaincode
type ThanksCcChaincode struct {
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *ThanksCcChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("************************************************")
	fmt.Println("*        ThanksCcChaincode Initialize          *")
	fmt.Println("************************************************")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}

	// Put in the ledger the key/value hello/world
	err := stub.PutState("foo", []byte("This is a sample baz response"))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return a successful message
	return shim.Success(nil)
}

// Invoke - Request Invocation
func (t *ThanksCcChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("************************************************")
	fmt.Println("*         ThanksCcChaincode Invoke             *")
	fmt.Println("************************************************")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	// Check if number of args is 2
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}
	
	if args[0] == "retrieve" {
		// Get the foo key values
		state, err := stub.GetState("foo")
		if err != nil {
			return shim.Error("Failed to get state of hello")
		}

		// Return this value in response
		return shim.Success(state)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(ThanksCcChaincode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}