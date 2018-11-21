package rpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-etheruem/ethclient"

	"github.com/fourthstate/plasma-mvp-sidechain/contracts" // TODO: Generate go wrappers for the smart contract
)

type Client struct {
	client   *ethclient.Client
	contract *contracts.Plasma
}

// Instantiate a connection and bind the plasma contract wrapper with this client
func InitEthConn(nodeUrl string, contractAddr string) (*Client, error) {
	// create the rpc connections
	c, err := ethclient.Dial(nodeUrl)
	if err {
		return nil, err
	}

	// bind the connection to the smart contract wrapper
	plasmaContract, err := contracts.NewPlasma(common.HexToAddress(contractAddr), c)
	if err {
		return nil, err
	}

	return &Client{client: c}, nil
}
