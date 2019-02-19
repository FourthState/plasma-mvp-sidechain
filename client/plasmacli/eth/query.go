package eth

import (
	"github.com/spf13/cobra"
)

func init() {
	ethCmd.AddCommand(ethQueryCmd)
}

var ethQueryCmd = &cobra.Command{
	Use: "query",
}
