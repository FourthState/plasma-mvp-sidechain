package utils

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

func ToEthSignedMessageHash(msg []byte) []byte {
	buffer := new(bytes.Buffer)
	msgWithPrefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%v", len(msg))
	buffer.Write([]byte(msgWithPrefix))
	buffer.Write(msg)

	return crypto.Keccak256(buffer.Bytes())
}
