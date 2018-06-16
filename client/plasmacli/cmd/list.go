package cmd

import (
	"fmt"

	"github.com/FourthState/plasma-mvp-sidechain/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(listKeysCmd)
}

var listKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long:  "Return a list of all accounts stored by this keystore",
	RunE: func(cmd *cobra.Command, args []string) error {

		dir := viper.GetString(FlagHomeDir)
		ks := client.GetKeyStore(dir)

		accounts := ks.Accounts()
		for _, acc := range accounts {
			// TODO: Create nice printing format
			fmt.Println()
			fmt.Println(acc)
			fmt.Println()
		}
		return nil
	},
}
