package cmd

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagPrivateKey = "privatekey"
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().BoolP(flagPrivateKey, "P", false, "read the the private key directly from the argument in hexadecimal format")
	viper.BindPFlags(importCmd.Flags())
}

var importCmd = &cobra.Command{
	Use:   "import <keyfile>",
	Short: "Import a private key into a new account on the sidechain",
	Long: `
plasmacli import <keyfile>
plasmacli import --key <private key>

Imports an unencrypted private key from <keyfile> and creates a new account on the sidechain.
Prints the address. If the key flag is set, the private key will be read directly from the argument
in hexadecimal format.

The keyfile is assumed to contain an unencrypted private key in hexadecimal format.
The keyfile must also be an absolute path

The account is saved in encrypted format, you are prompted for a passphrase.
You must remember this passphrase to unlock your account in the future.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		arg := args[0]

		var key *ecdsa.PrivateKey
		var err error
		if viper.GetBool(flagPrivateKey) {
			key, err = crypto.HexToECDSA(arg)
			if err != nil {
				return fmt.Errorf("failed parsing private key: %s", err)
			}
		} else {
			key, err = crypto.LoadECDSA(arg)
			if err != nil {
				return fmt.Errorf("failed loading the keyfile : %s", err)
			}
		}

		acct, err := keystore.ImportECDSA(key)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully imported. Address: 0x%x\n", acct.Address)
		return nil
	},
}
