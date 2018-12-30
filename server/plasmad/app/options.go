package app

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"path/filepath"
	"strconv"
)

func SetPlasmaOptions(isOperator bool, privkeyFile, contractAddr, nodeURL, finality string) func(*PlasmaMVPChain) {
	var privkey *ecdsa.PrivateKey
	var blockFinality uint64

	if isOperator {
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
		pc.isOperator = isOperator
		pc.plasmaContractAddress = common.HexToAddress(contractAddr)
		pc.nodeURL = nodeURL
		pc.blockFinality = blockFinality
	}
}
