package query

import (
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	"github.com/spf13/cobra"
)

// RootCmd returns the query command for plasmacli
func RootCmd() *cobra.Command {
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
