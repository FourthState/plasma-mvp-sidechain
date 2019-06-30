package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
)

func QueryCmd() *cobra.Command {
	config.AddPersistentTMFlags(queryCmd)
	queryCmd.AddCommand(
		BalanceCmd(),
		BlockCmd(),
		BlocksCmd(),
		InfoCmd(),
	)

	return queryCmd
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query information related to the sidechain",
}

func TxOutput(ctx context.CLIContext, pos plasma.Position) (store.TxOutput, error) {
	// query for an output for the given position
	queryRoute := fmt.Sprintf("custom/utxo/output/%s", pos)
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
