package keys

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	flagName = "name"
)

func init() {
	keysCmd.AddCommand(updateCmd)
	updateCmd.Flags().String(flagName, "", "updated key name.")
	viper.BindPFlags(updateCmd.Flags())
}

var updateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "Update account passphrase or key name",
	Long: `Update local encrypted private keys to be encrypted with the new passphrase.
--name can be used to update the account name`,
	Args: cobra.ExactArgs(1),
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

		newName := viper.GetString(flagName)

		if newName != "" {
			if err = db.Delete(key, nil); err != nil {
				return err
			}

			key, err := rlp.EncodeToBytes(newName)
			if err != nil {
				return err
			}

			if err = db.Put(key, addr, nil); err != nil {
				return err
			}
			fmt.Println("Account name has been updated. ")
		} else {
			if err := keystore.Update(ethcmn.BytesToAddress(addr)); err != nil {
				return err
			}
			fmt.Println("Account passphrase has been updated.")
		}

		return nil
	},
}
