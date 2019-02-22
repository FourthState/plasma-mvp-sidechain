package eth

import (
	"fmt"

	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	queryCmd.AddCommand(balanceCmd)
	balanceCmd.Flags().StringP(addrF, "A", "", "query based on address")
	viper.BindPFlags(balanceCmd.Flags())
}

var balanceCmd = &cobra.Command{
	Use:   "balance <account>",
	Short: "Query for balance avaliable for withdraw from rootchain",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var addr ethcmn.Address

		if viper.GetString(addrF) != "" {
			addr = ethcmn.HexToAddress(viper.GetString(addrF))
		} else if len(args) > 0 {
			if addr, err = ks.Get(args[0]); err != nil {
				return fmt.Errorf("failed to retrieve account: { %s }", err)
			}
		} else {
			return fmt.Errorf("please provide an account or use the address flag")
		}

		balance, err := rc.contract.BalanceOf(nil, addr)
		if err != nil {
			return fmt.Errorf("failed to retrieve balance: { %s }", err)
		}

		fmt.Printf("Rootchain Balance: %d\n", balance)
		return nil
	},
}
