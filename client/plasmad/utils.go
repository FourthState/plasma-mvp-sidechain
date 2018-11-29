package main

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/FourthState/plasma-mvp-sidechain/client/plasmad/config"
	"github.com/cosmos/cosmos-sdk/server"
	tmcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	tmcfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
)

// Returns a PresistentPreRunE function that initaillizes a config object
// using the config files stored locally. Parses plasma.toml, allowing access
// to the validators plasma configurations
func PersistentPreRunEFn(context *server.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		config, err := loadConfig()
		if err != nil {
			return err
		}
		err = validateConfig(config)
		if err != nil {
			return err
		}
		logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
		logger, err = tmflags.ParseLogLevel(config.LogLevel, logger, tmcfg.DefaultLogLevel())
		if err != nil {
			return err
		}

		if viper.GetBool(cli.TraceFlag) {
			logger = log.NewTracingLogger(logger)
		}
		logger = logger.With("module", "main")
		context.Config = config
		context.Logger = logger
		return nil
	}
}

// update default tendermint and child chain settings if config is created
func loadConfig() (conf *tmcfg.Config, err error) {
	// use a tmpConf to get root dir
	tmpConf := tmcfg.DefaultConfig()
	err = viper.Unmarshal(tmpConf)
	if err != nil {
		return nil, errors.New("error in unmarshalling default tendermint config object")
	}
	rootDir := tmpConf.RootDir
	configFilePath := filepath.Join(rootDir, "config/config.toml")

	// check if the config.toml exists
	// if it does not exist, create a config file and set default configurations
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		conf, _ = tmcmd.ParseConfig()
		conf.ProfListenAddress = "listenhost:6060"
		conf.P2P.RecvRate = 5120000
		conf.P2P.SendRate = 5120000
		conf.TxIndex.IndexAllTags = true
		conf.Consensus.TimeoutCommit = 5 * time.Second
		tmcfg.WriteConfigFile(configFilePath, conf)
	}

	// parse existing file
	if conf == nil {
		conf, err = tmcmd.ParseConfig()
	}

	plasmaConfigFilePath := filepath.Join(rootDir, "config/plasma.toml")
	viper.SetConfigName("plasma")
	_ = viper.MergeInConfig()
	var plasmaConf *config.Config

	// check if plasma.toml exists
	// if it does not exist, create the config file and set default configurations
	if _, err := os.Stat(plasmaConfigFilePath); os.IsNotExist(err) {
		plasmaConf, _ := config.ParseConfig()
		config.WriteConfigFile(plasmaConfigFilePath, plasmaConf)
	}

	if plasmaConf == nil {
		_, err = config.ParseConfig()
	}

	return
}

func validateConfig(conf *tmcfg.Config) error {
	if conf.Consensus.CreateEmptyBlocks == false {
		return errors.New("config option CreateEmptyBlocks = false is currently unsupported")
	}
	return nil
}
