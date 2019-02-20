package eth

import (
	"fmt"

	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/spf13/cobra"
)

func init() {
	queryCmd.AddCommand(balanceCmd)
}

var balanceCmd = &cobra.Command{
	Use:   "balance <account>",
	Short: "Query for balance avaliable for withdraw from rootchain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, err := ks.Get(args[0])
		if err != nil {
			return err
		}

		balance, err := rc.session.BalanceOf(addr)
		if err != nil {
			return err
		}

		fmt.Printf("Rootchain Balance: %d\n", balance)
		return nil
	},
}
