package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "plasma-mvp-sidechain/types"
	db "plasma-mvp-sidechain/db"
	utils "plasma-mvp-sidechain/utils"
	"reflect"
)

func NewHandler(uk db.UTXOKeeper, txIndex *uint16) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.SpendMsg:
			return handleSpendMsg(ctx, uk, msg, txIndex)
		default:
			errMsg := "Unrecognized Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle SpendMsg.
func handleSpendMsg(ctx sdk.Context, uk db.UTXOKeeper, msg types.SpendMsg, txIndex *uint16) sdk.Result {
	position1 := types.Position{msg.Blknum1, msg.Txindex1, msg.Oindex1, msg.DepositNum1}
	utxo1 := uk.UM.GetUTXO(ctx, position1)
	var position2 types.Position
	var utxo2 types.UTXO
	err := uk.SpendUTXO(ctx, msg.Owner1, position1)
	if err != nil {
		return err.Result()
	}

	if msg.Owner2 != nil && !utils.ZeroAddress(msg.Owner2) {
		position2 = types.Position{msg.Blknum2, msg.Txindex2, msg.Oindex2, msg.DepositNum2}
		utxo2 = uk.UM.GetUTXO(ctx, position2)
		err := uk.SpendUTXO(ctx, msg.Owner2, position2)
		if err != nil {
			return err.Result()
		}
	}

	oldUTXOs := [2]types.UTXO{utxo1, utxo2}
	err2 := uk.RecieveUTXO(ctx, msg.Newowner1, msg.Denom1, oldUTXOs, 0, *txIndex)
	if err2 != nil {
		return err2.Result()
	}
	if msg.Newowner2 != nil && !utils.ZeroAddress(msg.Newowner2) {
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
