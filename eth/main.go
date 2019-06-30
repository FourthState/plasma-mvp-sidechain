package eth

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

// Client defines wrappers to a remote endpoint
type Client struct {
	rpc *rpc.Client
	ec  *ethclient.Client
}

// Instantiate a connection and bind the go plasma contract wrapper with this client.
// Will return an error if the ethereum client is not fully sycned
func InitEthConn(nodeUrl string) (Client, error) {
	// Connect to a remote etheruem client
	//
	// Ethclient wraps around the underlying rpc module and provides convenient functions. We still keep reference
	// to the underlying rpc module to make calls that the wrapper does not support
	c, err := rpc.Dial(nodeUrl)
	if err != nil {
		return Client{}, err
	}
	ec := ethclient.NewClient(c)

	// check if the client is synced
	client := Client{c, ec}
	if synced, err := client.Synced(); !synced || err != nil {
		if err != nil {
			return client, err
		} else {
			return client, fmt.Errorf("geth endpoint is not fully synced")
		}
	}

	return client, nil
}

// Synced checks of the status of the geth endpoint with it's network
func (client Client) Synced() (bool, error) {
	var res json.RawMessage
	if err := client.rpc.Call(&res, "eth_syncing"); err != nil {
		return false, fmt.Errorf("rpc: %s", err)
	}

	if string(res) != "false" {
		return false, nil
	}

	return true, nil
}

// LatestBlockNum retrieves the latest block height of the of the geth endpoint
func (client Client) LatestBlockNum() (*big.Int, error) {
	var hexStr string
	if err := client.rpc.Call(&hexStr, "eth_blockNumber"); err != nil {
		return nil, fmt.Errorf("rpc: %s", err)
	}

	hexStr = utils.RemoveHexPrefix(hexStr)

	// pad if hex length is odd
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	hexBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("hex: %s", err)
	}

	return new(big.Int).SetBytes(hexBytes), nil
}
