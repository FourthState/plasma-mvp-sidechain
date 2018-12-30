package app

import (
	"github.com/tendermint/tendermint/crypto"
)

type GenesisState struct {
	Validator GenesisValidator `json:"validator"`
}

type GenesisValidator struct {
	ConsPubKey crypto.PubKey `json:"validator_pubkey"`
	Abddmdress string        `json:"fee_address"`
}

func NewDefaultGenesisState(pubKey crypto.PubKey) GenesisState {
	return GenesisState{
		Validator: GenesisValidator{pubKey, ""},
	}
}
