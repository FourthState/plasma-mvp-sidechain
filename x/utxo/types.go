package utxo

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UTXO is a standard unspent transaction output
type UTXO interface {
	// Address that owns UTXO
	GetAddress() []byte
	SetAddress([]byte) error // errors if already set

	GetAmount() uint64
	SetAmount(uint64) error // errors if already set

	GetDenom() string
	SetDenom(string) error // errors if already set

	GetPosition() Position
	SetPosition(Position) error // errors if already set
}

// Positions must be unqiue or a collision may result when using mapper.go
type Position interface {
	// Position is a uint slice
	Get() []uint64      // get position int slice. Return nil if unset.
	Set([]uint64) error // errors if already set
	// returns true if the position is valid, false otherwise
	IsValid() bool
}

// SpendMsg is an interface that wraps sdk.Msg with additional information
// for the UTXO spend handler.
type SpendMsg interface {
	sdk.Msg

	Inputs() []Input
	Outputs() []Output
	Fee() []Output // Owner is nil
}

type Input struct {
	Owner []byte
	Position
}

type Output struct {
	Owner  []byte
	Denom  string
	Amount uint64
}
