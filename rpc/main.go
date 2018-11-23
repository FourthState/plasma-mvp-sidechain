package rpc

import (
	//"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	//"github.com/fourthstate/plasma-mvp-sidechain/contracts"
)

type Client struct {
	client *ethclient.Client
	//contract *contracts.Plasma
}

// Instantiate a connection and bind the go plasma contract wrapper with this client
// TODO: add contract address
func InitEthConn(nodeUrl string) (*Client, error) {
	// create the rpc connections
	c, err := ethclient.Dial(nodeUrl)
	if err != nil {
		return nil, err
	}

	// TODO: use the abigen tool to create these wrappers in the contract submodule
	/*
		// bind the connection to the smart contract wrapper
		plasmaContract, err := contracts.NewPlasma(common.HexToAddress(contractAddr), c)
		if err {
			return nil, err
		}
	*/

	return &Client{client: c}, nil
}
