package eth

import (
	"github.com/spf13/cobra"
)

func init() {
	ethQueryCmd.AddCommand(rootchainCmd)
}

var rootchainCmd = &cobra.Command{
	Use: "rootchain",
	RunE: func(cmd *cobra.Command, agrs []string) error {
		return nil
	},
}
