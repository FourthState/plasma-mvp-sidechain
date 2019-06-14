package eth

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/cmd/eth/query"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/FourthState/plasma-mvp-sidechain/eth"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/cosmos/cosmos-sdk/client"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strconv"
)

const (
	// Flags
	accountF   = "account"
	addrF      = "address"
	allF       = "all"
	depositsF  = "deposits"
	feeF       = "fee"
	gasLimitF  = "gas-limit"
	indexF     = "index"
	limitF     = "limit"
	ownerF     = "owner"
	positionF  = "position"
	proofF     = "proof"
	sigsF      = "signatures"
	trustNodeF = "trust-node"
	txBytesF   = "tx-bytes"
)

var plasmaContract *eth.Plasma

func EthCmd() *cobra.Command {
	ethCmd.AddCommand(
		ProveCmd(),
		ChallengeCmd(),
		ExitCmd(),
		FinalizeCmd(),
		DepositCmd(),
		WithdrawCmd(),
		client.LineBreak,

		query.QueryCmd(plasmaContract),
	)

	return ethCmd
}

var ethCmd = &cobra.Command{
	Use:   "eth",
	Short: "Interact with the plasma smart contract",
	Long: `Configurations for interacting with the rootchain contract can be specified in <dirpath>/plasma.toml.
An eth node instance needs to be running for this command to work.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		plasma, err := setupPlasmaConn()
		if err != nil {
			return err
		}

		plasmaContract = plasma
		return nil
	},
}

func setupPlasmaConn() (*eth.Plasma, error) {
	if plasmaContract != nil {
		return nil, nil
	}

	// Parse plasma.toml before every call to the eth command
	// Update ethereum client connection if params have changed
	plasmaConfigFilePath := filepath.Join(viper.GetString(store.DirFlag), "plasma.toml")

	if _, err := os.Stat(plasmaConfigFilePath); os.IsNotExist(err) {
		plasmaConfig := config.DefaultPlasmaConfig()
		config.WritePlasmaConfigFile(plasmaConfigFilePath, plasmaConfig)
	}

	viper.SetConfigName("plasma")
	if err := viper.MergeInConfig(); err != nil {
		return nil, err
	}

	conf, err := config.ParsePlasmaConfigFromViper()
	if err != nil {
		return nil, err
	}

	// Check to see if the eth connection params have changed
	dir := viper.GetString(store.DirFlag)
	if conf.EthNodeURL == "" {
		return nil, fmt.Errorf("please specify a node url for eth connection in %s/plasma.toml", dir)
	} else if conf.EthPlasmaContractAddr == "" || !ethcmn.IsHexAddress(conf.EthPlasmaContractAddr) {
		return nil, fmt.Errorf("please specic a valid contract address in %s/plasma.toml", dir)
	}

	ethClient, err := eth.InitEthConn(conf.EthNodeURL)
	if err != nil {
		return nil, err
	}
	blockFinality, err := strconv.ParseUint(conf.EthBlockFinality, 10, 64)
	if err != nil {
		return nil, err
	}
	plasma, err := eth.InitPlasma(ethcmn.HexToAddress(conf.EthPlasmaContractAddr), ethClient, blockFinality)
	if err != nil {
		return nil, err
	}

	return plasma, nil
}

func HasTxExited(pos plasma.Position) (bool, error) {
	p, err := setupPlasmaConn()
	if err != nil {
		return true, err
	}

	exited, err := p.HasTxBeenExited(nil, pos)
	if err != nil {
		return true, err
	}

	return exited, nil
}
