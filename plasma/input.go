package plasma

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// Input represents the input to a spend
type Input struct {
	Position
	Owner             common.Address `json:"Owner"`
	Signature         [65]byte       `json:"Signature"`
	ConfirmSignatures [][65]byte     `json:"ConfirmSignature"`
}

func NewInput(blkNum *big.Int, txIndex uint16, oIndex uint8, nonce *big.Int, owner common.Address, sig [65]byte, confirmsigs [][65]byte) Input {
	return Input{
		Position:          NewPosition(blkNum, txIndex, oIndex, nonce),
		Owner:             owner,
		Signature:         sig,
		ConfirmSignatures: confirmsigs,
	}
}
