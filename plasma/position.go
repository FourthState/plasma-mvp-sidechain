package plasma

import (
	"bytes"
	"math/big"
)

// Position for an input/output
type Position struct {
	BlockNum     *big.Int `json:"BlockNum"`
	TxIndex      uint16   `json:"TxIndex"`
	OutputIndex  uint8    `json:"OutputIndex"`
	DepositNonce *big.Int `json:"DepositNonce"`
}

func NewPosition(blkNum *big.Int, txIndex uint16, oIndex uint8, depositNonce *big.Int) Position {
	return Position{
		BlockNum:     blkNum,
		TxIndex:      txIndex,
		OutputIndex:  oIndex,
		DepositNonce: depositNonce,
	}
}

func (p Position) Bytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(p.BlockNum.Bytes())
	buffer.WriteByte(byte(p.TxIndex))
	buffer.WriteByte(byte(p.OutputIndex))
	buffer.Write(p.DepositNonce.Bytes())

	return buffer.Bytes()
}

func (p Position) IsDeposit() bool {
	return p.DepositNonce.Sign() == 0
}
