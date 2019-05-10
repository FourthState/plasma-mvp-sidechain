package query

import (
	"encoding/json"
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func init() {
	queryCmd.AddCommand(balanceCmd)
}

var balanceCmd = &cobra.Command{
	Use:   "balance <name>",
	Short: "Query account balance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext().WithCodec(codec.New()).WithTrustNode(true)
		name := args[0]

		addr, err := ks.GetAccount(name)
		if err != nil {
			return err
		}

		queryRoute := fmt.Sprintf("custom/utxo/balance/%s", addr.Hex())
		data, err := ctx.Query(queryRoute, nil)
		if err != nil {
			return err
		}

		var total string
		if err := json.Unmarshal(data, &total); err != nil {
			return err
		}

		fmt.Printf("Address: %0x\n", addr)
		fmt.Printf("Total: %s\n", total)
		return nil
	},
}
