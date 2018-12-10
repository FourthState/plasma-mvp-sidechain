package eth

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/tendermint/tendermint/libs/log"
)

// Client defines wrappers to a remote endpoint
type Client struct {
	rpc    *rpc.Client
	ec     *ethclient.Client
	logger log.Logger
}

// Instantiate a connection and bind the go plasma contract wrapper with this client
func InitEthConn(nodeUrl string, logger log.Logger) (*Client, error) {
	// Connect to a remote etheruem client
	//
	// Ethclient wraps around the underlying rpc module and provides convenient functions. We still keep reference
	// to the underlying rpc module to make calls that the wrapper does not support
	c, err := rpc.Dial(nodeUrl)
	if err != nil {
		return nil, err
	}
	ec := ethclient.NewClient(c)

	return &Client{c, ec, logger}, nil
}

// SubscribeToHeads returns a channel that funnels new ethereum headers to the returned channel
func (client *Client) SubscribeToHeads() (<-chan *types.Header, error) {
	c := make(chan *types.Header)

	sub, err := client.ec.SubscribeNewHead(context.Background(), c)
	if err != nil {
		return nil, err
	}

	// close the channel if an error arises in the subscription
	go func() {
		for {
			err = <-sub.Err()
			client.logger.Error("Etheruem client header subscription error -", err)
			close(c)
		}
	}()

	return c, nil
}

// used for testing when running against a local client like ganache
func (client *Client) accounts() ([]common.Address, error) {
	var res json.RawMessage
	err := client.rpc.Call(&res, "eth_accounts")
	if err != nil {
		return nil, err
	}

	var addrs []string
	if err := json.Unmarshal(res, &addrs); err != nil {
		return nil, err
	}

	// convert to the correct type
	result := make([]common.Address, len(addrs))
	for i, addr := range addrs {
		result[i] = common.HexToAddress(addr)
	}

	return result, nil
}
