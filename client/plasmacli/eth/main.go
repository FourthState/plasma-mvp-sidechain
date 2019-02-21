package eth

import (
	"fmt"
	"os"
	"path/filepath"

	config "github.com/FourthState/plasma-mvp-sidechain/client"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	contracts "github.com/FourthState/plasma-mvp-sidechain/contracts/wrappers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags
const (
	allF        = "all"
	limitF      = "limit"
	feeF        = "fee"
	txBytesF    = "tx-bytes"
	proofF      = "proof"
	sigsF       = "signatures"
	trustNodeF  = "trust-node"
	depositsF   = "deposits"
	indexF      = "index"
	accountF    = "account"
	addrF       = "address"
	minExitBond = 200000
)

type Plasma struct {
	ec *ethclient.Client

	session  *contracts.PlasmaMVPSession
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
func persistentPreRunEFn() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		plasmaConfigFilePath := filepath.Join(viper.GetString(ks.DirFlag), "plasma.toml")

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
		if rc.nodeURL != conf.EthNodeURL {
			if err := initEthConn(conf); err != nil {
				return err
			}
		} else if conf.EthNodeURL == "" {
			return fmt.Errorf("please specify a node url for eth connection in %s/plasma.toml", viper.GetString(ks.DirFlag))
		}

		return nil
	}
}

// Create a connection to an eth node based on
// the params set in plasma.toml
func initEthConn(conf config.PlasmaConfig) error {
	c, err := rpc.Dial(conf.EthNodeURL)
	if err != nil {
		return err
	}
	ec := ethclient.NewClient(c)

	// Create a session with the contract and operator account
	plasmaContract, err := contracts.NewPlasmaMVP(ethcmn.HexToAddress(conf.EthPlasmaContractAddr), ec)
	if err != nil {
		return err
	}

	session := &contracts.PlasmaMVPSession{
		Contract: plasmaContract,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
	}

	rc = Plasma{
		ec:       ec,
		session:  session,
		contract: plasmaContract,
		nodeURL:  conf.EthNodeURL,
	}

	return nil
}
