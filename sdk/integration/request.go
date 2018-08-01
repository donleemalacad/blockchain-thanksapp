package integration

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// AddPerson - Add another user to the ledger
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
		Fcn:         args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])},
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
		Fcn:         args[0], Args: [][]byte{},
	})

	// If Query is unsuccessful
	if err != nil {
		return "", fmt.Errorf("Failed to query: %v", err)
	}

	fmt.Println("-- Original Byte Value --")
	fmt.Println(query.Payload)

	return string(query.Payload), nil
}

// GetSpecificUserDetails - Retrieve History for Specified User
func (bridge *SdkSetup) GetSpecificUserDetails(user string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "getHistoryOfPerson")
	args = append(args, user)

	// Access func Query in Chaincode and pass necessary parameters
	query, err := bridge.channelClient.Query(channel.Request{
		ChaincodeID: bridge.ChaincodeName,
		Fcn:         args[0], Args: [][]byte{[]byte(user)},
	})

	// If Query is unsuccessful
	if err != nil {
		return "", fmt.Errorf("Failed to query: %v", err)
	}

	return string(query.Payload), nil
}

// ReplenishPoints - adds points to everyone every month by cron
func (bridge *SdkSetup) ReplenishPoints() (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "addPointToAll")

	eventID := "eventInvoke"
	fmt.Printf("replenishpoints request")
	reg, notifier, err := bridge.eventClient.RegisterChaincodeEvent(bridge.ChaincodeName, eventID)
	if err != nil {
		return "", err
	}

	defer bridge.eventClient.Unregister(reg)

	// Access func Query in Chaincode and pass necessary parameters
	query, err := bridge.channelClient.Execute(channel.Request{
		ChaincodeID: bridge.ChaincodeName,
		Fcn:         args[0], Args: [][]byte{},
	})
	fmt.Printf("query")
	fmt.Print(query)
	fmt.Printf("err")
	fmt.Println(err)

	if err != nil {
		fmt.Println(err)

		return "replenish err", err
	}

	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("Chaincode Event failed for eventId(%s)", eventID)
	}

	return string(query.Payload), nil
}

// TransferPoint - Point Transferring Between two Users
func (bridge *SdkSetup) TransferPoint(from string, to string, message string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "transfer")
	args = append(args, from)
	args = append(args, to)
	args = append(args, message)

	eventID := "eventInvoke"

	reg, notifier, err := bridge.eventClient.RegisterChaincodeEvent(bridge.ChaincodeName, eventID)
	if err != nil {
		return "", err
	}

	defer bridge.eventClient.Unregister(reg)

	// Access func Query in Chaincode and pass necessary parameters
	query, err := bridge.channelClient.Execute(channel.Request{
		ChaincodeID: bridge.ChaincodeName,
		Fcn:         args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])},
	})
	fmt.Print("query")
	fmt.Print(query)
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("Chaincode Event failed for eventId(%s)", eventID)
	}

	return string(query.Payload), nil
}
