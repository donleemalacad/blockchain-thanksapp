package integration

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// SampleDataUsingSDK demonstration of data retrieval using Fabric SDK
func (bridge *SdkSetup) SampleDataUsingSDK() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "retrieve")

	// Access func Query in Chaincode and pass necessary parameters
	query, err := bridge.channelClient.Query(channel.Request{
		ChaincodeID: bridge.ChaincodeName, 
		Fcn: args[0], Args: [][]byte{[]byte(args[1])},
	})

	// If Query is unsuccessful
	if err != nil {
		return "", fmt.Errorf("Failed to query: %v", err)
	}

	return string(query.Payload), nil
}