package plasma

import (
	"math/big"
)

// Block represents a plasma block.
type Block struct {
	Header    [32]byte
	TxnCount  uint16
	FeeAmount *big.Int
	Height    *big.Int
}

// NewBlock creates a Block object.
func NewBlock(header [32]byte, txnCount uint16, feeAmount, height *big.Int) Block {
	return Block{
		Header:    header,
		TxnCount:  txnCount,
		FeeAmount: feeAmount,
		Height:    height,
	}
}
