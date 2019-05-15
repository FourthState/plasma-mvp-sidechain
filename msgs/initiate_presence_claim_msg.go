package msgs

import (
	//"fmt"
	//"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	InitiatePresenceClaimMsgRoute = "initiatePresenceClaim"
)

type InitiatePresenceClaimMsg struct {
	ZoneID       [32]byte `json:"zoneID"`
	UTXOPosition [4]int64 `json:"uxoPosition"`
	Signature    []byte   `json:"signature"`
}

func (msg InitiatePresenceClaimMsg) Type() string { return "initiate_presence_claim" }

func (msg InitiatePresenceClaimMsg) Route() string { return InitiatePresenceClaimMsgRoute }

func (msg InitiatePresenceClaimMsg) TxHash() []byte {

	messageNoSig := InitiatePresenceClaimMsg{}

	messageNoSig.ZoneID = msg.ZoneID
	messageNoSig.UTXOPosition = msg.UTXOPosition

	bytes, _ := rlp.EncodeToBytes(&messageNoSig)

	return crypto.Keccak256(bytes)
}

// GetSigners will attempt to retrieve the signers of the message.
// CONTRACT: a nil slice is returned if recovery fails
func (msg InitiatePresenceClaimMsg) GetSigners() []sdk.AccAddress {
	txHash := utils.ToEthSignedMessageHash(msg.TxHash())

	pubKey, err := crypto.SigToPub(txHash, msg.Signature[:])
	if err != nil {
		return nil
	}

	return []sdk.AccAddress{sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes())}

}

func (msg InitiatePresenceClaimMsg) GetSignBytes() []byte {
	return msg.TxHash()
}

func (msg InitiatePresenceClaimMsg) ValidateBasic() sdk.Error {
	// TODO validate message here
	return nil
}

func (msg InitiatePresenceClaimMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}
