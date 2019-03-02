package eth

import (
	"fmt"
	config "github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/store"
	contracts "github.com/FourthState/plasma-mvp-sidechain/contracts/wrappers"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
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

	minExitBond = 200000 // specified by rootchain contract
)

type Plasma struct {
	ec *ethclient.Client

	contract *contracts.PlasmaMVP

	nodeURL string // current url set for the eth client
}

var rc Plasma

var ethCmd = &cobra.Command{
	Use:   "eth",
	Short: "Interact with Ethereum rootchain contract",
	Long: `Configurations for interacting with the rootchain contract can be specified in <dirpath>/plasma.toml.
An eth node instance needs to be running for this command to work.`,
	PersistentPreRunE: persistentPreRunEFn(),
}

func EthCmd() *cobra.Command {
	return ethCmd
}

// Parse plasma.toml before every call to the eth command
// Update ethereum client connection if params have changed
func persistentPreRunEFn() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		plasmaConfigFilePath := filepath.Join(viper.GetString(store.DirFlag), "plasma.toml")

		if _, err := os.Stat(plasmaConfigFilePath); os.IsNotExist(err) {
			plasmaConfig := config.DefaultPlasmaConfig()
			config.WritePlasmaConfigFile(plasmaConfigFilePath, plasmaConfig)
		}

		viper.SetConfigName("plasma")
		if err := viper.MergeInConfig(); err != nil {
			return err
		}

		conf, err := config.ParsePlasmaConfigFromViper()
		if err != nil {
			return err
		}

		// Check to see if the eth connection params have changed
		if conf.EthNodeURL == "" {
			return fmt.Errorf("please specify a node url for eth connection in %s/plasma.toml", viper.GetString(store.DirFlag))
		} else if rc.nodeURL != conf.EthNodeURL {
			if err := initEthConn(conf); err != nil {
				return err
			}
		}

		return nil
	}
}

// Create a connection to an eth node based on
// the params set in plasma.toml
func initEthConn(conf config.PlasmaConfig) error {
	c, err := rpc.Dial(conf.EthNodeURL)
	if err != nil {
		return fmt.Errorf("failed to dial node url: { %s }", err)
	}
	ec := ethclient.NewClient(c)

	// Bind rootchain contract and operator account
	plasmaContract, err := contracts.NewPlasmaMVP(ethcmn.HexToAddress(conf.EthPlasmaContractAddr), ec)
	if err != nil {
		return fmt.Errorf("failed to bind to contract: { %s }", err)
	}

	rc = Plasma{
		ec:       ec,
		contract: plasmaContract,
		nodeURL:  conf.EthNodeURL,
	}

	return nil
}
