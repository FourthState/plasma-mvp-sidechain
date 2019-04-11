package plasma

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"math/big"

	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// Transaction represents a spend of inputs. Fields should not be accessed directly
type Transaction struct {
	Input0  Input    `json:"input0"`
	Input1  Input    `json:"input1"`
	Output0 Output   `json:"output0"`
	Output1 Output   `json:"output1"`
	Fee     *big.Int `json:"fee"`
}

type txList struct {
	BlkNum0           [32]byte
	TxIndex0          [32]byte
	OIndex0           [32]byte
	DepositNonce0     [32]byte
	Input0ConfirmSigs [130]byte
	BlkNum1           [32]byte
	TxIndex1          [32]byte
	OIndex1           [32]byte
	DepositNonce1     [32]byte
	Input1ConfirmSigs [130]byte
	NewOwner0         common.Address
	Amount0           [32]byte
	NewOwner1         common.Address
	Amount1           [32]byte
	Fee               [32]byte
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

	confirmSigs0 := parseSig(t.Tx.Input0ConfirmSigs)
	confirmSigs1 := parseSig(t.Tx.Input1ConfirmSigs)

	tx.Input0 = NewInput(NewPosition(big.NewInt(new(big.Int).SetBytes(t.Tx.BlkNum0[:]).Int64()), uint16(new(big.Int).SetBytes(t.Tx.TxIndex0[:]).Int64()), uint8(new(big.Int).SetBytes(t.Tx.OIndex0[:]).Int64()), big.NewInt(new(big.Int).SetBytes(t.Tx.DepositNonce0[:]).Int64())),
		t.Sigs[0], confirmSigs0)
	tx.Input1 = NewInput(NewPosition(big.NewInt(new(big.Int).SetBytes(t.Tx.BlkNum1[:]).Int64()), uint16(new(big.Int).SetBytes(t.Tx.TxIndex1[:]).Int64()), uint8(new(big.Int).SetBytes(t.Tx.OIndex1[:]).Int64()), big.NewInt(new(big.Int).SetBytes(t.Tx.DepositNonce1[:]).Int64())),
		t.Sigs[1], confirmSigs1)
	// set signatures if applicable
	tx.Output0 = NewOutput(t.Tx.NewOwner0, big.NewInt(new(big.Int).SetBytes(t.Tx.Amount0[:]).Int64()))
	tx.Output1 = NewOutput(t.Tx.NewOwner1, big.NewInt(new(big.Int).SetBytes(t.Tx.Amount1[:]).Int64()))
	tx.Fee = big.NewInt(new(big.Int).SetBytes(t.Tx.Fee[:]).Int64())

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

	return nil
}

func (tx Transaction) TxBytes() []byte {
	bytes, _ := rlp.EncodeToBytes(&tx)
	return bytes
}

// TxHash returns the bytes the signatures are signed over
func (tx Transaction) TxHash() []byte {
	txList := tx.toTxList()
	bytes, _ := rlp.EncodeToBytes(&txList)
	return crypto.Keccak256(bytes)
}

// MerkleHash returns the bytes that is included in the merkle tree
func (tx Transaction) MerkleHash() []byte {
	hash := sha256.Sum256(tx.TxBytes())
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

	// pointer safety if a transaction
	// object was ever created with Transaction{}
	txList := txList{}
	if tx.Input0.BlockNum == nil {
		tx.Input0.BlockNum = utils.Big0
	} else if tx.Input1.BlockNum == nil {
		tx.Input1.BlockNum = utils.Big0
	}

	if tx.Input0.DepositNonce == nil {
		tx.Input0.DepositNonce = utils.Big0
	} else if tx.Input1.DepositNonce == nil {
		tx.Input1.DepositNonce = utils.Big0
	}

	if tx.Output0.Amount == nil {
		tx.Output0.Amount = utils.Big0
	} else if tx.Output1.Amount == nil {
		tx.Output1.Amount = utils.Big0
	}

	if tx.Fee == nil {
		tx.Fee = utils.Big0
	}

	// fill in txList with values
	// Input 0
	if len(tx.Input0.BlockNum.Bytes()) > 0 {
		copy(txList.BlkNum0[32-len(tx.Input0.BlockNum.Bytes()):], tx.Input0.BlockNum.Bytes())
	}
	txList.TxIndex0[31] = byte(tx.Input0.TxIndex)
	txList.TxIndex0[30] = byte(tx.Input0.TxIndex >> 8)
	txList.OIndex0[31] = byte(tx.Input0.OutputIndex)
	if len(tx.Input0.DepositNonce.Bytes()) > 0 {
		copy(txList.DepositNonce0[32-len(tx.Input0.DepositNonce.Bytes()):], tx.Input0.DepositNonce.Bytes())
	}
	switch len(tx.Input0.ConfirmSignatures) {
	case 1:
		copy(txList.Input0ConfirmSigs[:65], tx.Input0.ConfirmSignatures[0][:])
	case 2:
		copy(txList.Input0ConfirmSigs[:65], tx.Input0.ConfirmSignatures[0][:])
		copy(txList.Input0ConfirmSigs[65:], tx.Input0.ConfirmSignatures[1][:])
	}

	// Input 1
	if len(tx.Input1.BlockNum.Bytes()) > 0 {
		copy(txList.BlkNum1[32-len(tx.Input1.BlockNum.Bytes()):], tx.Input1.BlockNum.Bytes())
	}
	txList.TxIndex1[31] = byte(tx.Input1.TxIndex)
	txList.TxIndex1[30] = byte(tx.Input1.TxIndex >> 8)
	txList.OIndex1[31] = byte(tx.Input1.OutputIndex)
	if len(tx.Input1.DepositNonce.Bytes()) > 0 {
		copy(txList.DepositNonce1[32-len(tx.Input1.DepositNonce.Bytes()):], tx.Input1.DepositNonce.Bytes())
	}

	switch len(tx.Input1.ConfirmSignatures) {
	case 1:
		copy(txList.Input1ConfirmSigs[:65], tx.Input1.ConfirmSignatures[0][:])
	case 2:
		copy(txList.Input1ConfirmSigs[:65], tx.Input1.ConfirmSignatures[0][:])
		copy(txList.Input1ConfirmSigs[65:], tx.Input1.ConfirmSignatures[1][:])
	}

	// Outputs and Fee
	txList.NewOwner0 = tx.Output0.Owner
	if len(tx.Output0.Amount.Bytes()) > 0 {
		copy(txList.Amount0[32-len(tx.Output0.Amount.Bytes()):], tx.Output0.Amount.Bytes())
	}
	txList.NewOwner1 = tx.Output1.Owner
	if len(tx.Output1.Amount.Bytes()) > 0 {
		copy(txList.Amount1[32-len(tx.Output1.Amount.Bytes()):], tx.Output1.Amount.Bytes())
	}
	if len(tx.Fee.Bytes()) > 0 {
		copy(txList.Fee[32-len(tx.Fee.Bytes()):], tx.Fee.Bytes())
	}
	return txList
}

// Helpers
// Convert 130 byte input confirm sigs to 65 byte slices
func parseSig(sig [130]byte) [][65]byte {
	if bytes.Equal(sig[:65], make([]byte, 65)) {
		return [][65]byte{}
	} else if bytes.Equal(sig[65:], make([]byte, 65)) {
		newSig := make([][65]byte, 1)
		copy(newSig[0][:], sig[:65])
		return newSig
	} else {
		newSig := make([][65]byte, 2)
		copy(newSig[0][:], sig[:65])
		copy(newSig[1][:], sig[65:])
		return newSig
	}
}
