package keys

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long:  "Return a list of all account addresses stored by this keystore",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Index:\t\tAddress:")
		keys := keystore.Accounts()
		for i, key := range keys {
			fmt.Printf("%d\t\t%s\n", i, key.Address.Hex())
		}

		return nil
	},
}
