package keys

import (
	"errors"
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
)

func init() {
	keysCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long:  "Return a list of all account addresses stored by the local keystore",
	RunE: func(cmd *cobra.Command, args []string) error {

		iter, db := ks.AccountIterator()
		if iter == nil || db == nil {
			return errors.New("unexpected error encountered when opening account data")
		}
		defer db.Close()
		fmt.Printf("NAME:\t\tADDRESS:\n")
		for iter.Next() {
			var name string
			if err := rlp.DecodeBytes(iter.Key(), &name); err != nil {
				return err
			}
			fmt.Printf("%s\t\t%v\n", name, ethcmn.BytesToAddress(iter.Value()).Hex())
		}
		iter.Release()

		return nil
	},
}
