package keys

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/spf13/cobra"
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

		address, err := store.AddAccount(name)
		if err != nil {
			return err
		}

		fmt.Println("\n**Important** do not lose your passphrase.")
		fmt.Println("It is the only way to recover your account")
		fmt.Println("You should export this account and store it in a secure location")
		fmt.Printf("NAME: %s\tADDRESS: 0x%x\n", name, address)
		return nil
	},
}
