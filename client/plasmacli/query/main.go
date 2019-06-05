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

func Output(ctx context.CLIContext, pos plasma.Position) (store.Output, error) {
	// query for an output for the given position
	queryRoute := fmt.Sprintf("custom/utxo/output/%s", pos)
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return store.Output{}, err
	}

	var output store.Output
	if err := json.Unmarshal(data, &output); err != nil {
		return store.Output{}, err
	}

	return output, nil
}
