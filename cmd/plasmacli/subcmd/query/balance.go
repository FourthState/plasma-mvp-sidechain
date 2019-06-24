package query

import (
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BalanceCmd() *cobra.Command {
	balanceCmd.Flags().StringP(addrF, "A", "", "query based on address")
	return balanceCmd
}

var balanceCmd = &cobra.Command{
	Use:   "balance <name>",
	Short: "Query plasma chain balance",
	Long: `Query for the total balance across utxos.
	
Usage: 
	plasmacli eth query balance <account>
	plasmacli eth query balance --address <address>`,
	SilenceUsage: true,
	Args:         cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		viper.BindPFlags(cmd.Flags())
		ctx := context.NewCLIContext()
		var addr ethcmn.Address

		if viper.GetString(addrF) != "" {
			addr = ethcmn.HexToAddress(viper.GetString(addrF))
		} else if len(args) > 0 {
			if addr, err = ks.GetAccount(args[0]); err != nil {
				return fmt.Errorf("failed to retrieve account: { %s }", err)
			}
		} else {
			return fmt.Errorf("please provide an account or use the address flag")
		}

		queryPath := fmt.Sprintf("custom/utxo/balance/%s", addr.Hex())
		total, err := ctx.Query(queryPath, nil)
		if err != nil {
			return err
		}

		fmt.Printf("Address: %0x\n", addr)
		fmt.Printf("Total: %s\n", string(total))
		return nil
	},
}
