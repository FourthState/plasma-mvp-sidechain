package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	ethcmn "github.com/ethereum/go-ethereum/common"
)

// TxOutput executes a query to the data store for a transaction output
func TxOutput(ctx context.CLIContext, pos plasma.Position) (store.TxOutput, error) {
	// query for an output for the given position
	queryRoute := fmt.Sprintf("custom/%s/output/%s", store.QuerierRouteName, pos)
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return store.TxOutput{}, err
	}

	var output store.TxOutput
	if err := json.Unmarshal(data, &output); err != nil {
		return store.TxOutput{}, err
	}

	return output, nil
}

// TxInput executes a query to the data store for a transaction input
func TxInput(ctx context.CLIContext, pos plasma.Position) (store.TxInput, error) {
	// query for input info on the given position
	queryRoute := fmt.Sprintf("custom/%s/input/%s", store.QuerierRouteName, pos)
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return store.TxInput{}, err
	}

	var input store.TxInput
	if err := json.Unmarshal(data, &input); err != nil {
		return store.TxInput{}, err
	}

	return input, nil
}

// Tx executes a query to the data store for a transaction with the provided hash.
func Tx(ctx context.CLIContext, hash []byte) (store.Transaction, error) {
	// query for a transaction using the provided hash
	queryRoute := fmt.Sprintf("custom/%s/tx/%s", store.QuerierRouteName, hash)
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return store.Transaction{}, err
	}

	var tx store.Transaction
	if err := json.Unmarshal(data, &tx); err != nil {
		return store.Transaction{}, err
	}

	return tx, nil
}

// Info executes a query to the data store for the information on an address
func Info(ctx context.CLIContext, addr ethcmn.Address) ([]store.TxOutput, error) {
	// query for all utxos owned by this address
	queryRoute := fmt.Sprintf("custom/%s/info/%s", store.QuerierRouteName, addr.Hex())
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return nil, err
	}

	var utxos []store.TxOutput
	if err := json.Unmarshal(data, &utxos); err != nil {
		return nil, err
	}

	return utxos, nil
}
