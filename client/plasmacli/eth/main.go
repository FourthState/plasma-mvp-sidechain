package eth

import (
	"encoding/hex"
	"fmt"
	config "github.com/FourthState/plasma-mvp-sidechain/client"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	contracts "github.com/FourthState/plasma-mvp-sidechain/contracts/wrappers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"os"
	"path/filepath"
)

type Plasma struct {
	ec *ethclient.Client

	operatorSession *contracts.PlasmaMVPSession
	contract        *contracts.PlasmaMVP

	nodeURL string // current url set for the eth client
}

var rc Plasma

var ethCmd = &cobra.Command{
	Use:               "eth",
	Short:             "Interact with Ethereum rootchain contract",
	PersistentPreRunE: persistentPreRunEFn(),
}

func EthCmd() *cobra.Command {

	return ethCmd
}

// Parse plasma.toml before every call to the eth command
func persistentPreRunEFn() func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// custom plasma config
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
			return fmt.Errorf("please specify a node url for eth connection")
		}

		return nil
	}
}

// Create a connection to an eth node based on the params set
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

	d, err := hex.DecodeString(conf.EthPrivateKey)
	if err != nil {
		return fmt.Errorf("Could not parse private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(d)
	if err != nil {
		return fmt.Errorf("Could not load the private key: %v", err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	operatorSession := &contracts.PlasmaMVPSession{
		Contract: plasmaContract,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: 3141592, // aribitrary
			Value:    big.NewInt(0),
		},
	}

	rc = Plasma{
		ec:              ec,
		operatorSession: operatorSession,
		contract:        plasmaContract,
		nodeURL:         conf.EthNodeURL,
	}

	return nil
}
