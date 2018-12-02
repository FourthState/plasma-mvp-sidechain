package app

import (
	"crypto/ecdsa"
	"fmt"
	"path/filepath"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func SetEthConfig(isValidator bool, privkey_file string, rootchain_addr string) func(*ChildChain) {
	var privkey *ecdsa.PrivateKey
	var rootchain ethcmn.Address
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

		rootchain = ethcmn.HexToAddress(rootchain_addr)
	}
	return func(cc *ChildChain) {
		cc.validatorPrivKey = privkey
		cc.isValidator = isValidator
		cc.rootchain = rootchain
	}
}
