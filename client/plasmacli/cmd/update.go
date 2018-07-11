package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/ethereum/go-ethereum/accounts"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(updateKeysCmd)
}

var updateKeysCmd = &cobra.Command{
	Use:   "update <address>",
	Short: "Change the password used to protect the account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		addrStr := args[0]
		addr, err := client.StrToAddress(addrStr)
		if err != nil {
			return err
		}

		buf := client.BufferStdin()
		oldpass, err := client.GetPassword("Enter the current passphrase:", buf)
		if err != nil {
			return err
		}
		newpass, err := client.GetCheckPassword("Enter the new passphrase:", "Repeat the new passphrase:", buf)
		if err != nil {
			return err
		}

		dir := viper.GetString(FlagHomeDir)
		ks := client.GetKeyStore(dir)

		acc := accounts.Account{
			Address: addr,
		}
		err = ks.Update(acc, oldpass, newpass)
		if err != nil {
			return err
		}
		fmt.Println("Password successfully updated!")
		return nil
	},
}
