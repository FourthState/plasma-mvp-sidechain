package plasma

import (
	"crypto/sha256"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

// Transaction represents a spend of inputs. Fields should not be accessed directly
type Transaction struct {
	Input0  Input
	Input1  Input
	Output0 Output
	Output1 Output
	Fee     *big.Int
}

type txList struct {
	BlkNum0           []byte
	TxIndex0          uint16
	OIndex0           uint8
	DepositNonce0     []byte
	Input0ConfirmSigs [][65]byte
	BlkNum1           []byte
	TxIndex1          uint16
	OIndex1           uint8
	DepositNonce1     []byte
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
	t := &rawTx{tx.toTxList(), [2][65]byte{tx.Input0.Signature, tx.Input1.Signature}}

	return rlp.Encode(w, t)
}

// DecodeRLP satisfies the rlp interface for Transaction
func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	var t rawTx
	if err := s.Decode(&t); err != nil {
		return err
	}

	tx.Input0 = NewInput(NewPosition(new(big.Int).SetBytes(t.Tx.BlkNum0), t.Tx.TxIndex0, t.Tx.OIndex0, new(big.Int).SetBytes(t.Tx.DepositNonce0)),
		t.Sigs[0], t.Tx.Input0ConfirmSigs)
	tx.Input1 = NewInput(NewPosition(new(big.Int).SetBytes(t.Tx.BlkNum1), t.Tx.TxIndex1, t.Tx.OIndex1, new(big.Int).SetBytes(t.Tx.DepositNonce1)),
		t.Sigs[1], t.Tx.Input1ConfirmSigs)
	// set signatures if applicable
	tx.Output0 = NewOutput(t.Tx.NewOwner0, new(big.Int).SetBytes(t.Tx.Amount0))
	tx.Output1 = NewOutput(t.Tx.NewOwner1, new(big.Int).SetBytes(t.Tx.Amount1))
	tx.Fee = new(big.Int).SetBytes(t.Tx.Fee)

	return nil
}

func (tx Transaction) ValidateBasic() error {
	// validate inputs
	if err := tx.Input0.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid first input { %s }", err)
	}
	if tx.Input0.Position.IsNilPosition() {
		return fmt.Errorf("first input cannot be nil")
	}
	if err := tx.Input1.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid second input { %s }", err)
	}
	if tx.Input0.Position.String() == tx.Input1.Position.String() {
		return fmt.Errorf("same position cannot be spent twice")
	}

	// validate outputs
	if err := tx.Output0.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid first output { %s }", err)
	}
	if utils.IsZeroAddress(tx.Output0.Owner) || tx.Output0.Amount.Sign() == 0 {
		return fmt.Errorf("first output must have a valid address and non-zero amount")
	}
	if err := tx.Output1.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid second output { %s }", err)
	}

	// validate signatures
	txHash := tx.TxHash()
	_, err := crypto.SigToPub(txHash, tx.Input0.Signature[:])
	if err != nil {
		return fmt.Errorf("first input has an invalid signature { %s }", err)
	}
	if tx.HasSecondInput() {
		_, err := crypto.SigToPub(txHash, tx.Input1.Signature[:])
		if err != nil {
			return fmt.Errorf("second input has an invalid signature { %s }", err)
		}
	}

	return nil
}

// TxHash returns the bytes the signatures are signed over
func (tx Transaction) TxHash() []byte {
	txList := tx.toTxList()
	bytes, _ := rlp.EncodeToBytes(txList)
	return crypto.Keccak256(bytes)
}

// MerkleHash returns the bytes that is included in the merkle tree
func (tx Transaction) MerkleHash() []byte {
	bytes, _ := rlp.EncodeToBytes(tx)
	hash := sha256.Sum256(bytes)

	return hash[:]
}

// HasSecondInput is an indicator for the existence of a second input
func (tx Transaction) HasSecondInput() bool {
	return !tx.Input1.Position.IsNilPosition()
}

// HasSecondOutput is an indicator if the second output is used in this transaction
func (tx Transaction) HasSecondOutput() bool {
	return !utils.IsZeroAddress(tx.Output1.Owner)
}

// OutputAt is a selector for the outputs
func (tx Transaction) OutputAt(i uint8) Output {
	if i == 0 {
		return tx.Output0
	}

	return tx.Output1
}

// InputAt is a selector for the inputs
func (tx Transaction) InputAt(i uint8) Input {
	if i == 0 {
		return tx.Input0
	}

	return tx.Input1
}

func (tx Transaction) String() string {
	return fmt.Sprintf("Input0: %s\nInput1: %s\nOutput0: %s\nOutput1: %s\nFee: %s\n",
		tx.Input0, tx.Input1, tx.Output0, tx.Output1, tx.Fee)
}

/* Helpers */

func (tx Transaction) toTxList() txList {
	return txList{
		BlkNum0:           tx.Input0.BlockNum.Bytes(),
		TxIndex0:          tx.Input0.TxIndex,
		OIndex0:           tx.Input0.OutputIndex,
		DepositNonce0:     tx.Input0.DepositNonce.Bytes(),
		Input0ConfirmSigs: tx.Input0.ConfirmSignatures,
		BlkNum1:           tx.Input1.BlockNum.Bytes(),
		TxIndex1:          tx.Input1.TxIndex,
		OIndex1:           tx.Input1.OutputIndex,
		DepositNonce1:     tx.Input1.DepositNonce.Bytes(),
		Input1ConfirmSigs: tx.Input1.ConfirmSignatures,
		NewOwner0:         tx.Output0.Owner,
		Amount0:           tx.Output0.Amount.Bytes(),
		NewOwner1:         tx.Output1.Owner,
		Amount1:           tx.Output1.Amount.Bytes(),
		Fee:               tx.Fee.Bytes(),
	}
}
