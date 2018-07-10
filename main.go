package main

import (
	"fmt"
	"github.com/thanksapp/sdk/integration"
	"os"
)

func main() {
	setup := integration.SdkSetup{
		// Network parameters 
		OrdererName: "orderer.hf.excite.ph",

		// Channel parameters
		ChannelName:   "thanksapp",
		ChannelTxPath: os.Getenv("GOPATH") + "/src/github.com/thanksapp/fixtures/artifacts/channel.tx",

		// Chaincode parameters
		ChaincodeName: "thanksapp",
		ChaincodePath: "github.com/thanksapp/chaincode/",
		GoPath:        os.Getenv("GOPATH"),
		PeerAdmin:     "Admin",
		PeerName:      "org1",
		SdkConfigFile: "sdkConfig.yaml",

		// User parameters
		PeerUser: "User1",
	}

	// Initialization of the Fabric SDK
	err := setup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer setup.CloseSDK()
}