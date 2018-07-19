package main

import (
	"fmt"
	"github.com/thanksapp/sdk/integration"
	"github.com/thanksapp/web"
	"github.com/thanksapp/web/controllers"
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
		ChaincodeName: "thankscc",
		ChaincodePath: "github.com/thanksapp/chaincode/go/",
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

	err = setup.InvokeChaincode()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}

	// Query the chaincode
	response, err := setup.SampleDataUsingSDK()
	if err != nil {
		fmt.Printf("Unable to query foo on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the query foo: %s\n", response)
	}

	// Setup WebApp
	app := &controllers.Application {
		Fabric: &setup,
	}
	web.Serve(app)
}