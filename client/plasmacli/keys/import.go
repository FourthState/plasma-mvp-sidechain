package keys

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	flagPrivateKey = "privatekey"
	flagFile       = "file"
)

func init() {
	keysCmd.AddCommand(importCmd)
	importCmd.Flags().StringP(flagPrivateKey, "P", "", "read the the private key directly from the argument in hexadecimal format")
	importCmd.Flags().String(flagFile, "", "read the private key from the specified keyfile")
	viper.BindPFlags(importCmd.Flags())
}

var importCmd = &cobra.Command{
	Use:   "import <name>",
	Short: "Import a private key",
	Long: `
plasmacli import <name> --file <keyfile>
plasmacli import <name> -P <private key>

Imports an unencrypted private key from <keyfile> and creates a new account on the sidechain.
Prints the address. If the privatekey flag is set, the private key will be read directly from the argument
in hexadecimal format.

The keyfile is assumed to contain an unencrypted private key in hexadecimal format.
The keyfile must also be an absolute path

The account is saved in encrypted format, you are prompted for a passphrase.
You must remember this passphrase to unlock your account in the future.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		var key *ecdsa.PrivateKey
		var err error
		privateKey := viper.GetString(flagPrivateKey)
		if privateKey != "" {
			key, err = crypto.HexToECDSA(privateKey)
			if err != nil {
				return fmt.Errorf("failed parsing private key: %s", err)
			}
		} else {
			key, err = crypto.LoadECDSA(viper.GetString(flagFile))
			if err != nil {
				return fmt.Errorf("failed loading the keyfile : %s", err)
			}
		}

		dir := accountDir()
		db, err := leveldb.OpenFile(dir, nil)
		if err != nil {
			return err
		}

		acct, err := keystore.ImportECDSA(key)
		if err != nil {
			return err
		}

		nameKey, err := rlp.EncodeToBytes(name)
		if err != nil {
			return err
		}

		if err = db.Put(nameKey, acct.Address.Bytes(), nil); err != nil {
			return err
		}

		fmt.Println("Successfully imported.")
		printAccount(name, acct.Address)
		return nil
	},
}
