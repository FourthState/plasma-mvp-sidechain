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
func handleSpendMsg(ctx sdk.Context, uk utxo.UTXOKeeper, msg types.SpendMsg, txIndex *uint16) sdk.Result {
	newInputAddrs := handleInputs(ctx, uk, msg)

	handleOutputs(ctx, uk, msg, txIndex, newInputAddrs)

	// Increment txIndex
	if !ctx.IsCheckTx() {
		(*txIndex)++
	}
	// TODO: add some tags so we can search it!
	return sdk.Result{} // TODO
}

// spend the inputs of transaction
func handleInputs(ctx sdk.Context, uk utxo.UTXOKeeper, msg types.SpendMsg) [2]common.Address {
	// spend first input from the spend msg
	position1 := types.Position{msg.Blknum1, msg.Txindex1, msg.Oindex1, msg.DepositNum1}
	inputAddr1 := msg.Owner1
	utxo1 := uk.GetUTXO(ctx, inputAddr1, position1)
	uk.SpendUTXO(ctx, inputAddr1, position1)

	utxo2 := types.BaseUTXO{}

	// spend second input if it exists
	inputAddr2 := msg.Owner2
	if !utils.ZeroAddress(inputAddr2) {
		position2 := types.Position{msg.Blknum2, msg.Txindex2, msg.Oindex2, msg.DepositNum2}
		utxo2 = uk.GetUTXO(ctx, inputAddr2, position2)
		uk.SpendUTXO(ctx, inputAddr2, position2)
	}

	return [2]common.Address{utxo1.GetAddress(), utxo2.GetAddress()}

}

func handleOutputs(ctx sdk.Context, uk utxo.UTXOKeeper, msg types.SpendMsg, txIndex *uint16, inputAddrs [2]common.Address) {
	// create first output
	newUTXO := types.BaseUTXO{
		InputAddresses: inputAddrs,
		Address:        msg.Newowner1,
		Amount:         msg.Denom1,
		Denom:          "plasma", // TODO: change (add field into spendMsg?)
		Position:       types.NewPosition(uint64(ctx.BlockHeight()), txIndex, 0, 0),
	}
	uk.RecieveUTXO(ctx, newUTXO)

	// create second output only if the fields for a second input are valid
	if !utils.ZeroAddress(msg.Newowner2) {
		newUTXO = types.BaseUTXO{
			InputAddresses: inputAddrs,
			Address:        msg.Newowner2,
			Amount:         msg.Denom2,
			Denom:          "plasma",
			Position:       types.NewPosition(uint64(ctx.BlockHeight()), txIndex, 1, 0),
		}
		uk.RecieveUTXO(ctx, newUTXO)
	}
}
