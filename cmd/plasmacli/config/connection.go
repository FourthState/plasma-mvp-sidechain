package config

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/FourthState/plasma-mvp-sidechain/eth"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strconv"
)

var plasmaContract *eth.Plasma

func GetContractConn() (*eth.Plasma, error) {
	if plasmaContract != nil {
		return plasmaContract, nil
	}

	conn, err := setupContractConn()
	if err != nil {
		return nil, fmt.Errorf("unable to enable contract connection: %s", err)
	}

	plasmaContract = conn
	return plasmaContract, nil
}

func setupContractConn() (*eth.Plasma, error) {
	// Parse plasma.toml before every call to the eth command
	// Update ethereum client connection if params have changed
	plasmaConfigFilePath := filepath.Join(viper.GetString(store.DirFlag), "plasma.toml")

	if _, err := os.Stat(plasmaConfigFilePath); os.IsNotExist(err) {
		plasmaConfig := DefaultPlasmaConfig()
		WritePlasmaConfigFile(plasmaConfigFilePath, plasmaConfig)
	}

	viper.SetConfigName("plasma")
	if err := viper.MergeInConfig(); err != nil {
		return nil, err
	}

	conf, err := ParsePlasmaConfigFromViper()
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
