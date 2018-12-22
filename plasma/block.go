package plasma

import (
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

// Block represent one plasma block
type Block struct {
	Header   [32]byte `json:"header"`
	TxnCount uint64   `json:"txnCount"`
	TotalFee *big.Int `json:"totalFee"`
}

type block struct {
	Header   [32]byte
	TxnCount uint64
	Fee      []byte
}

// NewBlock is a constructor for Block
func NewBlock(header [32]byte, txnCount uint64, totalFee *big.Int) *Block {
	return &Block{header, txnCount, totalFee}
}

// EncodeRLP satisfies the rlp interface for Block
func (b *Block) EncodeRLP(w io.Writer) error {
	blk := &block{b.Header, b.TxnCount, b.TotalFee.Bytes()}

	return rlp.Encode(w, blk)
}

// DecodeRLP satisfies the rlp interface for Block
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var blk block
	if err := s.Decode(&blk); err != nil {
		return err
	}

	b.Header = blk.Header
	b.TxnCount = blk.TxnCount
	b.TotalFee = new(big.Int).SetBytes(blk.Fee)

	return nil
}
