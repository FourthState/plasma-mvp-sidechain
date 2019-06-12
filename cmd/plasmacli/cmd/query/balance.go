package query

import (
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
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

		// no more additional argument validation
		cmd.SilenceUsage = true

		addr, err := ks.GetAccount(name)
		if err != nil {
			return err
		}

		queryPath := fmt.Sprintf("custom/utxo/balance/%s", addr.Hex())
		total, err := ctx.Query(queryPath, nil)
		if err != nil {
			return err
		}

		fmt.Printf("Address: %0x\n", addr)
		fmt.Printf("Total: %s\n", string(total))
		return nil
	},
}
