package app

import (
	"github.com/tendermint/tendermint/crypto"
)

// GenesisState specifies the validator of the chain
type GenesisState struct {
	Validator GenesisValidator `json:"validator"`
}

// GenesisValidator holds the consensus public key and fee address of
// the validator. ConsPubKey is tendermint Ed25119 public key.
type GenesisValidator struct {
	ConsPubKey crypto.PubKey `json:"validator_pubkey"`
	Address    string        `json:"fee_address"`
}

// NewDefaultGenesisState returns a GenesisState instance
func NewDefaultGenesisState(pubKey crypto.PubKey) GenesisState {
	return GenesisState{
		Validator: GenesisValidator{pubKey, ""},
	}
}
