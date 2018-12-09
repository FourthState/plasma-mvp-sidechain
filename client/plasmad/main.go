package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/server"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/FourthState/plasma-mvp-sidechain/app"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmad/cmd"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "plasmad",
		Short:             "Plasma Daemon (server)",
		PersistentPreRunE: PersistentPreRunEFn(ctx),
	}
	rootCmd.AddCommand(cmd.InitCmd(ctx, cdc, app.PlasmaAppInit()))

	server.AddCommands(ctx, cdc, rootCmd, app.PlasmaAppInit(),
		newApp,
		exportAppState)

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.plasmad")
	executor := cli.PrepareBaseCmd(rootCmd, "PC", rootDir)
	executor.Execute()
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	key_file := viper.GetString("ethereum_privkey_file")
	isValidator := viper.GetBool("is_validator")
	rootchain := viper.GetString("ethereum_rootchain")
	nodeURL := viper.GetString("ethereum_nodeurl")
	minFees := viper.GetString("minimum_fees")
	key_file = viper.GetString(cli.HomeFlag) + "/config/" + key_file
	finality := viper.GetString("ethereum_finality")

	return app.NewChildChain(logger, db, traceStore,
		app.SetEthConfig(isValidator, key_file, rootchain, nodeURL, minFees, finality),
	)
}

// non-functional
func exportAppState(logger log.Logger, db dbm.DB, traceStore io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	papp := app.NewChildChain(logger, db, traceStore)
	return papp.ExportAppStateJSON()
}
