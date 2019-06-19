package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

// keys
var (
	walletKey  = []byte{0x0}
	depositKey = []byte{0x1}
	feeKey     = []byte{0x2}
	txKey      = []byte{0x3}
	outputKey  = []byte{0x4}
)

/* Returns key to retrieve wallet for given address */
func GetWalletKey(addr common.Address) []byte {
	return prefixKey(walletKey, addr.Bytes())
}

/* Returns key to retrieve deposit for given nonce */
func GetDepositKey(nonce *big.Int) []byte {
	return prefixKey(depositKey, nonce.Bytes())
}

/* Returns key to retrieve fee for given position */
func GetFeeKey(pos plasma.Position) []byte {
	return prefixKey(feeKey, pos.Bytes())
}

/* Returns key to retrieve UTXO for given position */
func GetOutputKey(pos plasma.Position) []byte {
	return prefixKey(outputKey, pos.Bytes())
}

/* Returns Transaction for given hash */
func GetTxKey(hash []byte) []byte {
	return prefixKey(txKey, hash)
}

func prefixKey(prefix, key []byte) []byte {
	return append(prefix, key...)
}
