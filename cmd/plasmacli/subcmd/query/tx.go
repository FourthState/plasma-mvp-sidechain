package query

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/flags"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

// TxCommand -
func TxCommand() *cobra.Command {
	return txCmd
}

var txCmd *cobra.Command = &cobra.Command{
	Use:   "tx <txHash/position>",
	Short: "Query for information about a single transaction",
	Long: `Query for information about a single transaction.
If --verbose is set, additional information about the transaction will also be displayed`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *Cobra.Command, args []string) error {
		ctx := context.NewCLIContext()

		// argument validation
		arg := utils.RemoveHexPrefix(strings.TrimSpace(args[0]))
		if pos, err := plasma.FromPositionString(arg); err != nil {
			// position
		} else {
			// transaction hash
		}
	},
}
