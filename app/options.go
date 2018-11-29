package app

import (
	"github.com/ethereum/go-ethereum/crypto"
	"path/filepath"
)

func SetEthPrivKey(privkey_file string) func(*ChildChain) {
	path, err := filepath.Abs(privkey_file)
	if err != nil {
		panic(err)
	}

	privkey, err := crypto.LoadECDSA(path)
	if err != nil {
		panic(err)
	}
	return func(cc *ChildChain) {
		cc.validatorPrivKey = privkey
	}
}
