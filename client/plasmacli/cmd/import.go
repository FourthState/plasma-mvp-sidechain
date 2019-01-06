package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import <keyfile>",
	Short: "Import a private key into a new account on the sidechain",
	Long: `
plasmacli import <keyfile>

Imports an unencrypted private key from <keyfile> and creates a new account on the sidechain.
Prints the address.

The keyfile is assumed to contain an unencrypted private key in hexadecimal format.

The account is saved in encrypted format, you are prompted for a passphrase.

You must remember this passphrase to unlock your account in the future.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyfile := args[0]
		if len(keyfile) == 0 {
			return fmt.Errorf("keyfile must be given as argument")
		}
		key, err := crypto.LoadECDSA(keyfile)
		if err != nil {
			return err
		}

		acct, err := keystore.ImportECDSA(key)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully imported. Address: {%x}\n", acct.Address)
		return nil
	},
}
