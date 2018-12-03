package utxo

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
)

// UTXO is a standard unspent transaction output
// When spent, it becomes invalid and spender keys are filled in
type UTXO struct {
	InputKeys   [][]byte // Keys in store for input UTXOs that created this output
	Address     []byte
	Amount      uint64
	Denom       string
	Valid       bool
	Position    Position
	TxHash      []byte   // transaction that created this UTXO
	SpenderKeys [][]byte // Keys in store for UTXOs that spent this output
}

func NewUTXO(owner []byte, amount uint64, denom string, position Position) UTXO {
	return UTXO{
		Address:  owner,
		Amount:   amount,
		Denom:    denom,
		Position: position,
		Valid:    true,
	}
}

func NewUTXOwithInputs(owner []byte, amount uint64, denom string, position Position, txHash []byte, inputKeys [][]byte) UTXO {
	return UTXO{
		InputKeys: inputKeys,
		Address:   owner,
		Amount:    amount,
		Denom:     denom,
		Position:  position,
		Valid:     true,
		TxHash:    txHash,
	}
}

func (utxo UTXO) StoreKey(cdc *amino.Codec) []byte {
	encPos := cdc.MustMarshalBinaryBare(utxo.Position)
	return append(utxo.Address, encPos...)
}

// Recovers InputAddresses from Input keys.
// Assumes all addresses on-chain have same length
func (utxo UTXO) InputAddresses() [][]byte {
	var addresses [][]byte
	addrLen := len(utxo.Address)
	for _, key := range utxo.InputKeys {
		addresses = append(addresses, key[:addrLen])
	}
	return addresses
}

// Recovers InputPositions from Input keys.
// Assumes all addresses on-chain have same length
func (utxo UTXO) InputPositions(cdc *amino.Codec, proto ProtoPosition) []Position {
	var inputs []Position
	addrLen := len(utxo.Address)
	for _, key := range utxo.InputKeys {
		encodedPos := key[addrLen:]
		pos := proto()
		cdc.MustUnmarshalBinaryBare(encodedPos, pos)
		inputs = append(inputs, pos)
	}
	return inputs
}

// Recovers Spender Addresses from Spender keys.
// Assumes all addresses on-chain have same length
func (utxo UTXO) SpenderAddresses() [][]byte {
	var addresses [][]byte
	addrLen := len(utxo.Address)
	for _, key := range utxo.SpenderKeys {
		addresses = append(addresses, key[:addrLen])
	}
	return addresses
}

// Recovers Spender Positions from Spender keys.
// Assumes all addresses on-chain have same length
func (utxo UTXO) SpenderPositions(cdc *amino.Codec, proto ProtoPosition) []Position {
	var spenders []Position
	addrLen := len(utxo.Address)
	for _, key := range utxo.SpenderKeys {
		encodedPos := key[addrLen:]
		pos := proto()
		cdc.MustUnmarshalBinaryBare(encodedPos, pos)
		spenders = append(spenders, pos)
	}
	return spenders
}

// Create a prototype Position.
// Must return pointer to struct implementing Position
type ProtoPosition func() Position

// Positions must be unqiue or a collision may result when using mapper.go
type Position interface {
	// Position is a uint slice
	Get() []sdk.Uint // get position int slice. Return nil if unset.

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
