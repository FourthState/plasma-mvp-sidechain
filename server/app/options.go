package app

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/server/plasmad/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
)

func SetPlasmaOptionsFromConfig(conf config.PlasmaConfig) func(*PlasmaMVPChain) {
	var privateKey *ecdsa.PrivateKey
	var blockFinality uint64

	if conf.IsOperator {
		d, err := hex.DecodeString(conf.EthOperatorPrivateKey)
		if err != nil {
			errMsg := fmt.Sprintf("Could not parse private key: %v", err)
			panic(errMsg)
		}

		privateKey, err = crypto.ToECDSA(d)
		if err != nil {
			errMsg := fmt.Sprintf("Could not load the private key: %v", err)
			panic(errMsg)
		}
	}

	blockFinality, err := strconv.ParseUint(conf.EthBlockFinality, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("Could not parse block finality: %v", err)
		panic(errMsg)
	}

	fmt.Println("contract addr: ", conf.EthPlasmaContractAddr)
	if !common.IsHexAddress(conf.EthPlasmaContractAddr) {
		panic("invalid contract address. please use hex format")
	}

	return func(pc *PlasmaMVPChain) {
		pc.operatorPrivateKey = privateKey
		pc.isOperator = conf.IsOperator
		pc.plasmaContractAddress = common.HexToAddress(conf.EthPlasmaContractAddr)
		pc.nodeURL = conf.EthNodeURL
		pc.blockFinality = blockFinality
	}
}
