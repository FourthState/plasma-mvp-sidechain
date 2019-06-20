package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/app"
	pConfig "github.com/FourthState/plasma-mvp-sidechain/cmd/plasmad/config"
	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmConfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmCommon "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/privval"
	"os"
	"path/filepath"
)

const (
	flagOverwrite = "overwrite"
	flagMoniker   = "moniker"
	flagChainID   = "chainId"
	flagTest      = "test"
)

type chainInfo struct {
	Moniker    string          `json:"moniker"`
	ChainID    string          `json:"chain_id"`
	NodeID     string          `json:"node_id"`
	AppMessage json.RawMessage `json:"app_message"`
}

// get cmd to initialize all files for tendermint and application
// nolint
func InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			overwrite := viper.GetBool(flagOverwrite)

			/* Tendermint configuration */

			chainID := viper.GetString(flagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", tmCommon.RandStr(7))
			}
			if viper.GetString(flagMoniker) != "" {
				config.Moniker = viper.GetString(flagMoniker)
			}

			nodeID, _, err := gaiaInit.InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}
			var appState json.RawMessage
			genFile := config.GenesisFile()
			if tmCommon.FileExists(genFile) {
				if !overwrite {
					return fmt.Errorf("genesis.json file already exists: %v", genFile)
				} else {
					fmt.Printf("overwriting genesis.json...\n")
				}
			}
			// read of create the private key file for this config
			var privValidator *privval.FilePV
			privValFile := config.PrivValidatorKeyFile()

			if tmCommon.FileExists(privValFile) {
				privValidator = privval.LoadFilePV(privValFile, config.PrivValidatorStateFile())
			} else {
				privValidator = privval.GenFilePV(privValFile, config.PrivValidatorStateFile())
				privValidator.Save()
			}

			valPubKey := privValidator.GetPubKey()

			// create genesis and write to disk
			appState, err = codec.MarshalJSONIndent(cdc, app.NewDefaultGenesisState(valPubKey))
			if err != nil {
				return err
			}

			if err = gaiaInit.ExportGenesisFile(genFile, chainID, nil, appState); err != nil {
				return err
			}

			// write tendermint and plasma config files to disk
			tmConfig.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			var plasmaConfig pConfig.PlasmaConfig
			plasmaPath := filepath.Join(config.RootDir, "config", "plasma.toml")

			if viper.GetBool(flagTest) {
				plasmaConfig = pConfig.TestPlasmaConfig()
			} else {
				plasmaConfig = pConfig.DefaultPlasmaConfig()
			}
			if tmCommon.FileExists(plasmaPath) {
				if !overwrite {
					return fmt.Errorf("plasma.toml already exists at '%s'. Use -o flag to overwrite", plasmaPath)
				} else {
					fmt.Println("overwriting plasma.toml...")
				}
			}
			pConfig.WritePlasmaConfigFile(plasmaPath, plasmaConfig)

			// display chain info
			info, err := json.MarshalIndent(chainInfo{
				ChainID:    chainID,
				Moniker:    config.Moniker,
				NodeID:     nodeID,
				AppMessage: appState,
			}, "", "\t")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", string(info))

			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, os.ExpandEnv("$HOME/.plasmad"), "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().BoolP(flagTest, "t", false, "write default testing configuration")
	cmd.Flags().String(flagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(flagMoniker, "m", "set the validator's moniker")
	return cmd
}
