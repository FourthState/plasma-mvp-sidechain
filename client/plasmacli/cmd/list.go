package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(listKeysCmd)
}

var listKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long:  "Return a list of all account addresses stored by this keystore",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := viper.GetString(client.FlagHomeDir)
		keystore.InitKeystore(dir)

		fmt.Println("Index:\t\tAddress:")
		accounts := keystore.Accounts()
		for i, acc := range accounts {
			fmt.Printf("%d\t\t%s\n", i, acc.Address.Hex())
		}

		return nil
	},
}
