package store

import (
	"bytes"
)

const (
	separator = "::"
)

func prefixKey(prefix string, key []byte) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(prefix))
	buffer.Write([]byte(separator))
	buffer.Write(key)

	return buffer.Bytes()
}
