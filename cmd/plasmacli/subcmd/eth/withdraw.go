package eth

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

// WithdrawCmd returns the eth withdraw command
func WithdrawCmd() *cobra.Command {
	withdrawCmd.Flags().StringP(gasLimitF, "g", "150000", "gas limit for ethereum transaction")
	return withdrawCmd
}

var withdrawCmd = &cobra.Command{
	Use:   "withdraw <account>",
	Short: "Withdraw all available funds from rootchain contract",
	Long: `Withdraw all available funds from the rootchain contract

Usage:
	plasmacli eth withdraw <account> --gas-limit 30000`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())

		key, err := store.GetKey(args[0])
		if err != nil {
			return fmt.Errorf("failed to retrieve account: %s", err)
		}

		gasLimit, err := strconv.ParseUint(viper.GetString(gasLimitF), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse gas limit: %s", err)
		}

		cmd.SilenceUsage = true

		// bind key, generate transact opts
		auth := bind.NewKeyedTransactor(key)
		transactOpts := &bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: gasLimit,
		}

		tx, err := plasmaContract.Withdraw(transactOpts)
		if err != nil {
			return fmt.Errorf("failed to withdraw: %s", err)
		}

		fmt.Printf("Successfully sent withdraw transaction\nTransaction Hash: 0x%x\n", tx.Hash())
		return nil
	},
}
