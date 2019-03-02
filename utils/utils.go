package utils

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	Big0        = big.NewInt(0)
	Big1        = big.NewInt(1)
	Big2        = big.NewInt(2)
	ZeroAddress = common.Address{}
)

// IsZeroAddress is an indicator if the address is the "0x0" address
func IsZeroAddress(addr common.Address) bool {
	return bytes.Equal(addr[:], ZeroAddress[:])
}

func RemoveHexPrefix(hexStr string) string {
	if len(hexStr) < 2 {
		return hexStr
	} else if hexStr[:2] == "0x" {
		return hexStr[2:]
	}

	return hexStr
}
