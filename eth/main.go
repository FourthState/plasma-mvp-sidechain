package eth

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
)

// Client defines wrappers to a remote endpoint
type Client struct {
	rpc    *rpc.Client
	ec     *ethclient.Client
	logger log.Logger
}

// Instantiate a connection and bind the go plasma contract wrapper with this client
func InitEthConn(nodeUrl string, logger log.Logger) (Client, error) {
	// Connect to a remote etheruem client
	//
	// Ethclient wraps around the underlying rpc module and provides convenient functions. We still keep reference
	// to the underlying rpc module to make calls that the wrapper does not support
	c, err := rpc.Dial(nodeUrl)
	if err != nil {
		return Client{}, err
	}
	ec := ethclient.NewClient(c)

	return Client{c, ec, logger}, nil
}

func (client Client) LatestBlockNum() (*big.Int, error) {
	var hexStr string
	if err := client.rpc.Call(&hexStr, "eth_blockNumber"); err != nil {
		return nil, fmt.Errorf("RPC error { %s }", err)
	}

	hexStr = utils.RemoveHexPrefix(hexStr)

	// pad if hex length is odd
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	hexBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("Hex decoding error: { %s }", err)
	}

	return new(big.Int).SetBytes(hexBytes), nil
}

// SubscribeToHeads returns a channel that funnels new ethereum headers to the returned channel
func (client Client) SubscribeToHeads() (<-chan *types.Header, error) {
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
func (client Client) accounts() ([]common.Address, error) {
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
