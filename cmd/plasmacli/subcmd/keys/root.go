package keys

import (
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/spf13/cobra"
)

// flags
const (
	nameF = "name"
	fileF = "file"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage local private keys",
	Long:  `Keys allows you to manage your local keystore.`,
}

func KeysCmd() *cobra.Command {
	store.InitKeystore()

	keysCmd.AddCommand(
		AddCmd(),
		DeleteCmd(),
		ImportCmd(),
		ListCmd(),
		UpdateCmd(),
	)

	return keysCmd
}
