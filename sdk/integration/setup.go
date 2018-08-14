package integration

import (
	"fmt"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	packager  "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
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
	channelClient   *channel.Client
	eventClient     *event.Client
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (integrate *SdkSetup) Initialize(restartKey int) error {

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

	if (restartKey == 0) {
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
	}

	// If it passes this point, SDK Initialization is successful
	fmt.Println("SDK Initialization Successful")
	integrate.Initialized = true
	return nil
}

// Invoke Chaincode
func (integrate *SdkSetup) InvokeChaincode(restartKey int) error {
	if (restartKey == 0) {
		// Create chaincode package
		chaincodePackage, err := packager.NewCCPackage(integrate.ChaincodePath, integrate.GoPath)
		if err != nil {
			return errors.WithMessage(err, "Failed to create chaincode package")
		}
		fmt.Println("Chaincode Package created")

		// Install chaincode to peers
		chaincodeProper := resmgmt.InstallCCRequest{Name: integrate.ChaincodeName, Path: integrate.ChaincodePath, Version: "0", Package: chaincodePackage}
		_, err = integrate.resClient.InstallCC(chaincodeProper, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			return errors.WithMessage(err, "Failed to install chaincode")
		}
		fmt.Println("Chaincode installed")

		// Setup Chaincode Policy using peer id
		chaincodePolicy := cauthdsl.SignedByAnyMember([]string{"org1.hf.excite.ph"})

		// Org resource manager will instantiate chaincode on channel
		resp, err := integrate.resClient.InstantiateCC(integrate.ChannelName, resmgmt.InstantiateCCRequest{ Name: integrate.ChaincodeName, Path: integrate.GoPath, Version: "0", Args: [][]byte{[]byte("init")}, Policy: chaincodePolicy, },)

		if err != nil || resp.TransactionID == "" {
			return errors.WithMessage(err, "Failed to instantiate the chaincode")
		}
		
		fmt.Println("Chaincode instantiated")
	}

	// Create client context
	var err error
	clientContext := integrate.sdk.ChannelContext(integrate.ChannelName, fabsdk.WithUser(integrate.PeerUser))
	integrate.channelClient, err = channel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new channel client")
	}
	fmt.Println("Channel client created")

	// Event Client
	integrate.eventClient, err = event.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new event client")
	}
	fmt.Println("Event client created")

	fmt.Println("Chaincode Installation & Instantiation Successful")

	return nil
}

// Close SDK
func (integrate *SdkSetup) CloseSDK() {
	integrate.sdk.Close()
}
