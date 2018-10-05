package utils

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"reflect"
)

func ZeroAddress(addr common.Address) bool {
	return new(big.Int).SetBytes(addr.Bytes()).Sign() == 0
}

func ValidAddress(addr common.Address) bool {
	return !reflect.DeepEqual(addr, common.Address{})
}

func PrivKeyToAddress(p *ecdsa.PrivateKey) common.Address {
	return ethcrypto.PubkeyToAddress(ecdsa.PublicKey(p.PublicKey))
}

func GenerateAddress() common.Address {
	priv, err := ethcrypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	return PrivKeyToAddress(priv)
}

// helper function for tests
func GetIndex(index int64) int64 {
	if index >= 0 {
		return index
	} else {
		return 0
	}
}
