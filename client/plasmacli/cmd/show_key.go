package cmd

import (
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/spf13/cobra"
	keys "github.com/tendermint/go-crypto/keys"
)

func init() {
	rootCmd.AddCommand(showKeysCmd)
}

var showKeysCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show key info for the given name",
	Long:  `Return public details of on local key.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		info, err := getKey(name)
		if err == nil {
			client.printInfo(info)
		}
		return err
	},
}

func getKey(name string) (keys.Info, error) {
	kb, err := client.GetKeyBase()
	if err != nil {
		return keys.Info{}, err
	}
	return kb.Get(name)
}
