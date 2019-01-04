package app

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
)

func SetPlasmaOptions(isOperator bool, privKey, contractAddr, nodeURL, finality string) func(*PlasmaMVPChain) {
	var privateKey *ecdsa.PrivateKey
	var blockFinality uint64

	if isOperator {
		d, err := hex.DecodeString(privKey)
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

	blockFinality, err := strconv.ParseUint(finality, 10, 64)
	if err != nil {
		panic(err)
	}

	return func(pc *PlasmaMVPChain) {
		pc.operatorPrivateKey = privateKey
		pc.isOperator = isOperator
		pc.plasmaContractAddress = common.HexToAddress(contractAddr)
		pc.nodeURL = nodeURL
		pc.blockFinality = blockFinality
	}
}
