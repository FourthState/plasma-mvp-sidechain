package eth

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
)

// Client defines wrappers to a remote endpoint
type Client struct {
	rpc    *rpc.Client
	ec     *ethclient.Client
	logger log.Logger
}

// Instantiate a connection and bind the go plasma contract wrapper with this client.
// Will return an error if the ethereum client is not fully sycned
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

	// check if the client is synced
	client := Client{c, ec, logger}
	if synced, err := client.Synced(); !synced || err != nil {
		if err != nil {
			return client, err
		} else {
			return client, errors.New("geth endpoint is not fully synced")
		}
	}

	return client, nil
}

func (client Client) Synced() (bool, error) {
	var res json.RawMessage
	if err := client.rpc.Call(&res, "eth_syncing"); err != nil {
		return false, errors.Wrap(err, "rpc")
	}

	if string(res) != "false" {
		return false, nil
	}

	return true, nil
}

func (client Client) LatestBlockNum() (*big.Int, error) {
	var hexStr string
	if err := client.rpc.Call(&hexStr, "eth_blockNumber"); err != nil {
		return nil, errors.Wrap(err, "rpc")
	}

	hexStr = utils.RemoveHexPrefix(hexStr)

	// pad if hex length is odd
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	hexBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, errors.Wrap(err, "hex decoding")
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
	if err := client.rpc.Call(&res, "eth_accounts"); err != nil {
		return nil, errors.Wrap(err, "rpc")
	}

	var addrs []string
	if err := json.Unmarshal(res, &addrs); err != nil {
		return nil, errors.Wrap(err, "json unmarshaling")
	}

	// convert to the correct type
	result := make([]common.Address, len(addrs))
	for i, addr := range addrs {
		result[i] = common.HexToAddress(addr)
	}

	return result, nil
}
