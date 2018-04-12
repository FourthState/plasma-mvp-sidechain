package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
)

// Handle all "bank" type messages.
func NewHandler(uk UTXOKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SpendMsg:
			return handleSpendMsg(ctx, uk, msg)
		case FinalizeMsg:
			return handleFinalizeMsg(ctx, uk, msg)
		default:
			errMsg := "Unrecognized Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle SendMsg.
func handleSpendMsg(ctx sdk.Context, uk UTXOKeeper, msg SpendMsg) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked
	// TODO: Implement
	if msg.IsDeposit() {
		return handleDepositMsg(ctx, uk, msg)
	}
	if msg.Owner1 != nil {
		position := [3]uint{msg.Blknum1, msg.Txindex1, msg.Oindex1}
		err := uk.SpendUTXO(ctx, msg.Owner1, position)
		if err != nil {
			return err.Result()
		}
	}
	if msg.Owner2 != nil && !ZeroAddress(msg.Owner2) {
		position := [3]uint{msg.Blknum2, msg.Txindex2, msg.Oindex2}
		err := uk.SpendUTXO(ctx, msg.Owner2, position)
		if err != nil {
			return err.Result()
		}
	}
	if msg.Newowner1 != nil {
		err := uk.RecieveUTXO(ctx, msg.Newowner1, msg.Denom1)
		if err != nil {
			return err.Result()
		}
	}
	if msg.Newowner2 != nil && !ZeroAddress(msg.Newowner2) {
		err := uk.RecieveUTXO(ctx, msg.Newowner2, msg.Denom2)
		if err != nil {
			return err.Result()
		}
	}
	// TODO: add some tags so we can search it!
	return sdk.Result{} // TODO
}

// Handle IssueMsg.
func handleDepositMsg(ctx sdk.Context, uk UTXOKeeper, msg SpendMsg) sdk.Result {
	err := uk.RecieveUTXO(ctx, msg.Newowner1, msg.Denom1)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleFinalizeMsg(ctx sdk.Context, uk UTXOKeeper, msg FinalizeMsg) sdk.Result {
	err := uk.FinalizeUTXO(ctx, msg.Spend.Newowner1, msg.Spend.Denom1, msg.Position, msg.ConfirmSigs)
	if err != nil {
		return sdk.NewError(100, err.Error()).Result()
	}
	if msg.Spend.Newowner2 != nil && ZeroAddress(msg.Spend.Newowner2) {
		err = uk.FinalizeUTXO(ctx, msg.Spend.Newowner2, msg.Spend.Denom2, msg.Position, msg.ConfirmSigs)
	}
	if err != nil {
		return sdk.NewError(100, err.Error()).Result()
	}
	return sdk.Result{}
}
