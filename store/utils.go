package store

import (
	"bytes"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	ethcmn "github.com/ethereum/go-ethereum/common"
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

func GetUTXOStoreKey(addr ethcmn.Address, pos plasma.Position) []byte {
	return append(addr.Bytes(), pos.Bytes()...)
}

func GetStoreKey(utxo UTXO) []byte {
	return GetUTXOStoreKey(utxo.Output.Owner, utxo.Position)
}
