package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/spf13/cobra"

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
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
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
	return app.NewChildChain(logger, db, traceStore)
}

// non-functional
func exportAppState(logger log.Logger, db dbm.DB, traceStore io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	papp := app.NewChildChain(logger, db, traceStore)
	return papp.ExportAppStateJSON()
}
