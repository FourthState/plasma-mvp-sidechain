package eth

import (
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"math/big"
	"strconv"
)

func init() {
	ethCmd.AddCommand(depositCmd)
}

var depositCmd = &cobra.Command{
	Use:   "deposit <amount> <account>",
	Short: "Deposit to rootchain contract",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := ks.GetKey(args[1])
		if err != nil {
			return err
		}

		amt, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("Could not parse specified amount: %v", err)
		}

		opts := rc.operatorSession.TransactOpts
		auth := bind.NewKeyedTransactor(key)
		defer func() {
			rc.operatorSession.TransactOpts = opts
		}()
		rc.operatorSession.TransactOpts = bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: 3141592, // aribitrary
			Value:    big.NewInt(amt),
		}

		if _, err := rc.operatorSession.Deposit(crypto.PubkeyToAddress(key.PublicKey)); err != nil {
			return err
		}

		fmt.Printf("Successfully deposited\n")
		return nil
	},
}
