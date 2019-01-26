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
	hexStrWithPrefix := "0x123"
	require.Equal(t, hexStrWithPrefix[2:], RemoveHexPrefix(hexStrWithPrefix))

	hexStrWithoutPrefix := hexStrWithPrefix[2:]
	require.Equal(t, hexStrWithoutPrefix, RemoveHexPrefix(hexStrWithoutPrefix))

	require.Equal(t, "", RemoveHexPrefix(""))
}
