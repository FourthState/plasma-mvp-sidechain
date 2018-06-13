package cmd

import (
	"github.com/spf13/cobra"

	"github.com/FourthState/plasma-mvp-sidechain/client"
)

func init() {
	rootCmd.AddCommand(updateKeysCmd)
}

var updateKeysCmd = &cobra.Command{
	Use:   "update <name",
	Short: "Change the password used to protect private key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		buf := client.BufferStdin()
		oldpass, err := client.GetPassword("Enter the current passphrase:", buf)
		if err != nil {
			return err
		}
		newpass, err := client.GetCheckPassword("Enter the new passphrase:", "Repeat the new passphrase:", buf)
		if err != nil {
			return err
		}

		kb, err := client.GetKeyBase()
		if err != nil {
			return err
		}
		err = kb.Update(name, oldpass, newpass)
		if err != nil {
			return err
		}
		fmt.Println("Password successfully updated!")
		return nil
	},
}
