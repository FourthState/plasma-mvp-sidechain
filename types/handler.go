package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	"reflect"
)

func NewHandler(uk UTXOKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SpendMsg:
			return handleSpendMsg(ctx, uk, msg)
		default:
			errMsg := "Unrecognized Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle SpendMsg.
func handleSpendMsg(ctx sdk.Context, uk UTXOKeeper, msg SpendMsg) sdk.Result {
	position1 := Position{msg.Blknum1, msg.Txindex1, msg.Oindex1}
	position2 := Position{msg.Blknum2, msg.Txindex2, msg.Oindex2}
	parent1 := uk.um.GetUTXO(ctx, position1)
	parent2 := uk.um.GetUTXO(ctx, position2) 
	
	// Inputs
	if msg.Owner1 != nil && !ZeroAddress(msg.Owner1) {
		err := uk.SpendUTXO(ctx, msg.Owner1, position1)
		if err != nil {
			return err.Result()
		}
	}
	if msg.Owner2 != nil && !ZeroAddress(msg.Owner2) {
		err := uk.SpendUTXO(ctx, msg.Owner2, position2)
		if err != nil {
			return err.Result()
		}
	}

	var csAddress [2]crypto.Address
	var csPubKey [2]crypto.PubKey
	// Outputs
	if msg.Newowner1 != nil && !ZeroAddress(msg.Newowner1) && parent1 != nil {
		if parent2 == nil {
			csAddress = [2]crypto.Address{parent1.GetAddress(), crypto.Address{}}
			csPubKey = [2]crypto.PubKey{parent1.GetPubKey(), crypto.PubKeySecp256k1{}}
		} else {
			csAddress = [2]crypto.Address{parent1.GetAddress(), parent2.GetAddress()}
			csPubKey = [2]crypto.PubKey{parent1.GetPubKey(), parent2.GetPubKey()}
		}

		err := uk.RecieveUTXO(ctx, msg.Newowner1, msg.Denom1, csAddress, csPubKey, 0)
		if err != nil {
			return err.Result()
		}
	}
	if msg.Newowner2 != nil && !ZeroAddress(msg.Newowner2) && parent2 != nil {
		if parent1 == nil {
			csAddress = [2]crypto.Address{crypto.Address{}, parent2.GetAddress()}
			csPubKey = [2]crypto.PubKey{crypto.PubKeySecp256k1{}, parent2.GetPubKey()}
		} else {
			csAddress = [2]crypto.Address{parent1.GetAddress(), parent2.GetAddress()}
			csPubKey = [2]crypto.PubKey{parent1.GetPubKey(), parent2.GetPubKey()}
		}

		err := uk.RecieveUTXO(ctx, msg.Newowner2, msg.Denom2, csAddress, csPubKey, 1)
		if err != nil {
			return err.Result()
		}
	}
	
	// TODO: add some tags so we can search it!
	return sdk.Result{} // TODO
}