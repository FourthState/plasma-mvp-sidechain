package client

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

// TxOutput retrieves the output located at `pos` and contextual transaction information
func TxOutput(ctx context.CLIContext, pos plasma.Position) (store.TxOutput, error) {
	queryRoute := fmt.Sprintf("custom/%s/%s/%s",
		store.QuerierRouteName, store.QueryTxOutput, pos)
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return store.TxOutput{}, err
	}

	var output store.TxOutput
	if err := json.Unmarshal(data, &output); err != nil {
		return store.TxOutput{}, fmt.Errorf("json: %s", err)
	}

	return output, nil
}

// TxInput retrieves the tx hash and the inputs that created the output located at `pos`
func TxInput(ctx context.CLIContext, pos plasma.Position) (store.TxInput, error) {
	queryRoute := fmt.Sprintf("custom/%s/%s/%s",
		store.QuerierRouteName, store.QueryTxInput, pos)
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return store.TxInput{}, err
	}

	var input store.TxInput
	if err := json.Unmarshal(data, &input); err != nil {
		return store.TxInput{}, fmt.Errorf("json: %s", err)
	}

	return input, nil
}

// Tx locates a transaction and given it's hash
// @param hash 32-byte hexadecimal string
func Tx(ctx context.CLIContext, hash []byte) (store.Transaction, error) {
	queryRoute := fmt.Sprintf("custom/%s/%s/%x",
		store.QuerierRouteName, store.QueryTx, hash)
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return store.Transaction{}, err
	}

	var tx store.Transaction
	if err := json.Unmarshal(data, &tx); err != nil {
		return store.Transaction{}, fmt.Errorf("json: %s", err)
	}

	return tx, nil
}

// Info retrieves the unspent utxo set of an owned address
func Info(ctx context.CLIContext, addr ethcmn.Address) ([]store.TxOutput, error) {
	queryRoute := fmt.Sprintf("custom/%s/%s/%s",
		store.QuerierRouteName, store.QueryInfo, addr.Hex())
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return nil, err
	}

	var utxos []store.TxOutput
	if err := json.Unmarshal(data, &utxos); err != nil {
		return nil, fmt.Errorf("json: %s", err)
	}

	return utxos, nil
}

// Balance retrieves the aggregate value across unspent utxos of an address
func Balance(ctx context.CLIContext, addr ethcmn.Address) (string, error) {
	queryRoute := fmt.Sprintf("custom/%s/%s/%s",
		store.QuerierRouteName, store.QueryBalance, addr.Hex())
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return "", err
	}

	// balance returned in string format
	return string(data), err
}

// Height retrieves the current plasma block height
func Height(ctx context.CLIContext) (string, error) {
	queryRoute := fmt.Sprintf("custom/%s/%s",
		store.QuerierRouteName, store.QueryHeight)
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return "", err
	}

	// data returned in stirng format
	return string(data), nil
}

// Block retrieves block information specified at `height`
func Block(ctx context.CLIContext, height *big.Int) (store.Block, error) {
	if height == nil || height.Sign() <= 0 {
		return store.Block{}, fmt.Errorf("block numbering starts at 1")
	}

	queryPath := fmt.Sprintf("custom/%s/%s/%s",
		store.QuerierRouteName, store.QueryBlock, height)
	data, err := ctx.Query(queryPath, nil)
	if err != nil {
		return store.Block{}, err
	}

	var block store.Block
	if err := json.Unmarshal(data, &block); err != nil {
		return block, fmt.Errorf("json: %s", err)
	}

	return block, nil
}

// Blocks retrieves 10 blocks from `startingHeight`. if `startingHeight == nil`, the latest 10 are retrieved
func Blocks(ctx context.CLIContext, startingHeight *big.Int) ([]store.Block, error) {
	if startingHeight != nil && startingHeight.Sign() <= 0 {
		return nil, fmt.Errorf("block height starts at 1")
	}

	var queryPath string
	if startingHeight == nil {
		queryPath = "latest"
	} else {
		queryPath = startingHeight.String()
	}

	// add prefix
	queryPath = fmt.Sprintf("custom/%s/%s/%s",
		store.QuerierRouteName, store.QueryBlocks, queryPath)
	data, err := ctx.Query(queryPath, nil)
	if err != nil {
		return nil, err
	}

	var blocks []store.Block
	if err = json.Unmarshal(data, &blocks); err != nil {
		return nil, fmt.Errorf("json: %s", err)
	}

	return blocks, nil
}

// Returns confirm sig results for given position
// Trusts connected full node
func ConfirmSignatures(ctx context.CLIContext, position plasma.Position) ([]byte, error) {
	key := store.GetOutputKey(position)
	hash, err := ctx.QueryStore(key, store.DataStoreName)
	if err != nil {
		return nil, err
	}

	txKey := store.GetTxKey(hash)
	txBytes, err := ctx.QueryStore(txKey, store.DataStoreName)

	var tx store.Transaction
	if err := rlp.DecodeBytes(txBytes, &tx); err != nil {
		return nil, fmt.Errorf("transaction decoding failed: %s", err.Error())
	}

	// Look for confirmation signatures
	// Ignore error if no confirm sig currently exists in store
	var sigs []byte
	if len(tx.SpenderTxs[position.OutputIndex]) > 0 {
		queryPath := fmt.Sprintf("custom/%s/%s/%s",
			store.QuerierRouteName, store.QueryTx, tx.SpenderTxs[position.OutputIndex])
		data, err := ctx.Query(queryPath, nil)
		if err != nil {
			return nil, err
		}

		var spenderTx store.Transaction
		if err := json.Unmarshal(data, &spenderTx); err != nil {
			return nil, fmt.Errorf("unmarshaling json query response: %s", err)
		}

		for _, input := range spenderTx.Transaction.Inputs {
			if input.Position.String() == position.String() {
				for _, sig := range input.ConfirmSignatures {
					sigs = append(sigs, sig[:]...)
				}
			}
		}
	}

	return sigs, nil
}

