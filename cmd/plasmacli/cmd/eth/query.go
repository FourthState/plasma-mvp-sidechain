package eth

import (
	"github.com/spf13/cobra"
)

func init() {
	ethCmd.AddCommand(queryCmd)
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query for rootchain related information",
}
