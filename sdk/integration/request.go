package integration

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"time"
)

// SampleDataUsingSDK demonstration of data retrieval using Fabric SDK
func (bridge *SdkSetup) AddPerson(name string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "addPerson")
	args = append(args, name)
	args = append(args, "System")

	eventID := "eventInvoke"

	reg, notifier, err := bridge.eventClient.RegisterChaincodeEvent(bridge.ChaincodeName, eventID)
	if err != nil {
		return "", err
	}

	defer bridge.eventClient.Unregister(reg)


	// Access func Query in Chaincode and pass necessary parameters
	query, err := bridge.channelClient.Execute(channel.Request{
		ChaincodeID: bridge.ChaincodeName, 
		Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])},
	})

	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("Chaincode Event failed for eventId(%s)", eventID)
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