package keys

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
)

func init() {
	keysCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long:  "Return a list of all account addresses stored by the local keystore",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := accountDir()
		db, err := leveldb.OpenFile(dir, nil)
		if err != nil {
			return err
		}

		iter := db.NewIterator(nil, nil)
		fmt.Printf("NAME:\t\tADDRESS:\n")
		for iter.Next() {
			var name string
			if err := rlp.DecodeBytes(iter.Key(), &name); err != nil {
				return err
			}
			printAccount(name, ethcmn.BytesToAddress(iter.Value()))
		}

		return nil
	},
}
