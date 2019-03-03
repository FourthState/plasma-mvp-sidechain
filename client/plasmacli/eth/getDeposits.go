package eth

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"strconv"
)

func init() {
	queryCmd.AddCommand(getDepositsCmd)
	getDepositsCmd.Flags().Bool(allF, false, "all deposits will be displayed")
	getDepositsCmd.Flags().String(limitF, "1", "amount of deposits to be displayed")
}

var getDepositsCmd = &cobra.Command{
	Use:   "deposit <nonce>",
	Short: "Query for a deposit that occured on the rootchain",
	Long: `Queries for deposits that occured on the rootchain.

Usage:
	plasmacli eth query deposit <nonce>
	plasmacli eth query deposit <nonce> --limit <number>
	plasmacli eth query deposit --all`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		viper.BindPFlags(cmd.Flags())

		var curr, lim int64

		if lim, err = strconv.ParseInt(viper.GetString(limitF), 10, 64); err != nil {
			return fmt.Errorf("failed to parse limit - %s", err)
		}

		if viper.GetBool(allF) { // Print all deposits
			curr = 1
			lastNonce, err := rc.contract.DepositNonce(nil)
			if err != nil {
				return fmt.Errorf("failed to trying to get last deposit nonce: { %s }", err)
			}
			lim = lastNonce.Int64() - 1
		} else if len(args) > 0 { // Use command line arg as starting nonce
			if curr, err = strconv.ParseInt(args[0], 10, 64); err != nil {
				return fmt.Errorf("failed to parse nonce - %s", err)
			}

			end = curr + lim
		} else {
			return fmt.Errorf("please provide a nonce")
		}

		if err = displayDeposits(curr, lim); err != nil {
			return fmt.Errorf("failed while displaying deposits - %s", err)
		}

		return err
	},
}

func displayDeposits(curr, lim int64) error {
	for lim > 0 {
		deposit, err := rc.contract.Deposits(nil, big.NewInt(curr))
		if err != nil {
			return err
		}

		if deposit.EthBlockNum.Int64() == 0 {
			break
		}

		fmt.Printf("Owner: 0x%x\nAmount: %d\nNonce: %d\nRootchain Block: %d\n\n", deposit.Owner, deposit.Amount, curr, deposit.EthBlockNum)
		curr++
		lim--
	}

	return nil
}
