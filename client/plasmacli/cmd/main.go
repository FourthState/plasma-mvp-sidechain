package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/cmd/keys"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "plasmacli",
	Short: "Plasma Client",
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
	viper.BindPFlags(rootCmd.Flags())

	// TODO: Fix this
	viper.Set(client.FlagTrustNode, true)
	viper.Set(client.FlagListenAddr, "tcp://localhost:1317")

	rootCmd.AddCommand(keys.KeysCmd())
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	viper.AutomaticEnv()
}
