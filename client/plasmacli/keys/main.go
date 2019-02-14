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
	Long:  `Keys allows you to manage your local keystore.`,
}

func KeysCmd() *cobra.Command {
	keystore.InitKeystore(filepath.Join(viper.GetString(cli.DirFlag), "keys"))
	return keysCmd
}
