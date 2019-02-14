package keys

import (
	"errors"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
)

func init() {
	keysCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Create a new account",
	Long:  `Add an encrypted account to your local keystore.`,
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

		if _, err = db.Get(key, nil); err == nil {
			return errors.New("You are trying to override an existing private key name. Please delete it first.")
		}

		address, err := keystore.NewAccount()
		if err != nil {
			return err
		}

		if err = db.Put(key, address.Bytes(), nil); err != nil {
			return err
		}

		fmt.Println("\n**Important** do not lose your passphrase.")
		fmt.Println("It is the only way to recover your account")
		fmt.Println("You should export this account and store it in a secure location")
		fmt.Printf("Your account data is stored in  %v\n", dir)
		fmt.Printf("NAME: %v\tADDRESS: %v\n", name, address.Hex())
		return nil
	},
}
