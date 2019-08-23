package store

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/big"
)

const (
	// RouteName to mount the querier
	QuerierRouteName = DataStoreName

	/*** block routes ***/

	// QueryHeight retrieves the current block height
	QueryHeight = "height"

	// QueryBlocks retrieves full information about a
	// speficied block
	QueryBlock = "block"

	// QueryBlocks retrieves metadata about 10 blocks from
	// a specified start point or the last 10 from the latest
	// block
	QueryBlocks = "blocks"

	/*** output routes ***/

	// QueryBalance retrieves the aggregate value of
	// the set of owned by the specified address
	QueryBalance = "balance"

	// QueryInfo retrieves the entire output set owned
	// by the specified address
	QueryInfo = "info"

	// QueryTxOutput retrieves a single output at
	// the given position and returns it with transactional
	// information
	QueryTxOutput = "output"

	// QueryTxInput retrieves basic transaction data at
	// given position along with input information
	QueryTxInput = "input"

	// QueryTx retrieves a transaction at the given hash
	QueryTx = "tx"
)

// NewQuerier returns an SDK querier to interact with the store
func NewQuerier(ds DataStore) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		if len(path) == 0 {
			return nil, ErrInvalidPath("path not specified")
		}

		switch path[0] {
		case QueryHeight:
			return queryHeight(ctx, ds)
		case QueryBlock:
			return queryBlock(ctx, ds, path[1:])
		case QueryBlocks:
			return queryBlocks(ctx, ds, path[1:])
		case QueryBalance:
			return queryBalance(ctx, ds, path[1:])
		case QueryInfo:
			return queryInfo(ctx, ds, path[1:])
		case QueryTxOutput:
			return queryTxOutput(ctx, ds, path[1:])
		case QueryTxInput:
			return queryTxInput(ctx, ds, path[1:])
		case QueryTx:
			return queryTx(ctx, ds, path[1:])
		default:
			return nil, ErrInvalidPath("unregistered query path")
		}
	}
}

func queryHeight(ctx sdk.Context, ds DataStore) ([]byte, sdk.Error) {
	height := ds.PlasmaBlockHeight(ctx)
	if height == nil {
		height = utils.Big0
	}

	return []byte(height.String()), nil
}

func queryBlock(ctx sdk.Context, ds DataStore, path []string) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, ErrInvalidPath("expected %s/<height>", QueryBlock)
	}

	height, err := parseHeight(path[0])
	if err != nil {
		return nil, err
	}

	block, ok := ds.GetBlock(ctx, height)
	if !ok {
		return nil, ErrDNE("plasma block %s does not exist", height)
	}

	return marshalResponse(block)
}

func queryBlocks(ctx sdk.Context, ds DataStore, path []string) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, ErrInvalidPath("expected %s/<height> or %s/latest", QueryBlocks, QueryBlocks)
	}

	var height *big.Int
	if path[0] == "latest" {
		// latest 10 blocks
		height = ds.PlasmaBlockHeight(ctx)
		if height == nil {
			return nil, ErrDNE("no blocks")
		}

		bigNine := big.NewInt(9)
		if height.Cmp(bigNine) <= 0 {
			height = big.NewInt(1)
		} else {
			height = height.Sub(height, bigNine)
		}
	} else {
		// predefined starting point
		h, err := parseHeight(path[0])
		if err != nil {
			return nil, err
		}
		height = h
	}

	var blocks []Block
	for i := 0; i < 10; i++ {
		block, ok := ds.GetBlock(ctx, height)
		if !ok {
			break
		}
		blocks = append(blocks, block)
		height = height.Add(height, utils.Big1)
	}

	if len(blocks) != 0 {
		return nil, ErrDNE("no blocks")
	}

	return marshalResponse(blocks)
}

func queryBalance(ctx sdk.Context, ds DataStore, path []string) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, ErrInvalidPath("expected %s/<address>", QueryBalance)
	}

	addr, err := parseAddress(path[0])
	if err != nil {
		return nil, err
	}

	acc, ok := ds.GetWallet(ctx, addr)
	if !ok {
		return nil, ErrDNE("no wallet exists for the address provided: 0x%x", addr)
	}

	return []byte(acc.Balance.String()), nil
}

