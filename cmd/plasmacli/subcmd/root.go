package subcmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/subcmd/eth"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/subcmd/keys"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/subcmd/query"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/subcmd/tx"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// default directory
var homeDir string = os.ExpandEnv("$HOME/.plasmacli/")

// Flags
const ()

func init() {
	// initConfig to be ran when Execute is called
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	viper.AutomaticEnv()
}

func RootCmd() *cobra.Command {
	cobra.EnableCommandSorting = false
	rootCmd.PersistentFlags().StringP(store.DirFlag, "d", homeDir, "directory for plasmacli")
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(homeDir)
	plasmaDir := filepath.Join(homeDir, "plasma.toml")
	if _, err := os.Stat(plasmaDir); os.IsNotExist(err) {
		config.WritePlasmaConfigFile(plasmaDir, config.DefaultPlasmaConfig())
	}

	rootCmd.AddCommand(
		tx.TxCmd(),
		eth.EthCmd(),
		query.QueryCmd(),
		client.LineBreak,

		RestServerCmd(),
		client.LineBreak,

		keys.KeysCmd(),
		client.LineBreak,

		VersionCmd(),
	)

	return rootCmd
}

var rootCmd = &cobra.Command{
	Use:           "plasmacli",
	Short:         "Plasma Client",
	SilenceErrors: true,
}
