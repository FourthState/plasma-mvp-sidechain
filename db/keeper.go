package db

import (
	types "github.com/FourthState/plasma-mvp-sidechain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

type UTXOKeeper struct {
	UM types.UTXOMapper
}

// NewUTXOKeeper returns a new UTXOKeeper
func NewUTXOKeeper(um types.UTXOMapper) UTXOKeeper {
	return UTXOKeeper{UM: um}
}

// Delete's utxo from utxo store
// AnteHandler should have already checked existence of the utxo
func (uk UTXOKeeper) SpendUTXO(ctx sdk.Context, addr common.Address, position types.Position) {
	uk.UM.DeleteUTXO(ctx, addr, position)
}

// Creates a new utxo and adds it to the utxo store
func (uk UTXOKeeper) RecieveUTXO(ctx sdk.Context, addr common.Address, denom uint64,
	oldutxos [2]types.UTXO, oindex uint8, txIndex uint16) {

	inputAddr1 := oldutxos[0].GetAddress()
	var inputAddr2 common.Address

	// oldutxo[1] may be nil
	if oldutxos[1] != nil {
		inputAddr2 = oldutxos[1].GetAddress()
	}

	inputAddresses := [2]common.Address{inputAddr1, inputAddr2}
	position := types.Position{uint64(ctx.BlockHeight()), txIndex, oindex, 0}
	utxo := types.NewBaseUTXO(addr, inputAddresses, denom, position)
	uk.UM.AddUTXO(ctx, utxo) // Adds utxo to utxo store
}
