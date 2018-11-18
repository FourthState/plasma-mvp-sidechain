package cmd

import (
	"errors"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"

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
			return errors.New("keyfile must be given as argument")
		}
		key, err := crypto.LoadECDSA(keyfile)
		if err != nil {
			return err
		}

		buf := client.BufferStdin()
		passphrase, err := client.GetCheckPassword("Please set a passphrase for your imported account.\nPassphrase:", "Repeat the passphrase:", buf)
		if err != nil {
			return err
		}

		dir := viper.GetString(FlagHomeDir)
		ks := client.GetKeyStore(dir)

		acct, err := ks.ImportECDSA(key, passphrase)
		if err != nil {
			return err
		}

		fmt.Printf("Address: {%x}\n", acct.Address)
		return nil
	},
}
