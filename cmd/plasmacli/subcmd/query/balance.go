package query

import (
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// BalanceCmd returns the query balance command
func BalanceCmd() *cobra.Command {
	return balanceCmd
}

var balanceCmd = &cobra.Command{
	Use:          "balance <account/address>",
	Short:        "Total plasma chain balance across utxos",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()
		var (
			addr ethcmn.Address
			err  error
		)

		if !ethcmn.IsHexAddress(args[0]) {
			if addr, err = ks.GetAccount(args[0]); err != nil {
				return fmt.Errorf("failed local account retrieval: %s", err)
			}
		} else {
			addr = ethcmn.HexToAddress(args[0])
		}

		queryPath := fmt.Sprintf("custom/data/balance/%s", addr.Hex())
		total, err := ctx.Query(queryPath, nil)
		if err != nil {
			return err
		}

		fmt.Printf("Address: %0x\n", addr)
		fmt.Printf("Total: %s\n", string(total))
		return nil
	},
}
