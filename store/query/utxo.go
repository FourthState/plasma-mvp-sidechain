package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/big"
)

const (
	// QueryBalance retrieves the aggregate value of
	// the set of owned by the specified address
	QueryBalance = "balance"

	// QueryInfo retrieves the entire utxo set owned
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

func NewOutputQuerier(outputStore store.OutputStore) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		if len(path) == 0 {
			return nil, ErrInvalidPath("path not specified")
		}

		switch path[0] {
		case QueryBalance:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("expected balance/<address>")
			}
			addr := common.HexToAddress(path[1])
			total, err := queryBalance(ctx, outputStore, addr)
			if err != nil {
				return nil, sdk.ErrInternal(fmt.Sprintf("failed query balance for 0x%x", addr))
			}
			return []byte(total.String()), nil

		case QueryInfo:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("expected info/<address>")
			}
			addr := common.HexToAddress(path[1])
			utxos, err := queryInfo(ctx, outputStore, addr)
			if err != nil {
				return nil, err
			}
			data, e := json.Marshal(utxos)
			if e != nil {
				return nil, sdk.ErrInternal("serialization error")
			}
			return data, nil
		case QueryTxOutput:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("expected txo/<position>")
			}
			pos, e := plasma.FromPositionString(path[1])
			if e != nil {
				return nil, sdk.ErrInternal("position decoding error")
			}
			txo, err := queryTxOutput(ctx, outputStore, pos)
			if err != nil {
				return nil, err
			}
			data, e := json.Marshal(txo)
			if e != nil {
				return nil, sdk.ErrInternal("serialization error")
			}
			return data, nil
		case QueryTxInput:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("expected input/<position>")
			}
			pos, e := plasma.FromPositionString(path[1])
			if e != nil {
				return nil, sdk.ErrInternal("position decoding error")
			}
			txInput, err := queryTxInput(ctx, outputStore, pos)
			if err != nil {
				return nil, err
			}
			data, e := json.Marshal(txInput)
			if e != nil {
				return nil, sdk.ErrInternal("serialization error")
			}
			return data, nil
		case QueryTx:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("expected tx/<hash>")
			}
			tx, ok := outputStore.GetTx(ctx, []byte(path[1]))
			if !ok {
				return nil, ErrTxDNE(fmt.Sprintf("no transaction exists for the hash provided: %x", []byte(path[1])))
			}
			data, e := json.Marshal(tx)
			if e != nil {
				return nil, sdk.ErrInternal("serialization error")
			}
			return data, nil

		default:
			return nil, ErrInvalidPath("unregistered endpoint")
		}
	}
}

func queryBalance(ctx sdk.Context, outputStore store.OutputStore, addr common.Address) (*big.Int, sdk.Error) {
	acc, ok := outputStore.GetWallet(ctx, addr)
	if !ok {
		return nil, ErrWalletDNE(fmt.Sprintf("no wallet exists for the address provided: 0x%x", addr))
	}

	return acc.Balance, nil
}

func queryInfo(ctx sdk.Context, outputStore store.OutputStore, addr common.Address) ([]store.TxOutput, sdk.Error) {
	acc, ok := outputStore.GetWallet(ctx, addr)
	if !ok {
		return nil, ErrWalletDNE(fmt.Sprintf("no wallet exists for the address provided: 0x%x", addr))
	}
	return outputStore.GetUnspentForWallet(ctx, acc), nil
}

func queryTxOutput(ctx sdk.Context, outputStore store.OutputStore, pos plasma.Position) (store.TxOutput, sdk.Error) {
	output, ok := outputStore.GetOutput(ctx, pos)
	if !ok {
		return store.TxOutput{}, ErrOutputDNE(fmt.Sprintf("no output exists for the position provided: %s", pos))
	}

	tx, ok := outputStore.GetTxWithPosition(ctx, pos)
	if !ok {
		return store.TxOutput{}, ErrTxDNE(fmt.Sprintf("no transaction exists for the position provided: %s", pos))
	}

	txo := store.NewTxOutput(output.Output, pos, tx.ConfirmationHash, tx.Transaction.TxHash(), output.Spent, output.SpenderTx)

	return txo, nil
}

func queryTxInput(ctx sdk.Context, outputStore store.OutputStore, pos plasma.Position) (store.TxInput, sdk.Error) {
	output, ok := outputStore.GetOutput(ctx, pos)
	if !ok {
		return store.TxInput{}, ErrOutputDNE(fmt.Sprintf("no output exists for the position provided: %s", pos))
	}

	tx, ok := outputStore.GetTxWithPosition(ctx, pos)
	if !ok {
		return store.TxInput{}, ErrTxDNE(fmt.Sprintf("no transaction exists for the position provided: %s", pos))
	}

	inputPositions := tx.Transaction.InputPositions()
	var inputAddresses []common.Address
	for _, inPos := range inputPositions {
		input, ok := outputStore.GetOutput(ctx, inPos)
		if !ok {
			panic(fmt.Sprintf("Corrupted store: input position for given transaction does not exist: %s", pos))
		}
		inputAddresses = append(inputAddresses, input.Output.Owner)
	}

	return store.NewTxInput(output.Output, pos, tx.Transaction.TxHash(), inputAddresses, inputPositions), nil
}
