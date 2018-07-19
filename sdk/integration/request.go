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
	args = append(args, "Nobuyuki Mugima")
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

func (bridge *SdkSetup) GetAllBlocks() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "getHistoryOfPerson")
	args = append(args, "Nobuyuki Mugima")

	// Access func Query in Chaincode and pass necessary parameters
	query, err := bridge.channelClient.Query(channel.Request{
		ChaincodeID: bridge.ChaincodeName, 
		Fcn: args[0], Args: [][]byte{[]byte(args[1])},
	})

	// If Query is unsuccessful
	if err != nil {
		return "", fmt.Errorf("Failed to query: %v", err)
	}

	fmt.Println(query.Payload)
	return string(query.Payload), nil
}