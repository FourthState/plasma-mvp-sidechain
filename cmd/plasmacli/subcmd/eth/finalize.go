package eth

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

// FinalizeCmd returns the eth finalize command
func FinalizeCmd() *cobra.Command {
	finalizeCmd.Flags().BoolP(depositsF, "D", false, "indicate that deposit exits should be finalized")
	finalizeCmd.Flags().StringP(gasLimitF, "g", "240000", "gas limit for ethereum transaction")
	return finalizeCmd
}

var finalizeCmd = &cobra.Command{
	Use:   "finalize <account>",
	Short: "Finalize exit queue on rootchain",
	Long: `Defaults to finalizing transaction exits. Use deposit flag to finalize deposit exit queue

Usage:
	plasmacli eth finalize <account> --gas-limit 30000
	plasmacli eth finalize <account> --deposits`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		viper.BindPFlags(cmd.Flags())

		key, err := store.GetKey(args[0])
		if err != nil {
			return fmt.Errorf("failed to retrieve account: { %s }", err)
		}

		gasLimit, err := strconv.ParseUint(viper.GetString(gasLimitF), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse gas limit: { %s }", err)
		}

		// bind key, generate transact opts
		auth := bind.NewKeyedTransactor(key)
		transactOpts := &bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: gasLimit,
		}

		var tx *eth.Transaction
		if viper.GetBool(depositsF) {
			tx, err = plasmaContract.FinalizeDepositExits(transactOpts)
		} else {
			tx, err = plasmaContract.FinalizeTransactionExits(transactOpts)
		}
		if err != nil {
			return fmt.Errorf("failed to finalize exits: { %s }", err)
		}

		fmt.Printf("Successfully sent finalize exits transaction\nTransaction Hash: 0x%x\n", tx.Hash())
		return nil
	},
}
