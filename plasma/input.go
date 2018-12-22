package plasma

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// Input represents the input to a spend
type Input struct {
	BlockNum          *big.Int       `json:"BlockNum"`
	TxIndex           *big.Int       `json:"TxIndex"`
	OutputIndex       *big.Int       `json:"OutputIndex"`
	DepositNonce      *big.Int       `json:"DepositNonce"`
	Owner             common.Address `json:"Owner"`
	ConfirmSignatures [][65]byte     `json:"ConfirmSignature"`
}

func newInput(blkNum, txIndex, oIndex, nonce []byte, owner common.Address, confirmsigs [][65]byte) *Input {
	return &Input{
		BlockNum:          new(big.Int).SetBytes(blkNum),
		TxIndex:           new(big.Int).SetBytes(txIndex),
		OutputIndex:       new(big.Int).SetBytes(oIndex),
		DepositNonce:      new(big.Int).SetBytes(nonce),
		Owner:             owner,
		ConfirmSignatures: confirmsigs,
	}
}
