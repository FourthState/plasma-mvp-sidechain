package keys

import (
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
)

var accounts *leveldb.DB

func KeysCmd() *cobra.Command {
	keysCmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage local private keys",
		Long: `Keys allows you to add, import, delete, and view your local keystore
		
The default keystore location is $HOME/.plasmacli/keys`,
	}
	keysCmd.AddCommand(
		addCmd,
		deleteCmd,
		importCmd,
		listCmd,
		updateCmd,
	)
	keystore.InitKeystore(os.ExpandEnv("$HOME/.plasmacli/keys"))

	return keysCmd
}
