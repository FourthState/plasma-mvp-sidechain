package eth

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Client defines wrappers to a remote endpoint
type Client struct {
	rpc *rpc.Client
	ec  *ethclient.Client
}

// Instantiate a connection and bind the go plasma contract wrapper with this client
func InitEthConn(nodeUrl string) (*Client, error) {
	// Connect to a remote etheruem client
	//
	// Ethclient wraps around the underlying rpc module and provides convenient functions. We still keep reference
	// to the underlying rpc module to make calls that the wrapper does not support
	c, err := rpc.Dial(nodeUrl)
	if err != nil {
		return nil, err
	}
	ec := ethclient.NewClient(c)

	return &Client{c, ec}, nil
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
