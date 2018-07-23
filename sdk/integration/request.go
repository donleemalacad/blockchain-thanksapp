package integration

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// SampleDataUsingSDK demonstration of data retrieval using Fabric SDK
func (bridge *SdkSetup) SampleDataUsingSDK() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "addPerson")
	args = append(args, "Donlee Malacad")
	args = append(args, "System")

	// Access func Query in Chaincode and pass necessary parameters
	query, err := bridge.channelClient.Query(channel.Request{
		ChaincodeID: bridge.ChaincodeName, 
		Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])},
	})

	// If Query is unsuccessful
	if err != nil {
		return "", fmt.Errorf("Failed to query: %v", err)
	}

	return string(query.Payload), nil
}

// GetAllUserDetailsInLedger - Retrieve All History for Users
func (bridge *SdkSetup) GetAllUserDetailsInLedger() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "getAllUsers")

	// Access func Query in Chaincode and pass necessary parameters
	query, err := bridge.channelClient.Query(channel.Request{
		ChaincodeID: bridge.ChaincodeName, 
		Fcn: args[0], Args: [][]byte{},
	})

	// If Query is unsuccessful
	if err != nil {
		return "", fmt.Errorf("Failed to query: %v", err)
	}

	fmt.Println("-- Original Byte Value --")
	fmt.Println(query.Payload)

	return string(query.Payload), nil
}