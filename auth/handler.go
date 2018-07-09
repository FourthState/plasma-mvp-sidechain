package auth

import (
	db "github.com/FourthState/plasma-mvp-sidechain/db"
	types "github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
// Spends inputs, creates new outputs
func handleSpendMsg(ctx sdk.Context, uk db.UTXOKeeper, msg types.SpendMsg, txIndex *uint16) sdk.Result {

	position1 := types.Position{msg.Blknum1, msg.Txindex1, msg.Oindex1, msg.DepositNum1}
	inputAddr1 := msg.Owner1
	utxo1 := uk.UM.GetUTXO(ctx, inputAddr1, position1)
	uk.SpendUTXO(ctx, inputAddr1, position1)

	var position2 types.Position
	var utxo2 types.UTXO
	inputAddr2 := msg.Owner2
	if !utils.ZeroAddress(inputAddr2) {
		position2 = types.Position{msg.Blknum2, msg.Txindex2, msg.Oindex2, msg.DepositNum2}
		utxo2 = uk.UM.GetUTXO(ctx, inputAddr2, position2)
		uk.SpendUTXO(ctx, inputAddr2, position2)
	}

	oldUTXOs := [2]types.UTXO{utxo1, utxo2}
	uk.RecieveUTXO(ctx, msg.Newowner1, msg.Denom1, oldUTXOs, 0, *txIndex)
	if !utils.ZeroAddress(msg.Newowner2) {
		uk.RecieveUTXO(ctx, msg.Newowner2, msg.Denom2, oldUTXOs, 1, *txIndex)
	}

	// Increment txIndex
	if !ctx.IsCheckTx() {
		(*txIndex)++
	}
	// TODO: add some tags so we can search it!
	return sdk.Result{} // TODO
}
