package utxo

// UTXO is a standard unspent transaction output
type UTXO interface {
	// Address that owns UTXO
	GetAddress() []byte
	SetAddress([]byte) error // errors if already set

	GetAmount() uint64
	SetAmount(uint64) error // errors if already set

	GetDenom() uint64
	SetDenom(string) error // errors if already set

	GetPosition() Position
	SetPosition(Position) error // errors if already set
}

// Positions must be unqiue or a collision may result when using mapper.go
type Position interface {
	// utilized for mapping from position to utxo
	GetSignBytes() []byte
	// returns true if the position is valid, false otherwise
	IsValid() bool
}
