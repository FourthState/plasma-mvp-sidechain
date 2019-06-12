package keys

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/spf13/cobra"
)

func init() {
	keysCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete the given address",
	Long:  `Deletes the account from the keystore if the passphrase is correct.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if err := store.DeleteAccount(name); err != nil {
			return err
		}

		fmt.Println("Account deleted.")
		return nil
	},
}
