package keys

import (
	cli "github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
)

const (
	AccountDir = "accounts.ldb"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage local private keys",
	Long: `Keys allows you to add, import, delete, and view your local keystore
         
The default keystore location is $HOME/.plasmacli/keys`,
}

func KeysCmd() *cobra.Command {
	keysCmd.AddCommand(
		deleteCmd,
		importCmd,
		updateCmd,
	)
	keystore.InitKeystore(filepath.Join(viper.GetString(cli.DirFlag), "keys"))
	return keysCmd
}
