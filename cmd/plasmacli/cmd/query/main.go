package query

import (
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query information related to the sidechain",
}

func QueryCmd() *cobra.Command {
	return queryCmd
}
