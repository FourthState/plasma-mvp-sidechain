package keys

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

// ImportCmd returns the keys import command
func ImportCmd() *cobra.Command {
	importCmd.Flags().String(fileF, "", "read the private key from raw private keyfile (must be absolute path)")
	importCmd.Flags().String(encryptF, "", "read the private key from an encrypted geth-compatible keyfile (must be absolute path)")
	return importCmd
}

var importCmd = &cobra.Command{
	Use:   "import <name> <privatekey>",
	Short: "Import a private key",
	Long: `Imports an unencrypted private key read in hexadecimal format and creates a new account on the sidechain.
Prints the address. 

Usage:
	plasmacli import <name> <privatekey>
	plasmacli import <name> --file <filepath>
	plasmacli import <name> --encrypted-file <filepath>

If the file flag is set:
The keyfile is assumed to contain an unencrypted private key in hexadecimal format.
The keyfile must also be an absolute path

If the encrypted-file flag is set:
The keyfile is assumed to contain an encrypted geth-compatible keyfile.
The keyfile must also be an absolute path

The account is saved in encrypted format, you are prompted for a passphrase.
You must remember this passphrase to unlock your account in the future.
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())

		name := args[0]
		var key *ecdsa.PrivateKey
		var err error
		file := viper.GetString(fileF)
		efile := viper.GetString(encryptF)
		if efile != "" {
			keybuf, err := ioutil.ReadFile(efile)
			if err != nil {
				return fmt.Errorf("failed loading the keyfile: %s", err)
			}

			_, err = store.Import(name, keybuf)
			if err != nil {
				return fmt.Errorf("error Importing keyfile: %s", err)
			}
			fmt.Println("Successfully imported.")
			return nil
		}

		if file != "" {
			key, err = crypto.LoadECDSA(file)
			if err != nil {
				return fmt.Errorf("failed loading the keyfile: { %s }", err)
			}
		} else {
			if len(args) < 2 {
				return errors.New("please provide an unencrypted private key if the --file flag is not set")
			}
			key, err = crypto.HexToECDSA(args[1])
			if err != nil {
				return fmt.Errorf("failed parsing private key: { %s }", err)
			}
		}

		address, err := store.ImportECDSA(name, key)
		if err != nil {
			return err
		}

		fmt.Println("Successfully imported.")
		fmt.Printf("NAME: %s\t\tADDRESS: 0x%x\n", name, address)
		return nil
	},
}
