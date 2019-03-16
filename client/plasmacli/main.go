package main

import (
	"fmt"
	config "github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/eth"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/keys"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/query"
	"github.com/FourthState/plasma-mvp-sidechain/client/store"
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

var rootCmd = &cobra.Command{
	Use:   "plasmacli",
	Short: "Plasma Client",
}

func main() {
	cobra.EnableCommandSorting = false

	rootCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	rootCmd.PersistentFlags().StringP(store.DirFlag, "d", homeDir, "directory for plasmacli")
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Println(err)
	}

	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		fmt.Println(err)
	}

	viper.Set(client.FlagTrustNode, true)
	viper.Set(client.FlagListenAddr, "tcp://localhost:1317")

	viper.AddConfigPath(homeDir)
	plasmaDir := filepath.Join(homeDir, "plasma.toml")
	if _, err := os.Stat(plasmaDir); os.IsNotExist(err) {

		config.WritePlasmaConfigFile(plasmaDir, config.DefaultPlasmaConfig())
	}

	rootCmd.AddCommand(
		eth.EthCmd(),
		query.QueryCmd(),
		eth.ProveCmd(),
		client.LineBreak,
		signCmd,
		spendCmd,
		client.LineBreak,
		keys.KeysCmd(),
		client.LineBreak,
		versionCmd,
	)

	Execute()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// initConfig to be ran when Execute is called
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	viper.AutomaticEnv()
}
