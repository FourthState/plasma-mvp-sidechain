package utils

import (
	ecdsa "crypto/ecdsa"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	crypto "github.com/tendermint/go-crypto"
	"math/big"
)

func ZeroAddress(addr crypto.Address) bool {
	return new(big.Int).SetBytes(addr.Bytes()).Sign() == 0
}

func ValidAddress(addr crypto.Address) bool {
	return new(big.Int).SetBytes(addr.Bytes()).Sign() != 0 && len(addr) == 20
}

func EthPrivKeyToSDKAddress(p *ecdsa.PrivateKey) crypto.Address {
	return ethcrypto.PubkeyToAddress(ecdsa.PublicKey(p.PublicKey)).Bytes()
}

func GenerateAddress() crypto.Address {
	priv, err := ethcrypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	return EthPrivKeyToSDKAddress(priv)
}
