package eth

import (
	"github.com/spf13/cobra"
)

func init() {
	ethQueryCmd.AddCommand(exitsCmd)
}

var exitsCmd = &cobra.Command{
	Use: "exits",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
