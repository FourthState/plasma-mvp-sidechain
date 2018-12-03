package eth

import (
	"bytes"
)

const (
	prefixSeperator       = "::"
	depositPrefix         = "deposit"
	transactionExitPrefix = "txExit"
	depositExitPrefix     = "depositExit"
)

func prefixKey(prefix string, key []byte) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(prefix))
	buffer.Write([]byte(prefixSeperator))
	buffer.Write(key)
	return buffer.Bytes()
}
