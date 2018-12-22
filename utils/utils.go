package utils

import (
	"github.com/ethereum/go-ethereum/common"
)

// IsZeroAddress is an indicator if the address is the "0x0" address
func IsZeroAddress(addr common.Address) bool {
	return addr.Big().Sign() == 0
}
