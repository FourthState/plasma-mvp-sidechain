package utxo

import (
	types "github.com/FourthState/plasma-mvp-sidechain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

type UTXOKeeper struct {
	utxoMapper UTXOMapper
}

// NewUTXOKeeper returns a new UTXOKeeper
func NewUTXOKeeper(um types.UTXOMapper) UTXOKeeper {
	return UTXOKeeper{utxoMapper: um}
}

// Delete's utxo from utxo store
// AnteHandler should have already checked existence of the utxo
func (uk UTXOKeeper) SpendUTXO(ctx sdk.Context, addr []byte, position Position) {
	uk.utxoMapper.DeleteUTXO(ctx, addr, position)
}

// Creates a new utxo and adds it to the utxo store
func (uk UTXOKeeper) RecieveUTXO(ctx sdk.Context, utxo UTXO) {
	uk.utxoMapper.AddUTXO(ctx, utxo)
}

// Get UTXO from the mapper
func (uk UTXOKeeper) GetUTXO(ctx sdk.Context, addr []byte, position Position) UTXO {
	return uk.utxoMapper.GetUTXO(ctx, addr, position)
}
