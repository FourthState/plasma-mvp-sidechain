package msgs

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
)

const (
	SpendMsgRoute = "spend"
)

type SpendMsg struct {
	plasma.Transaction
}

// satisfy rlp interface
func (msg *SpendMsg) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &msg.Transaction)
}

// satisfy rlp interface
func (msg *SpendMsg) DecodeRLP(s rlp.Stream) error {
	tx := plasma.Transaction{}
	if err := s.Decode(&tx); err != nil {
		return nil
	}

	msg.Transaction = tx

	return nil
}

// Implement the sdk.Msg interface

func (msg SpendMsg) Type() string { return "spend_utxo" }

func (msg SpendMsg) Route() string { return SpendMsgRoute }

func (msg SpendMsg) GetSigners() []sdk.AccAddress {
	txHash := msg.TxHash()
	addrs := make([]sdk.AccAddress, 1)

	// recover first owner
	pubKey, _ := crypto.SigToPub(txHash[:], msg.Input0.Signature[:])
	addrs = append(addrs, sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes()))

	if msg.HasSecondInput() {
		// recover the second owner
		pubKey, _ = crypto.SigToPub(txHash[:], msg.Input1.Signature[:])
		addrs = append(addrs, sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes()))
	}

	return addrs
}

func (msg SpendMsg) GetSignBytes() []byte {
	txHash := msg.TxHash()
	return txHash[:]
}

func (msg SpendMsg) ValidateBasic() sdk.Error {
	if err := msg.Transaction.ValidateBasic(); err != nil {
		return ErrInvalidTransaction(DefaultCodespace, err.Error())
	}

	return nil
}

// Also satisfy the sdk.Tx interface
func (msg SpendMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}
