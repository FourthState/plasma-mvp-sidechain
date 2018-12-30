package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	app "github.com/FourthState/plasma-mvp-sidechain"
	"github.com/FourthState/plasma-mvp-sidechain/server/plasmad/cmd"
	"github.com/FourthState/plasma-mvp-sidechain/server/plasmad/config"
)

func main() {
	// codec only used for server.AddCommand
	cdc := codec.New()
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
	executor.Execute()
}

// wraps the default cosmos-sdk function with additional logic to handle plasma specific configuration
func persistentPreRunEFn(context *server.Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// run sdk/tendermint configuration
		if err := server.PersistentPreRunE(context)(cmd)(args); err != nil {
			return err
		}

		// custom plasma config
		plasmaConfigFilePath := filepath.Join(context.Config.RootDir, "config/plasma.toml")

		if _, err := os.Stat(plasmaConfigFilePath); os.IsNotExist(err) {
			plasmaConfig, _ := config.DefaultPlasmaConfig()
			config.WritePlasmaConfigFile(plasmaConfigFilePath, plasmaConfig)
		}

		viper.SetConfigName("plasma")
		if err = viper.MergeInConfig(); err != nil {
			return err
		}
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	keyFile := viper.GetString(cli.HomeFlag) + "/config/" + viper.GetString("ethereum_privkey_file")
	isOperator := viper.GetBool("is_operator")
	contractAddr := viper.GetString("ethereum_plasma_contract_addr")
	nodeURL := viper.GetString("ethereum_nodeurl")
	finality := viper.GetString("ethereum_finality")

	return app.NewPlasmaMVPChain(logger, db, traceStore,
		setPlasmaOptions(isOperator, keyFile, contractAddr, nodeURL, finality),
	)
}

func setPlasmaOptions(isValidator bool, privkeyFile, contractAddr, finality string) func(*PlasmaMVPChain) {
	var privkey *ecdsa.PrivateKey
	var contractAddr common.Address
	var blockFinality uint64

	if isValidator {
		path, err := filepath.Abs(privkeyFile)
		if err != nil {
			errMsg := fmt.Sprintf("Could not resolve provided private key file path: %v", err)
			panic(errMsg)
		}

		privkey, err = crypto.LoadECDSA(path)
		if err != nil {
			errMsg := fmt.Sprintf("Could not load provided private key file to ecdsa private key: %v", err)
			panic(errMsg)
		}
	}

	blockFinality, err := strconv.ParseUint(finality, 10, 64)
	if err != nil {
		panic(err)
	}

	return func(pc *PlasmaMVPChain) {
		pc.operatorPrivateKey = privkey
		pc.isOperator = isValidator
		pc.plasmaContractAddress = commmon.HexToAddress(contractAddr)
		pc.nodeURL = nodeURL
		pc.blockFinality = blockFinality
	}
}

// non-functional
func exportAppState(logger log.Logger, db dbm.DB, traceStore io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	papp := app.NewPlasmaMVPChain(logger, db, traceStore)
	return papp.ExportAppStateJSON()
}
