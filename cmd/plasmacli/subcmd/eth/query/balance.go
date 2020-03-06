package query

import (
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// BalanceCmd returns the eth query balance command
func BalanceCmd() *cobra.Command {
	return balanceCmd
}

var balanceCmd = &cobra.Command{
	Use:          "balance <account/address>",
	Short:        "Query for balance avaliable for withdraw from rootchain",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			addr ethcmn.Address
			err  error
		)

		if !ethcmn.IsHexAddress(args[0]) {
			fmt.Println("Hex address not provided, retrieving account..")
			if addr, err = ks.GetAccount(args[0]); err != nil {
				return fmt.Errorf("failed account retrieval: %s", err)
			}
		} else {
			addr = ethcmn.HexToAddress(args[0])
		}

		cmd.SilenceUsage = true

		balance, err := plasmaContract.BalanceOf(nil, addr)
		if err != nil {
			return fmt.Errorf("failed to retrieve balance: { %s }", err)
		}

		fmt.Printf("Rootchain Balance: %d\n", balance)
		return nil
	},
}
