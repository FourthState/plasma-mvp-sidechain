package types

import (
	"reflect"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Handle all "bank" type messages.
func NewHandler(uk UTXOKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SpendMsg:
			return handleSpendMsg(ctx, uk, msg)
		case DepositMsg:
			return handleDepositMsg(ctx, uk, msg)
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
	if msg.Owner1 != nil {
		//CHANGE
		err := uk.SpendUTXO(ctx, msg.Owner1, msg.Denom1)
		if err != nil {
			return err.Result()
		}
	}
	if msg.Owner2 != nil {
		//CHANGE
		err := uk.SpendUTXO(ctx, msg.Owner2, msg.Denom2)
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
	if msg.Newowner2 != nil {
		err := uk.RecieveUTXO(ctx, msg.Newowner2, msg.Denom2)
		if err != nil {
			return err.Result()
		}
	}
	// TODO: add some tags so we can search it!
	return sdk.Result{} // TODO
}

// Handle IssueMsg.
func handleDepositMsg(ctx sdk.Context, uk UTXOKeeper, msg DepositMsg) sdk.Result {
	// TODO: Implement 
	return sdk.Result{}
}