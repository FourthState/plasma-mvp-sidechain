package eth

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/store"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"strconv"
)

func init() {
	ethCmd.AddCommand(depositCmd)
	depositCmd.Flags().StringP(gasLimitF, "g", "30000", "gas limit for ethereum transaction")
}

var depositCmd = &cobra.Command{
	Use:   "deposit <amount> <account>",
	Short: "Deposit to rootchain contract",
	Long: `Deposit to the rootchain contract as specified in plasma.toml.

Usage:
	plasmacli eth deposit <amount> <account> --gas-limit 30000`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())

		key, err := store.GetKey(args[1])
		if err != nil {
			return fmt.Errorf("failed to retrieve account: { %s }", err)
		}

		amt, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse amount: { %s }", err)
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
			Value:    big.NewInt(amt),
		}

		tx, err := rc.contract.Deposit(transactOpts, crypto.PubkeyToAddress(key.PublicKey))
		if err != nil {
			return fmt.Errorf("failed to deposit: { %s }", err)
		}

		fmt.Printf("Successfully sent deposit transaction\nTransaction Hash: 0x%x\n", tx.Hash())
		return nil
	},
}
