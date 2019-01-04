package main

import (
	"encoding/json"
	"github.com/FourthState/plasma-mvp-sidechain/server/plasmad/app"
	"github.com/FourthState/plasma-mvp-sidechain/server/plasmad/cmd"
	"github.com/FourthState/plasma-mvp-sidechain/server/plasmad/config"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// codec only used for server.AddCommand
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "plasmad",
		Short:             "Plasma Daemon (server)",
		PersistentPreRunE: persistentPreRunEFn(ctx),
	}
	rootCmd.AddCommand(cmd.InitCmd(ctx, cdc))
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppState)

	// HomeFlag in tendermint cli will be set to `~/.plasmad`
	rootDir := os.ExpandEnv("$HOME/.plasmad")
	executor := cli.PrepareBaseCmd(rootCmd, "PD", rootDir)
	if err := executor.Execute(); err != nil {
		panic(err)
	}
}

// wraps the default cosmos-sdk function with additional logic to handle plasma specific configuration
func persistentPreRunEFn(context *server.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// run sdk/tendermint configuration
		if err := server.PersistentPreRunEFn(context)(cmd, args); err != nil {
			return err
		}

		// custom plasma config
		plasmaConfigFilePath := filepath.Join(context.Config.RootDir, "config/plasma.toml")

		if _, err := os.Stat(plasmaConfigFilePath); os.IsNotExist(err) {
			plasmaConfig := config.DefaultPlasmaConfig()
			config.WritePlasmaConfigFile(plasmaConfigFilePath, plasmaConfig)
		}

		viper.SetConfigName("plasma")
		if err := viper.MergeInConfig(); err != nil {
			return err
		}

		return nil
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	keyFile := viper.GetString(cli.HomeFlag) + "/config/" + viper.GetString("ethereum_privkey_file")
	isOperator := viper.GetBool("is_operator")
	contractAddr := viper.GetString("ethereum_plasma_contract_addr")
	nodeURL := viper.GetString("ethereum_nodeurl")
	finality := viper.GetString("ethereum_finality")

	return app.NewPlasmaMVPChain(logger, db, traceStore,
		app.SetPlasmaOptions(isOperator, keyFile, contractAddr, nodeURL, finality),
	)
}

// non-functional
func exportAppState(logger log.Logger, db dbm.DB, traceStore io.Writer, _ int64, _ bool) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	papp := app.NewPlasmaMVPChain(logger, db, traceStore)
	return papp.ExportAppStateJSON()
}
