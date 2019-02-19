package eth

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"strconv"
)

func init() {
	ethQueryCmd.AddCommand(depositsCmd)
	depositsCmd.Flags().Bool(allF, false, "all deposits will be displayed")
	depositsCmd.Flags().String(limitF, "1", "The amount of deposits to be displayed")
	viper.BindPFlags(depositsCmd.Flags())
}

var depositsCmd = &cobra.Command{
	Use:   "deposits <nonce>",
	Short: "Query for a deposit that occured on the rootchain",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var curr, end int64

		if viper.GetBool(allF) { // Print all deposits
			curr = 1
			lastNonce, err := rc.operatorSession.DepositNonce()
			if err != nil {
				return err
			}
			end = lastNonce.Int64()
		} else if len(args) > 0 { // Use command line arg as starting nonce
			if curr, err = strconv.ParseInt(args[0], 10, 64); err != nil {
				return fmt.Errorf("could not parse nonce: %v\n", err)
			}

			lim, err := strconv.ParseInt(viper.GetString(limitF), 10, 64)
			if err != nil {
				return err
			}

			end = curr + lim
		} else {
			return fmt.Errorf("please provide a starting nonce")
		}

		if err = displayDeposits(curr, end); err != nil {
			return err
		}

		return err
	},
}

func displayDeposits(curr int64, end int64) error {
	for curr < end {
		deposit, err := rc.operatorSession.Deposits(big.NewInt(curr))
		if err != nil {
			return err
		}
		fmt.Printf("Owner: 0x%x\nAmount: %d\nNonce: %d\n\n", deposit.Owner, deposit.Amount, curr)
		curr++
	}
	return nil
}
