package context

import (
	"github.com/FourthState/plasma-mvp-sidechain/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// redefine in utxo.go
type UTXODecoder func(utxoBytes []byte) (types.UTXO, error)

type ClientContext struct {
	Height          int64
	TrustNode       bool
	NodeURI         string
	FromAddressName string
	Client          rpcclient.Client
	Decoder         UTXODecoder
	UTXOStore       string
}

// Returns a copy of the context with an updated height
func (c ClientContext) WithHeight(height int64) ClientContext {
	c.Height = height
	return c
}

// Returns a copy of the context with an updated TrustNode flag
func (c ClientContext) WithTrustNode(trustNode bool) ClientContext {
	c.TrustNode = trustNode
	return c
}

// Returns a copy of the xontext with an updated node URI
func (c ClientContext) WithNodeURI(nodeURI string) ClientContext {
	c.NodeURI = nodeURI
	c.Client = rpcclient.NewHTTP(nodeURI, "/websocket")
	return c
}

// Returns a copy of the context with an updated from address
func (c ClientContext) WithFromAddressName(fromAddressName string) ClientContext {
	c.FromAddressName = fromAddressName
	return c
}

// Returns a copy of the context with an updated RPC client instance
func (c ClientContext) WithClient(client rpcclient.Client) ClientContext {
	c.Client = client
	return c
}

// Returns a copy of the context with an updated utxo decoder
func (c ClientContext) WithDecoder(decoder UTXODecoder) ClientContext {
	c.Decoder = decoder
	return c
}

// Returns a copy of the context with an updated UTXOStore
func (c ClientContext) WithUTXOStore(utxoStore string) ClientContext {
	c.UTXOStore = utxoStore
	return c
}
