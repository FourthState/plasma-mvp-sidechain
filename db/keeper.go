package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	types "github.com/FourthState/plasma-mvp-sidechain/types"
)

type UTXOKeeper struct {
	UM types.UTXOMapper
}

// NewUTXOKeeper returns a new UTXOKeeper
func NewUTXOKeeper(um types.UTXOMapper) UTXOKeeper {
	return UTXOKeeper{UM: um}
}

// Delete's utxo from utxo store
func (uk UTXOKeeper) SpendUTXO(ctx sdk.Context, addr crypto.Address, position types.Position) sdk.Error {

	utxo := uk.UM.GetUTXO(ctx, position) // Get the utxo that should be spent
	// Check to see if utxo exists, will be taken care of in ante handler
	if utxo == nil {
		return types.ErrInvalidUTXO(types.DefaultCodespace, "Unrecognized UTXO. Does not exist.")
	}
	uk.UM.DeleteUTXO(ctx, position) // Delete utxo from utxo store
	return nil
}

// Creates a new utxo and adds it to the utxo store
func (uk UTXOKeeper) RecieveUTXO(ctx sdk.Context, addr crypto.Address, denom uint64,
	oldutxos [2]types.UTXO, oindex uint8, txIndex uint16) sdk.Error {
	var inputAddr1 crypto.Address
	var inputAddr2 crypto.Address
	if oldutxos[0] != nil {
		inputAddr1 = oldutxos[0].GetAddress()
	}

	if oldutxos[1] != nil {
		inputAddr2 = oldutxos[1].GetAddress()
	}

	inputAddresses := [2]crypto.Address{inputAddr1, inputAddr2}
	position := types.Position{uint64(ctx.BlockHeight()), txIndex, oindex, 0}
	utxo := types.NewBaseUTXO(addr, inputAddresses, denom, position)
	uk.UM.AddUTXO(ctx, utxo) // Adds utxo to utxo store
	return nil
}
