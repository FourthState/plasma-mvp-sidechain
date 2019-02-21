package eth

import (
	"fmt"
	"math/big"
	"strconv"

	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

func init() {
	ethCmd.AddCommand(depositCmd)
}

var depositCmd = &cobra.Command{
	Use:   "deposit <amount> <account>",
	Short: "Deposit to rootchain contract",
	Long:  `Deposit to the rootchain contract as specified in plasma.toml.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := ks.GetKey(args[1])
		if err != nil {
			return err
		}

		amt, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse amount - %v", err)
		}

		auth := bind.NewKeyedTransactor(key)
		defer func() {
			rc.session.TransactOpts = bind.TransactOpts{}
		}()
		rc.session.TransactOpts = bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: 3141592, // aribitrary
			Value:    big.NewInt(amt),
		}

		if _, err := rc.session.Deposit(crypto.PubkeyToAddress(key.PublicKey)); err != nil {
			return fmt.Errorf("failed to deposit - %v", err)
		}

		fmt.Printf("Successfully deposited\n")
		return nil
	},
}
