package app

import (
	"crypto/ecdsa"
	"fmt"
	"path/filepath"
	"strconv"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func SetEthConfig(isValidator bool, privkey_file, rootchain_addr, nodeURL, minFees, finality string) func(*ChildChain) {
	var privkey *ecdsa.PrivateKey
	var rootchain ethcmn.Address
	var min_fees uint64
	var block_finality uint64

	if isValidator {
		path, err := filepath.Abs(privkey_file)
		if err != nil {
			errMsg := fmt.Sprintf("Could not resolve provided private key file path: %v", err)
			panic(errMsg)
		}

		privkey, err = crypto.LoadECDSA(path)
		if err != nil {
			errMsg := fmt.Sprintf("Could not load provided private key file to ecdsa private key: %v", err)
			panic(errMsg)
		}

		min_fees, err = strconv.ParseUint(minFees, 10, 64)
		if err != nil {
			panic(err)
		}

		block_finality, err = strconv.ParseUint(finality, 10, 64)
		if err != nil {
			panic(err)
		}
	}
	rootchain = ethcmn.HexToAddress(rootchain_addr)

	return func(cc *ChildChain) {
		cc.validatorPrivKey = privkey
		cc.isValidator = isValidator
		cc.rootchain = rootchain
		cc.nodeURL = nodeURL
		cc.min_fees = min_fees
		cc.block_finality = block_finality
	}
}
