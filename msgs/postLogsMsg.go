package msgs

import (
	//"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	PostLogsMsgRoute = "postLogs"
)

type PostLogsMsg struct {
	ClaimID   []byte           `json:"claimID"`
	LogsHash  []byte           `json:"logsHash"`
	Beacons   []ethcmn.Address `json:"beacons"`
	Signature []byte           `json:"signature"`
}

func (msg PostLogsMsg) Type() string  { return "post_logs" }
func (msg PostLogsMsg) Route() string { return PostLogsMsgRoute }

func (msg PostLogsMsg) TxHash() []byte {

	postLogsMsg := PostLogsMsg{}

	postLogsMsg.ClaimID = msg.ClaimID
	postLogsMsg.LogsHash = msg.LogsHash
	postLogsMsg.Beacons = msg.Beacons

	bytes, _ := rlp.EncodeToBytes(&postLogsMsg)

	return crypto.Keccak256(bytes)
}

// GetSigners will attempt to retrieve the signers of the message.
// CONTRACT: a nil slice is returned if recovery fails

func (msg PostLogsMsg) GetSigners() []sdk.AccAddress {
	txHash := utils.ToEthSignedMessageHash(msg.TxHash())

	pubKey, err := crypto.SigToPub(txHash, msg.Signature[:])
	if err != nil {
		return nil
	}

	return []sdk.AccAddress{sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes())}

}

func (msg PostLogsMsg) GetSignBytes() []byte {
	return msg.TxHash()
}

func (msg PostLogsMsg) ValidateBasic() sdk.Error {
	// TODO validate message here
	return nil
}

func (msg PostLogsMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}
