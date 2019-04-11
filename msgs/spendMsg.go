package msgs

import (
	//"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	SpendMsgRoute = "spend"
)

// SpendMsg implements the RLP interface through `Transaction`
type SpendMsg struct {
	plasma.Transaction
}

// Implement the sdk.Msg interface

func (msg SpendMsg) Type() string { return "spend_utxo" }

func (msg SpendMsg) Route() string { return SpendMsgRoute }

// GetSigners will attempt to retrieve the signers of the message.
// CONTRACT: a nil slice is returned if recovery fails
func (msg SpendMsg) GetSigners() []sdk.AccAddress {
	txHash := utils.ToEthSignedMessageHash(msg.TxHash())
	//fmt.Println("GET SIGNERS TX HASH", ethcmn.ToHex(txHash), "INPUT0 SIGNATURE ", ethcmn.ToHex(msg.Input0.Signature[:]))
	var addrs []sdk.AccAddress

	// recover first owner
	pubKey, err := crypto.SigToPub(txHash, msg.Input0.Signature[:])
	if err != nil {
		return nil
	}
	addrs = append(addrs, sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes()))

	if msg.HasSecondInput() {
		// recover the second owner
		pubKey, err = crypto.SigToPub(txHash, msg.Input1.Signature[:])
		if err != nil {
			return nil
		}
		addrs = append(addrs, sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes()))
	}

	return addrs
}

func (msg SpendMsg) GetSignBytes() []byte {
	return msg.TxHash()
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
