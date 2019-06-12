package keys

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags
const (
	nameF = "name"
)

func init() {
	keysCmd.AddCommand(updateCmd)
	updateCmd.Flags().String(nameF, "", "updated key name.")
	viper.BindPFlags(updateCmd.Flags())
}

var updateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "Update account passphrase or key name",
	Long: `Update local encrypted private keys to be encrypted with the new passphrase.

Usage:
	plasmacli keys update <name>
	plasmacli keys update <name> --name <updatedName>
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		updatedName := viper.GetString(nameF)
		msg, err := store.UpdateAccount(name, updatedName)
		if err != nil {
			return err
		}

		fmt.Println(msg)
		return nil
	},
}
