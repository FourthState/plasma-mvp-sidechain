package keys

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
)

func init() {
	keysCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete the given address",
	Long:  `Deletes the account from the keystore if the passphrase is correct.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// retrieve account
		dir := AccountDir()
		db, err := leveldb.OpenFile(dir, nil)
		if err != nil {
			return err
		}
		defer db.Close()

		key, err := rlp.EncodeToBytes(name)
		if err != nil {
			return err
		}

		addr, err := db.Get(key, nil)
		if err != nil {
			return err
		}

		// delete from the keystore
		if err := keystore.Delete(ethcmn.BytesToAddress(addr)); err != nil {
			return err
		}

		if err := db.Delete(key, nil); err != nil {
			return err
		}

		fmt.Println("Account deleted.")
		return nil
	},
}
