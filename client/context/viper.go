package context

import (
	"github.com/spf13/viper"

	rpcclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/FourthState/plasma-mvp-sidechain/client"
)

// Return a new context with parameters from the command line
func NewClientContextFromViper() ClientContext {
	nodeURI := viper.GetString(client.FlagNode)
	var rpc rpcclient.Client
	if nodeURI != "" {
		rpc = rpcclient.NewHTTP(nodeURI, "/websocket")
	}
	return ClientContext{
		Height:         viper.GetInt64(client.FlagHeight),
		TrustNode:      viper.GetBool(client.FlagTrustNode),
		InputAddresses: viper.GetString(client.FlagAddress),
		NodeURI:        nodeURI,
		Client:         rpc,
		Decoder:        nil,
		UTXOStore:      "main",
	}
}
