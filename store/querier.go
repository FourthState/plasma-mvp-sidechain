package store

import (
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

	// QueryBlocks retrieves full information about a
	// speficied block
	QueryBlock = "block"

	// QueryBlocs retrieves metadata about 10 blocks from
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

func queryBlock(ctx sdk.Context, ds DataStore, path []string) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, ErrInvalidPath(fmt.Sprintf("expected %s/<height>", QueryBlock))
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
		return nil, ErrInvalidPath(fmt.Sprintf("expected %s/<height> or %s/latest", QueryBlocks, QueryBlocks))
	}

	var height *big.Int
	if path[0] == "latest" {
		// latest 10 blocks
		height = ds.PlasmaBlockHeight(ctx)
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
		return nil, ErrInvalidPath(fmt.Sprintf("expected %s/<address>", QueryBalance))
	}

	addr, err := parseAddress(path[0])
	if err != nil {
		return nil, err
	}

	acc, ok := ds.GetWallet(ctx, addr)
	if !ok {
		return nil, ErrDNE(fmt.Sprintf("no wallet exists for the address provided: 0x%x", addr))
	}

	return []byte(acc.Balance.String()), nil
}

func queryInfo(ctx sdk.Context, ds DataStore, path []string) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, ErrInvalidPath(fmt.Sprintf("expected %s/<address>", QueryInfo))
	}

	addr, err := parseAddress(path[0])
	if err != nil {
		return nil, err
	}

	acc, ok := ds.GetWallet(ctx, addr)
	if !ok {
		return nil, ErrDNE(fmt.Sprintf("no wallet exists for the address provided: 0x%x", addr))
	}

	outputs := ds.GetUnspentForWallet(ctx, acc)
	return marshalResponse(outputs)
}

func queryTxOutput(ctx sdk.Context, ds DataStore, path []string) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, ErrInvalidPath(fmt.Sprintf("expected %s/<position>", QueryTxOutput))
	}

	pos, err := plasma.FromPositionString(path[0])
	if err != nil {
		return nil, ErrInvalidPath("position is encoded in the format (blocknum,txIndex,oIndex,depositNonce)")
	}

	o, ok := ds.GetOutput(ctx, pos)
	if !ok {
		return nil, ErrDNE(fmt.Sprintf("no output exists for the position provided: %s", pos))
	}

	tx, ok := ds.GetTxWithPosition(ctx, pos)
	if !ok {
		return nil, ErrDNE(fmt.Sprintf("no transaction exists for the position provided: %s", pos))
	}

	txo := NewTxOutput(o.Output, pos, tx.ConfirmationHash, tx.Transaction.TxHash(), o.Spent, o.SpenderTx)
	return marshalResponse(txo)
}

func queryTxInput(ctx sdk.Context, ds DataStore, path []string) ([]byte, sdk.Error) {
	if len(path) != 1 {
		return nil, ErrInvalidPath(fmt.Sprintf("expected %s/<position>", QueryTxInput))
	}

	pos, err := plasma.FromPositionString(path[0])
	if err != nil {
		return nil, ErrInvalidPath("position is encoded in the format (blocknum,txIndex,oIndex,depositNonce)")
	}

	tx, ok := ds.GetTxWithPosition(ctx, pos)
	if !ok {
		return nil, ErrDNE(fmt.Sprintf("no transaction exists for the position provided: %s", pos))
	}

	o, ok := ds.GetOutput(ctx, pos)
	if !ok {
		return nil, ErrDNE(fmt.Sprintf("no output exists for the position provided: %s", pos))
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
		return nil, ErrInvalidPath(fmt.Sprintf("expected %s/<txhash>", QueryTx))
	}
	txHash := path[0]
	if len(txHash) >= 2 && txHash[:2] == "0x" || txHash[:2] == "0X" {
		txHash = txHash[2:]
	}
	if _, ok := new(big.Int).SetString(txHash, 16); !ok || len(txHash) != 64 {
		return nil, ErrInvalidPath("txHash must be a 32-byte (64 character) hexadecimal string")
	}

	tx, ok := ds.GetTx(ctx, []byte(txHash))
	if !ok {
		return nil, ErrDNE(fmt.Sprintf("no transaction exists for the hash provided: %s", txHash))
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

func parseAddress(addrString string) (ethcmn.Address, sdk.Error) {
	if len(addrString) > 2 && addrString[:2] == "0x" {
		addrString = addrString[2:]
	}
	if len(addrString)%2 != 0 {
		addrString = "0" + addrString
	}

	if !ethcmn.IsHexAddress(addrString) {
		return utils.ZeroAddress, ErrInvalidPath(fmt.Sprintf("%s/<address> must be a valid ethereum address", QueryBalance))
	}

	return ethcmn.HexToAddress(addrString), nil
}
