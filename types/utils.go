package types

import (
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	"github.com/ethereum/go-ethereum/common"
)

// helper function for creating a utxo with the correct msg hash
func NewBaseUTXOWithMsgHash(addr common.Address, inputaddr [2]common.Address, amount uint64,
	denom string, position PlasmaPosition, msghash []byte) utxo.UTXO {
	return &BaseUTXO{
		MsgHash:        msghash,
		InputAddresses: inputaddr,
		Address:        addr,
		Amount:         amount,
		Denom:          denom,
		Position:       position,
	}
}
