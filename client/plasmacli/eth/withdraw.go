package eth

import (
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

func init() {
	ethCmd.AddCommand(withdrawCmd)
	withdrawCmd.Flags().StringP(gasLimitF, "g", "21000", "gas limit for ethereum transaction")
	viper.BindPFlags(withdrawCmd.Flags())
}

var withdrawCmd = &cobra.Command{
	Use:   "withdraw <account>",
	Short: "Withdraw all avaliable funds from rootchain contract",
	Long: `Withdraw all avaliable funds from the rootchain contract

Usage:
	plasmacli eth withdraw <account> --gas-limit 30000`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := ks.GetKey(args[0])
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

		tx, err := rc.contract.Withdraw(transactOpts)
		if err != nil {
			return fmt.Errorf("failed to withdraw: {%s }", err)
		}

		fmt.Printf("Successfully sent withdraw transaction\nTransaction Hash: 0x%x", tx.Hash())
		return nil
	},
}
