package app

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"path/filepath"
)

func SetEthPrivKey(privkey_file string, isValidator bool) func(*ChildChain) {
	var privkey *ecdsa.PrivateKey
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
	}
	return func(cc *ChildChain) {
		cc.validatorPrivKey = privkey
		cc.isValidator = isValidator
	}
}