func queryInfo(ctx sdk.Context, ds DataStore, path []string) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, ErrInvalidPath("expected %s/<address>", QueryInfo)
	}

	addr, err := parseAddress(path[0])
	if err != nil {
		return nil, err
	}

	acc, ok := ds.GetWallet(ctx, addr)
	if !ok {
		return nil, ErrDNE("no wallet exists for the address provided: 0x%x", addr)
	}

	outputs := ds.GetUnspentForWallet(ctx, acc)
	return marshalResponse(outputs)
}

func queryTxOutput(ctx sdk.Context, ds DataStore, path []string) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, ErrInvalidPath("expected %s/<position>", QueryTxOutput)
	}

	pos, err := plasma.FromPositionString(path[0])
	if err != nil {
		return nil, ErrInvalidPath("position is encoded in the format (blocknum,txIndex,oIndex,depositNonce)")
	}

	o, ok := ds.GetOutput(ctx, pos)
	if !ok {
		return nil, ErrDNE("no output exists for the position provided: %s", pos)
	}

	tx, ok := ds.GetTxWithPosition(ctx, pos)
	if !ok {
		return nil, ErrDNE("no transaction exists for the position provided: %s", pos)
	}

	txo := NewTxOutput(o.Output, pos, tx.ConfirmationHash, tx.Transaction.TxHash(), o.Spent, o.SpenderTx)
	return marshalResponse(txo)
}

func queryTxInput(ctx sdk.Context, ds DataStore, path []string) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, ErrInvalidPath("expected %s/<position>", QueryTxInput)
	}

	pos, err := plasma.FromPositionString(path[0])
	if err != nil {
		return nil, ErrInvalidPath("position is encoded in the format (blocknum,txIndex,oIndex,depositNonce)")
	}

	tx, ok := ds.GetTxWithPosition(ctx, pos)
	if !ok {
		return nil, ErrDNE("no transaction exists for the position provided: %s", pos)
	}

	o, ok := ds.GetOutput(ctx, pos)
	if !ok {
		return nil, ErrDNE("no output exists for the position provided: %s", pos)
	}

	inputPositions := tx.Transaction.InputPositions()
	var inputAddresses []ethcmn.Address
	for _, inPos := range inputPositions {
		input, ok := ds.GetOutput(ctx, inPos)
		if !ok {
			panic(fmt.Sprintf("Corrupted store: input position for given transaction does not exist: %s", pos))
		}
		inputAddresses = append(inputAddresses, input.Output.Owner)
	}

	txinput := NewTxInput(o.Output, pos, tx.Transaction.TxHash(), inputAddresses, inputPositions)
	return marshalResponse(txinput)
}

func queryTx(ctx sdk.Context, ds DataStore, path []string) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, ErrInvalidPath("expected %s/<txhash>", QueryTx)
	}
	txHash, err := hex.DecodeString(utils.RemoveHexPrefix(path[0]))
	if err != nil {
		return nil, ErrInvalidPath("tx hash expected in hexadecimal format. hex: %s", err)
	} else if len(txHash) != 32 {
		return nil, ErrInvalidPath("tx hash expected to be 32 bytes in length")
	}

	tx, ok := ds.GetTx(ctx, txHash)
	if !ok {
		return nil, ErrDNE("no transaction exists for the hash provided: %s", txHash)
	}

	return marshalResponse(tx)
}

/** helpers **/

func marshalResponse(resp interface{}) ([]byte, sdk.Error) {
	data, err := json.Marshal(resp)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("json: %s", err)).
			WithDefaultCodespace(DefaultCodespace)
	}
	return data, nil
}

func parseHeight(height string) (*big.Int, sdk.Error) {
	h, ok := new(big.Int).SetString(height, 10)
	if !ok || h.Sign() <= 0 {
		return nil, ErrInvalidPath("block height must start from 1 in decimal format")
	}

	return h, nil
}

func parseAddress(addr string) (ethcmn.Address, sdk.Error) {
	addr = utils.RemoveHexPrefix(addr)

	if !ethcmn.IsHexAddress(addr) {
		return utils.ZeroAddress, ErrInvalidPath("address expected to be a valid 20-byte ethereum address")
	}

	return ethcmn.HexToAddress(addr), nil
}
