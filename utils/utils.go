package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	Big0 = big.NewInt(0)
	Big1 = big.NewInt(1)
)

// IsZeroAddress is an indicator if the address is the "0x0" address
func IsZeroAddress(addr common.Address) bool {
	return addr.Big().Sign() == 0
}
