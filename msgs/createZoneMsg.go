package msgs

import (
	//"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	CreateZoneMsgRoute = "createZone"
)

type CreateZoneMsg struct {
	ZoneID    []byte           `json:"zoneID`
	Beacons   []common.Address `json:"beacons"`
	Geohash   string           `json:"geohash"`
	Signature []byte           `json:"signature"`
}

func (msg CreateZoneMsg) Type() string  { return "create_zone" }
func (msg CreateZoneMsg) Route() string { return CreateZoneMsgRoute }

func (msg CreateZoneMsg) TxHash() []byte {

	createZoneMsg := CreateZoneMsg{}
	createZoneMsg.ZoneID = msg.ZoneID
	createZoneMsg.Beacons = msg.Beacons
	createZoneMsg.Geohash = msg.Geohash

	bytes, _ := rlp.EncodeToBytes(&createZoneMsg)

	return crypto.Keccak256(bytes)
}

// GetSigners will attempt to retrieve the signers of the message.
// CONTRACT: a nil slice is returned if recovery fails

func (msg CreateZoneMsg) GetSigners() []sdk.AccAddress {
	txHash := utils.ToEthSignedMessageHash(msg.TxHash())

	pubKey, err := crypto.SigToPub(txHash, msg.Signature[:])
	if err != nil {
		return nil
	}

	return []sdk.AccAddress{sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes())}

}

func (msg CreateZoneMsg) GetSignBytes() []byte {
	return msg.TxHash()
}

func (msg CreateZoneMsg) ValidateBasic() sdk.Error {
	// TODO validate message here
	return nil
}

func (msg CreateZoneMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}
