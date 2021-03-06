package store

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// keys
var (
	walletKey      = []byte{0x0}
	depositKey     = []byte{0x1}
	feeKey         = []byte{0x2}
	txKey          = []byte{0x3}
	outputKey      = []byte{0x4}
	blockKey       = []byte{0x5}
	blockHeightKey = []byte{0x6}
)

// GetWalletKey returns the key to retrieve wallet for given address.
func GetWalletKey(addr common.Address) []byte {
	return prefixKey(walletKey, addr.Bytes())
}

// GetDepositKey returns the key to retrieve deposit for given nonce.
func GetDepositKey(nonce *big.Int) []byte {
	return prefixKey(depositKey, nonce.Bytes())
}

// GetFeeKey returns the key to retrieve fee for given position.
func GetFeeKey(pos plasma.Position) []byte {
	return prefixKey(feeKey, pos.Bytes())
}

// GetOutputKey returns key to retrieve Output for given position.
func GetOutputKey(pos plasma.Position) []byte {
	return prefixKey(outputKey, pos.Bytes())
}

// GetTxKey returns key to retrieve Transaction for given hash.
func GetTxKey(hash []byte) []byte {
	return prefixKey(txKey, hash)
}

// GetBlockKey returns the key for the specified height
func GetBlockKey(height *big.Int) []byte {
	return prefixKey(blockKey, height.Bytes())
}

// GetBlockHeightKey returns the key for the height counter
func GetBlockHeightKey() []byte {
	return blockHeightKey
}

func prefixKey(prefix, key []byte) []byte {
	return append(prefix, key...)
}
