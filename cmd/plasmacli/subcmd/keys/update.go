package keys

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// UpdateCmd returns the keys update command
func UpdateCmd() *cobra.Command {
	updateCmd.Flags().String(nameF, "", "updated key name.")
	return updateCmd
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
		viper.BindPFlags(cmd.Flags())

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
