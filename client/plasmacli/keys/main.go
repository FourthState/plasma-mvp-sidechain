package keys

import (
	"github.com/FourthState/plasma-mvp-sidechain/client/store"
	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage local private keys",
	Long:  `Keys allows you to manage your local keystore.`,
}

func KeysCmd() *cobra.Command {
	store.InitKeystore()
	return keysCmd
}
