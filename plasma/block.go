package plasma

import (
	"math/big"
)

// Block represents a plasma block.
type Block struct {
	Header    [32]byte
	TxnCount  uint16
	FeeAmount *big.Int
}

type block struct {
	Header    [32]byte
	TxnCount  uint16
	FeeAmount []byte
}

// NewBlock creates a Block object.
func NewBlock(header [32]byte, txnCount uint16, feeAmount *big.Int) Block {
	return Block{header, txnCount, feeAmount}
}
