package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homeDir string = os.ExpandEnv("$HOME/.plasmacli/keys")

const (
	FlagHomeDir = "home"
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
	rootCmd.PersistentFlags().StringP(FlagHomeDir, "", homeDir, "directory for keystore")
	viper.BindPFlags(rootCmd.Flags())
	viper.Set(FlagHomeDir, homeDir)
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	viper.AutomaticEnv()
}
