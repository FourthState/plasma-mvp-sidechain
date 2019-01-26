package utils

import (
	"bytes"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEthSignedMessageHash(t *testing.T) {
	hash := crypto.Keccak256([]byte("fourthstate"))

	ethSignedMessageHash := ToEthSignedMessageHash(hash)
	expectedMessage := crypto.Keccak256(append([]byte("\x19Ethereum Signed Message:\n32"), hash...))

	require.True(t,
		bytes.Equal(expectedMessage, ethSignedMessageHash),
		"did not create the appropriate ethereum signed hash")
}
