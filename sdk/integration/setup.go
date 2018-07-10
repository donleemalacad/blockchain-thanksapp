package integration

import (
	"fmt"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

// SDK Integration Struct
type SdkSetup struct {
	ChaincodeName   string
	ChaincodePath   string
	ChannelName     string
	ChannelTxPath  	string
	GoPath          string
	Initialized     bool
	OrdererName     string
	PeerAdmin       string
	PeerName        string
	PeerUser        string
	SdkConfigFile   string
	sdk             *fabsdk.FabricSDK
	resClient       *resmgmt.Client
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (integrate *SdkSetup) Initialize() error {

	// Check if SDK is Initialized.
	if integrate.Initialized {
		return errors.New("SDK is already initialized")
	}

	// If SDK is not initialized, Initialize SDK.
	sdk, err := fabsdk.New(config.FromFile(integrate.SdkConfigFile))
	if err != nil {
		return errors.WithMessage(err, "Failed to initialize SDK")
	}
	integrate.sdk = sdk

	// At this point if it passes, SDK is Initialized.
	fmt.Println("SDK created")

	// Prepare Client Context
	clientContext := integrate.sdk.Context(fabsdk.WithUser(integrate.PeerAdmin), fabsdk.WithOrg(integrate.PeerName))
	if err != nil {
		return errors.WithMessage(err, "Failed to prepare client context")
	}

	// Create resource management client.
	resourceManagementClient, err := resmgmt.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "Failed to create resource management client")
	}
	integrate.resClient = resourceManagementClient

	// If it passes this point, Resource management client is created.
	fmt.Println("Resource Management Client is Created")

	// The MSP client allow us to retrieve user information which we need to save channel
	// 1. Initialize MSP Client
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(integrate.PeerName))
	if err != nil {
		return errors.WithMessage(err, "Failed to create MSP client")
	}

	// 2. Get Signing Identity
	adminIdentity, err := mspClient.GetSigningIdentity(integrate.PeerAdmin)
	if err != nil {
		return errors.WithMessage(err, "Failed to get Signing Identity")
	}

	// 3. Save Channel Request
	req := resmgmt.SaveChannelRequest{ChannelID: integrate.ChannelName, ChannelConfigPath: integrate.ChannelTxPath, SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	channelTx, err := integrate.resClient.SaveChannel(req, resmgmt.WithOrdererEndpoint(integrate.OrdererName))
	if err != nil || channelTx.TransactionID == "" {
		return errors.WithMessage(err, "Failed to save channel")
	}

	// If it passes this point, Channel is created
	fmt.Println("Successfully created channel")

	// Join admin to newly created channel
	if err = integrate.resClient.JoinChannel(integrate.ChannelName, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(integrate.OrdererName)); err != nil {
		return errors.WithMessage(err, "Failed to join channel")
	}

	// If it passes this point, Admin was able to join the channel
	fmt.Println("Admin successfully joined channel")

	// If it passes this point, SDK Initialization is successful
	fmt.Println("SDK Initialization Successful")
	integrate.Initialized = true
	return nil
}

func (integrate *SdkSetup) CloseSDK() {
	integrate.sdk.Close()
}
