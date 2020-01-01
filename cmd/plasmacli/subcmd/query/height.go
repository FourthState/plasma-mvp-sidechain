package query

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
)

// HeightCmd returns the current block height of the plasmad connection
func HeightCmd() *cobra.Command {
	return heightcmd
}

var heightcmd = &cobra.Command {
	Use:   "height",
	Short: "Check current block height",
	Long: "Returns current block height of plasmad connection",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.NewCLIContext()
		height, err := client.Height(ctx)
		if err != nil {
			return fmt.Errorf("error retrieving current block height: %s", err)
		}
		fmt.Printf("current plasma block height: %s", height)
		return nil
	},
}

