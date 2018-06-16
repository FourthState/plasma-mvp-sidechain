package main

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/server"
	abci "github.com/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/FourthState/plasma-mvp-sidechain/app"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "plasmad",
		Short:             "Plasma Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, server.DefaultAppInit,
		server.ConstructAppCreator(newApp, "plasma"),
		server.ConstructAppExporter(exportAppState, "plasma"))

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.plasmad")
	executor := cli.PrepareBaseCmd(rootCmd, "PC", rootDir)
	executor.Execute()
}

func newApp(logger log.Logger, db dbm.DB) abci.Application {
	return app.NewChildChain(logger, db)
}

// non-functional
func exportAppState(logger log.Logger, db dbm.DB) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	papp := app.NewChildChain(logger, db)
	return papp.ExportAppStateJSON()
}