package query

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
)

// StatusCmd returns the current state of the eth connection (syncing, crashed etc)
func StatusCmd() *cobra.Command {
	return statuscmd
}

var statuscmd = &cobra.Command {
	Use:   "status",
	Short: "Check current block number",
	Long: `returns current block number of plasmad connection`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.NewCLIContext()
		height, err := client.Height(ctx)
		if err != nil {
			return fmt.Errorf("error retrieving current block height, {%s}", err)
		} else {
			fmt.Printf("current plasma block height = {%s}", height)
		}
		return nil
	},
}

