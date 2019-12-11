package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNilAddressDetector(t *testing.T) {
	nilAddress := common.Address{}
	require.True(t, IsZeroAddress(nilAddress), "marked 0x0 as not the nil address")

	addr := common.HexToAddress("1")
	require.False(t, IsZeroAddress(addr), "marked a non 0x0 address as nil")
}

func TestHexPrefixRemoval(t *testing.T) {
	str := "0x0123"
	require.Equal(t, "0123", RemoveHexPrefix(str))

	str = "0x123"
	require.Equal(t, "0123", RemoveHexPrefix(str))

	str = "0123"
	require.Equal(t, "0123", RemoveHexPrefix(str))

	str = "123"
	require.Equal(t, "0123", RemoveHexPrefix(str))

	require.Equal(t, "", RemoveHexPrefix(""))
}
