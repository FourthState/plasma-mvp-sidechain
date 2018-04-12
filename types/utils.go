package types

import (
	crypto "github.com/tendermint/go-crypto"
	"math/big"
)

func ZeroAddress(addr crypto.Address) bool {
	return new(big.Int).SetBytes(addr.Bytes()).Sign() == 0
}
