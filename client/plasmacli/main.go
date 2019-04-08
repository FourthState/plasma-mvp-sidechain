package main

import (
	"fmt"
	config "github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/eth"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/keys"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/query"
	plasmarest "github.com/FourthState/plasma-mvp-sidechain/client/rest"
	"github.com/FourthState/plasma-mvp-sidechain/client/store"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// default directory
var homeDir string = os.ExpandEnv("$HOME/.plasmacli/")

// Flags
const (
	flagAccount      = "accounts"
	flagAddress      = "address"
	flagOwner        = "owner"
	flagPositions    = "position"
	flagConfirmSigs0 = "Input0ConfirmSigs"
	flagConfirmSigs1 = "Input1ConfirmSigs"
	flagInputs       = "inputValues"
	flagSync         = "sync"
	flagFee          = "fee"
	flagReplay       = "replay"
)

var rootCmd = &cobra.Command{
	Use:   "plasmacli",
	Short: "Plasma Client",
}

func main() {
	Execute()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func MkCodec() *codec.Codec {
	var cdc = codec.New()
	codec.RegisterCrypto(cdc)
	plasmarest.RegisterCodec(cdc)
	return cdc
}

func registerRoutes(rs *lcd.RestServer) {
	plasmarest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
}

func init() {
	// initConfig to be ran when Execute is called
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	rootCmd.PersistentFlags().StringP(store.DirFlag, "d", homeDir, "directory for plasmacli")
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Println(err)
	}

	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		fmt.Println(err)
	}

	viper.Set(client.FlagTrustNode, true)
	viper.Set(client.FlagListenAddr, "tcp://0.0.0.0:1317")

	viper.AddConfigPath(homeDir)
	plasmaDir := filepath.Join(homeDir, "plasma.toml")
	if _, err := os.Stat(plasmaDir); os.IsNotExist(err) {
		config.WritePlasmaConfigFile(plasmaDir, config.DefaultPlasmaConfig())
	}

	cdc := MkCodec()

	rootCmd.AddCommand(
		keys.KeysCmd(),
		query.QueryCmd(),
		eth.EthCmd(),
		eth.ProveCmd(),
		lcd.ServeCommand(cdc, registerRoutes),
	)
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	viper.AutomaticEnv()
}
