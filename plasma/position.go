package plasma

import (
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

// Position for an input/output
type Position struct {
	BlockNum     *big.Int `json:"BlockNum"`
	TxIndex      uint16   `json:"TxIndex"`
	OutputIndex  uint8    `json:"OutputIndex"`
	DepositNonce *big.Int `json:"DepositNonce"`
}

type position struct {
	BlockNum     []byte
	TxIndex      uint16
	OutputIndex  uint8
	DepositNonce []byte
}

const (
	blockIndexFactor = 1000000
	txIndexFactor    = 10
)

func NewPosition(blkNum *big.Int, txIndex uint16, oIndex uint8, depositNonce *big.Int) Position {
	return Position{
		BlockNum:     blkNum,
		TxIndex:      txIndex,
		OutputIndex:  oIndex,
		DepositNonce: depositNonce,
	}
}

func (p *Position) EncodeRLP(w io.Writer) error {
	pos := position{p.BlockNum.Bytes(), p.TxIndex, p.OutputIndex, p.DepositNonce.Bytes()}

	return rlp.Encode(w, pos)
}

func (p *Position) DecodeRLP(s *rlp.Stream) error {
	var pos position
	if err := s.Decode(&pos); err != nil {
		return err
	}

	p.BlockNum = new(big.Int).SetBytes(pos.BlockNum)
	p.TxIndex = pos.TxIndex
	p.OutputIndex = pos.OutputIndex
	p.DepositNonce = new(big.Int).SetBytes(pos.DepositNonce)

	return nil
}

func (p Position) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(&p)
	return bytes
}

func (p Position) IsDeposit() bool {
	return p.DepositNonce.Sign() != 0
}

func (p Position) Priority() *big.Int {
	if p.IsDeposit() {
		return p.DepositNonce
	}

	bFactor := big.NewInt(blockIndexFactor)
	tFactor := big.NewInt(txIndexFactor)

	bFactor = bFactor.Mul(bFactor, p.BlockNum)
	tFactor = tFactor.Mul(tFactor, big.NewInt(int64(p.TxIndex)))

	temp := new(big.Int).Add(bFactor, tFactor)
	return temp.Add(temp, big.NewInt(int64(p.OutputIndex)))
}
