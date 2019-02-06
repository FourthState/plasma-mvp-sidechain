package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addKeyCmd)
}

var addKeyCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new account",
	Long:  `Add an encrypted account to the keystore.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := keystore.NewAccount()
		if err != nil {
			return err
		}

		fmt.Println("\n**Important** do not lose your passphrase.")
		fmt.Println("It is the only way to recover your account")
		fmt.Println("You should export this account and store it in a secure location")
		fmt.Printf("Your new account address is: %s\n", address.Hex())
		return nil
	},
}
