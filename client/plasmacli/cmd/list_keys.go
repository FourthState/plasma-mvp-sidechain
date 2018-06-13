package cmd

import (
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listKeysCmd)
}

var listKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List all keys",
	Long:  "Return a list of all public keys stored by this key manager along with their associated name and address",
	RuneE: func(cmd *cobra.Command, args []string) error {
		kb, err := client.GetKeyBase()
		if err != nil {
			return err
		}

		infoList, err := kb.List()
		if err == nil {
			for _, info := range infoList {
				client.printInfo(info)
			}
		}
		return err
	},
}
