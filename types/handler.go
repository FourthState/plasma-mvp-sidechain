package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
)

func NewHandler(uk UTXOKeeper, txIndex *uint16) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SpendMsg:
			return handleSpendMsg(ctx, uk, msg, txIndex)
		default:
			errMsg := "Unrecognized Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle SpendMsg.
func handleSpendMsg(ctx sdk.Context, uk UTXOKeeper, msg SpendMsg, txIndex *uint16) sdk.Result {
	position1 := Position{msg.Blknum1, msg.Txindex1, msg.Oindex1}
	utxo1 := uk.um.GetUTXO(ctx, position1)
	var position2 Position
	var utxo2 UTXO
	err := uk.SpendUTXO(ctx, msg.Owner1, position1)
	if err != nil {
		return err.Result()
	}

	if msg.Owner2 != nil && !ZeroAddress(msg.Owner2) {
		position2 = Position{msg.Blknum2, msg.Txindex2, msg.Oindex2}
		utxo2 = uk.um.GetUTXO(ctx, position2)
		err := uk.SpendUTXO(ctx, msg.Owner2, position2)
		if err != nil {
			return err.Result()
		}
	}

	oldUTXOs := [2]UTXO{utxo1, utxo2}
	err2 := uk.RecieveUTXO(ctx, msg.Newowner1, msg.Denom1, oldUTXOs, 0, *txIndex)
	if err2 != nil {
		return err2.Result()
	}
	if msg.Newowner2 != nil && !ZeroAddress(msg.Newowner2) {
		err := uk.RecieveUTXO(ctx, msg.Newowner2, msg.Denom2, oldUTXOs, 1, *txIndex)
		if err != nil {
			return err.Result()
		}
	}
	
	// Increment txIndex
	if !ctx.IsCheckTx() {
		(*txIndex)++
	}
	// TODO: add some tags so we can search it!
	return sdk.Result{} // TODO
}