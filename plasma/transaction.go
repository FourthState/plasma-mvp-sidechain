package plasma

import (
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

// Transaction represents a spend of inputs. Fields should not be accessed directly
type Transaction struct {
	Input0  *Input   `json:"Input0"`
	Sig0    [65]byte `json:"Sig0"`
	Input1  *Input   `json:"Input1"`
	Sig1    [65]byte `json:"Sig1"`
	Output0 *Output  `json:"Output0"`
	Output1 *Output  `json:"Output1"`
	Fee     *big.Int `json:"Fee"`
}

type txList struct {
	BlkNum0           []byte
	TxIndex0          []byte
	OIndex0           []byte
	DepositNonce0     []byte
	Owner0            common.Address
	Input0ConfirmSigs [][65]byte
	BlkNum1           []byte
	TxIndex1          []byte
	OIndex1           []byte
	DepositNonce1     []byte
	Owner1            common.Address
	Input1ConfirmSigs [][65]byte
	NewOwner0         common.Address
	Amount0           []byte
	NewOwner1         common.Address
	Amount1           []byte
	Fee               []byte
}

type rawTx struct {
	Tx   txList
	Sigs [2][65]byte
}

// EncodeRLP satisfies the rlp interface for Transaction
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	t := &rawTx{*tx.toTxList(), [2][65]byte{tx.Sig0, tx.Sig1}}

	return rlp.Encode(w, t)
}

// DecodeRLP satisfies the rlp interface for Transaction
func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	var t rawTx
	if err := s.Decode(&t); err != nil {
		return err
	}

	tx.Input0 = newInput(t.Tx.BlkNum0, t.Tx.TxIndex0, t.Tx.OIndex0, t.Tx.DepositNonce0, t.Tx.Owner0, t.Tx.Input0ConfirmSigs)
	tx.Input1 = newInput(t.Tx.BlkNum1, t.Tx.TxIndex1, t.Tx.OIndex1, t.Tx.DepositNonce1, t.Tx.Owner1, t.Tx.Input1ConfirmSigs)
	tx.Output0 = newOutput(t.Tx.NewOwner0, t.Tx.Amount0)
	tx.Output1 = newOutput(t.Tx.NewOwner1, t.Tx.Amount1)
	tx.Sig0 = t.Sigs[0]
	tx.Sig1 = t.Sigs[1]
	tx.Fee = new(big.Int).SetBytes(t.Tx.Fee)

	return nil
}

// TxHash returns the bytes the signatures are signed over
func (tx *Transaction) TxHash() [32]byte {
	txList := tx.toTxList()
	bytes, _ := rlp.EncodeToBytes(txList)
	bytes = crypto.Keccak256(bytes)

	var result [32]byte
	copy(result[:], bytes)
	return result
}

// MerkleHash returns the bytes that is included in the merkle tree
func (tx *Transaction) MerkleHash() [32]byte {
	bytes, _ := rlp.EncodeToBytes(tx)

	return sha256.Sum256(bytes)
}

// IndexAt returns the input specified by i, 0/1
func (tx *Transaction) IndexAt(i uint8) *Input {
	if i == 0 {
		return tx.Input0
	}

	return tx.Input1
}

// OutputAt returns the output specified by i, 0/1
func (tx *Transaction) OutputAt(i uint8) *Output {
	if i == 0 {
		return tx.Output0
	}

	return tx.Output1
}

// SigAt returns the transaction signature specified by i, 0/1
func (tx *Transaction) SigAt(i uint8) [65]byte {
	if i == 0 {
		return tx.Sig0
	}

	return tx.Sig1
}

/* Helpers */

func (tx *Transaction) toTxList() *txList {
	return &txList{
		BlkNum0:           tx.Input0.BlockNum.Bytes(),
		TxIndex0:          tx.Input0.TxIndex.Bytes(),
		OIndex0:           tx.Input0.OutputIndex.Bytes(),
		DepositNonce0:     tx.Input0.DepositNonce.Bytes(),
		Owner0:            tx.Input0.Owner,
		Input0ConfirmSigs: tx.Input0.ConfirmSignatures,
		BlkNum1:           tx.Input1.BlockNum.Bytes(),
		TxIndex1:          tx.Input1.TxIndex.Bytes(),
		OIndex1:           tx.Input1.OutputIndex.Bytes(),
		DepositNonce1:     tx.Input1.DepositNonce.Bytes(),
		Owner1:            tx.Input1.Owner,
		Input1ConfirmSigs: tx.Input1.ConfirmSignatures,
		NewOwner0:         tx.Output0.Owner,
		Amount0:           tx.Output0.Amount.Bytes(),
		NewOwner1:         tx.Output1.Owner,
		Amount1:           tx.Output1.Amount.Bytes(),
		Fee:               tx.Fee.Bytes(),
	}
}
