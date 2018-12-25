package utils

import (
	"bytes"
	"github.com/ethereum/go-ethereum/crypto"
)

func ToEthSignedMessageHash(msg [32]byte) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte("\x19Ethereum Signed Message:\n32"))
	buffer.Write(msg[:])
	return crypto.Keccak256(buffer.Bytes())
}
