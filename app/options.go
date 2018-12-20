package app

import (
	"crypto/ecdsa"
	"fmt"
	"path/filepath"
	"strconv"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func SetEthConfig(isValidator bool, privkeyFile, rootchainAddr, nodeURL, minFeesStr, finality string) func(*ChildChain) {
	var privkey *ecdsa.PrivateKey
	var rootchain ethcmn.Address
	var minFees uint64
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

		minFees, err = strconv.ParseUint(minFeesStr, 10, 64)
		if err != nil {
			panic(err)
		}
	}
	blockFinality, err := strconv.ParseUint(finality, 10, 64)
	if err != nil {
		panic(err)
	}
	rootchain = ethcmn.HexToAddress(rootchainAddr)

	return func(cc *ChildChain) {
		cc.validatorPrivKey = privkey
		cc.isValidator = isValidator
		cc.rootchain = rootchain
		cc.nodeURL = nodeURL
		cc.minFees = minFees
		cc.blockFinality = blockFinality
	}
}
