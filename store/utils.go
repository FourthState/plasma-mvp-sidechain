package store

func prefixKey(prefix, key []byte) []byte {
	return append(prefix, key...)
}
