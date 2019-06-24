package tx

import (
	"github.com/spf13/cobra"
)

const (
	accountF      = "accounts"
	addressF      = "address"
	asyncF        = "async"
	confirmSigs0F = "Input0ConfirmSigs"
	confirmSigs1F = "Input1ConfirmSigs"
	feeF          = "fee"
	inputsF       = "inputValues"
	ownerF        = "owner"
	positionF     = "position"
	replayF       = "replay"
)

func TxCmd() *cobra.Command {
	txCmd.AddCommand(
		IncludeCmd(),
		SpendCmd(),
		SignCmd(),
	)
	return txCmd
}

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Submit or interact with plasma chain txs",
}
