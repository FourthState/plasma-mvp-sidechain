package store

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

// Wallet holds reference to the total balance, unspent, and spent outputs
// at a given address
type Wallet struct {
	Balance *big.Int          // total amount avaliable to be spent
	Unspent []plasma.Position // position of unspent transaction outputs
	Spent   []plasma.Position // position of spent transaction outputs
}

// Deposit wraps a plasma deposit with spend information.
type Deposit struct {
	Deposit   plasma.Deposit
	Spent     bool
	SpenderTx []byte // transaction hash that spends this deposit
}

// Output wraps a plasma output with spend information.
type Output struct {
	Output    plasma.Output
	Spent     bool
	InputTx   []byte // transaction hash that created this output
	SpenderTx []byte // transaction hash that spent this output
}

// Transaction wraps a plasma transaction with spend information.
type Transaction struct {
	InputTxs         [][]byte // transaction hashes the created the inputs of this transaction
	Transaction      plasma.Transaction
	ConfirmationHash []byte
	Spent            []bool
	SpenderTxs       [][]byte // transaction hashes that spend the outputs of this transaction
	Position         plasma.Position
}

// TxOutput holds all transactional information related to an output.
type TxOutput struct {
	plasma.Output
	Position       plasma.Position
	Spent          bool
	SpenderTx      []byte
	InputAddresses []ethcmn.Address
	InputPositions []plasma.Position
}

// NewTxOutput creates a TxOutput object.
func NewTxOutput(output plasma.Output, pos plasma.Position, spent bool, spenderTx []byte, inputAddresses []ethcmn.Address, inputPositions []plasma.Position) TxOutput {
	return TxOutput{
		Output:         output,
		Position:       pos,
		Spent:          spent,
		SpenderTx:      spenderTx,
		InputAddresses: inputAddresses,
		InputPositions: inputPositions,
	}
}

// Block wraps a plasma block with the tendermint block height.
type Block struct {
	plasma.Block
	TMBlockHeight uint64
}

type block struct {
	PlasmaBlock   plasma.Block
	TMBlockHeight uint64
}

// EncodeRLP RLP encodes a Block struct.
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &block{b.Block, b.TMBlockHeight})
}

// DecodeRLP decodes the byte stream into a Block.
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var block block
	if err := s.Decode(&block); err != nil {
		return err
	}

	b.Block = block.PlasmaBlock
	b.TMBlockHeight = block.TMBlockHeight
	return nil
}
