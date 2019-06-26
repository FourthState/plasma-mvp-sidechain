package plasma

import (
	"github.com/ethereum/go-ethereum/rlp"
	"io"
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

// EncodeRLP satisfies the rlp interface for Block.
func (b *Block) EncodeRLP(w io.Writer) error {
	blk := &block{b.Header, b.TxnCount, b.FeeAmount.Bytes()}

	return rlp.Encode(w, blk)
}

// DecodeRLP satisfies the rlp interface for Block.
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var blk block
	if err := s.Decode(&blk); err != nil {
		return err
	}

	b.Header = blk.Header
	b.TxnCount = blk.TxnCount
	b.FeeAmount = new(big.Int).SetBytes(blk.FeeAmount)

	return nil
}
