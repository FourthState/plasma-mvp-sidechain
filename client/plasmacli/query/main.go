package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query information related to the sidechain",
}

func QueryCmd() *cobra.Command {
	return queryCmd
}

func OutputInfo(ctx context.CLIContext, pos plasma.Position) (store.OutputInfo, error) {
	// query for an output for the given position
	queryRoute := fmt.Sprintf("custom/utxo/output/%s", pos)
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return store.OutputInfo{}, err
	}

	var output store.OutputInfo
	if err := json.Unmarshal(data, &output); err != nil {
		return store.OutputInfo{}, err
	}

	return output, nil
}

func Tx(ctx context.CLIContext, hash []byte) (store.Transaction, error) {
	// query for a transaction using the provided hash
	queryRoute := fmt.Sprintf("custom/utxo/tx/%s", hash)
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
