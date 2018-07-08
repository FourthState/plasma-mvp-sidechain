package utils

import (
	"reflect"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"math/big"
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
