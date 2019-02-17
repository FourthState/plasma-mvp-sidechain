package main

import (
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/keys"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/query"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// default directory
var homeDir string = os.ExpandEnv("$HOME/.plasmacli/")

// Flags
const (
	flagAccount      = "accounts"
	flagOwner        = "owner"
	flagPositions    = "position"
	flagConfirmSigs0 = "Input0ConfirmSigs"
	flagConfirmSigs1 = "Input1ConfirmSigs"
	flagInputs       = "inputValues"
	flagSync         = "sync"
	flagFee          = "fee"
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

func init() {
	// initConfig to be ran when Execute is called
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	rootCmd.PersistentFlags().StringP(ks.DirFlag, "d", homeDir, "directory for plasmacli")
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Println(err)
	}

	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		fmt.Println(err)
	}

	viper.Set(client.FlagTrustNode, true)
	viper.Set(client.FlagListenAddr, "tcp://localhost:1317")

	rootCmd.AddCommand(keys.KeysCmd(), query.QueryCmd())
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	viper.AutomaticEnv()
}
