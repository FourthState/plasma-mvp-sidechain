package query

import (
	"github.com/spf13/cobra"
)

func QueryCmd() *cobra.Command {
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
