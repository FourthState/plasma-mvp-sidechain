package plasma

import (
	"github.com/ethereum/go-ethereum/common"
)

// Input represents the input to a spend
type Input struct {
	Position
	Owner             common.Address `json:"Owner"`
	Signature         [65]byte       `json:"Signature"`
	ConfirmSignatures [][65]byte     `json:"ConfirmSignature"`
}

func NewInput(position Position, owner common.Address, sig [65]byte, confirmsigs [][65]byte) Input {
	return Input{
		Position:          position,
		Owner:             owner,
		Signature:         sig,
		ConfirmSignatures: confirmsigs,
	}
}
