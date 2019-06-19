package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/cmd/eth"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/cmd/keys"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/cmd/query"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// default directory
var homeDir string = os.ExpandEnv("$HOME/.plasmacli/")

// Flags
const (
	accountF      = "accounts"
	addressF      = "address"
	asyncF        = "async"
	confirmSigs0F = "Input0ConfirmSigs"
	confirmSigs1F = "Input1ConfirmSigs"
	feeF          = "fee"
	inputsF       = "inputValues"
	ownerF        = "owner"
	positionF     = "position"
	replayF       = "replay"
)

func init() {
	// initConfig to be ran when Execute is called
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	viper.AutomaticEnv()
}

var rootCmd = &cobra.Command{
	Use:           "plasmacli",
	Short:         "Plasma Client",
	SilenceErrors: true,
}

func RootCmd() *cobra.Command {
	cobra.EnableCommandSorting = false
	rootCmd.PersistentFlags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	rootCmd.PersistentFlags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
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
		eth.EthCmd(),
		query.QueryCmd(),
		client.LineBreak,

		RestServerCmd(),
		client.LineBreak,

		SignCmd(),
		SpendCmd(),
		IncludeCmd(),
		client.LineBreak,

		keys.KeysCmd(),
		client.LineBreak,

		VersionCmd(),
	)

	return rootCmd
}
