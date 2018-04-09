package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
)

// NewAnteHandler returns an AnteHandler that checks signatures,
// and deducts fees from the first signer.

// TODO: what kind of Mapper to take in as param?
func NewAnteHandler(utxoMapper UTXOMapper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx,)
	(_ sdk.Context, _ sdk.Result, abort bool) 
	{

		// TODO: figure out Go syntax: tx.()
		baseTx, ok := tx.(BaseTx)
		if !ok {
			return ctx, sdk.ErrInternal("tx must be type BaseTx").Result(), true
		}

		sigs := tx.getSignatures()
		if len(sigs) == 0 {
			return ctx,
				sdk.ErrUnauthorized("no signers").Result(),
				true
		}


		msg := tx.GetMsg()

		// Assert that number of signatures is correct.
		// GetSigners returns list of crypto.Address
		var signerAddrs = msg.GetSigners()
		if len(sigs) != len(signerAddrs) {
			return ctx,
				sdk.ErrUnauthorized("wrong number of signers").Result(),
				true
		}

		// TODO: calculate sum(msg input) = sum(msg output) + fee
		owner1 := msg.Owner1
		position1 := [3]uint{msg.Blknum1, msg.Txindex1, msg.Oindex1}
		utxo1 := utxoMapper.GetUTXO(ctx, owner1, position1)
		if utxo1 == (BaseUTXO{}) {
			// TODO: return error
		}

		owner2 := msg.Owner2
		position2 := [3]uint{msg.Blknum2, msg.Txindex2, msg.Oindex2}
		utxo2 := utxoMapper.GetUTXO(ctx, owner2, position2)
		if utxo2 == (BaseUTXO{}) {
			// TODO: return error
		}

		// check that the sum of inputs is greater than sum of outputs
		// Should we still include fee in SpendMsg
		if utxo1.GetDenom() + utxo2.GetDenom() != msg.Denom1 + msg.Denom2 + msg.fee {
			// TODO: return error
		}



		// Get the sign bytes (requires all sequence numbers and the fee)
		// sequences := make([]int64, len(signerAddrs))
		// for i := 0; i < len(signerAddrs); i++ {
			// sequences[i] = sigs[i].Sequence
		// }

		// fee := stdTx.Fee

		// XXX: major hack; need to get ChainID
		// into the app right away (#565)
		// LL: not sure what ChainID and viper is?
		// chainID := ctx.ChainID()
		// if chainID == "" {
		// 	chainID = viper.GetString("chain-id")
		// }

		// // signBytes := sdk.StdSignBytes(ctx.ChainID(), sequences, fee, msg)
		// signBytes := msg.GetSignBytes()


		// // Check sig and nonce and collect signer accounts.
		// // LL: probably need to make these into UTXOs instead?
		// var signerAccs = make([]sdk.Account, len(signerAddrs))
		// for i := 0; i < len(sigs); i++ {
		// 	signerAddr, sig := signerAddrs[i], sigs[i]

		// 	// check signature, return account with incremented nonce
		// 	signerAcc, res := processSig(
		// 		ctx, accountMapper,
		// 		signerAddr, sig, signBytes,
		// 	)
		// 	if !res.IsOK() {
		// 		return ctx, res, true
		// 	}

		// 	// first sig pays the fees
		// 	if i == 0 {
		// 		// TODO: min fee
		// 		if !fee.Amount.IsZero() {
		// 			signerAcc, res = deductFees(signerAcc, fee)
		// 			if !res.IsOK() {
		// 				return ctx, res, true
		// 			}
		// 		}
		// 	}

		// 	// Save the account.
		// 	accountMapper.SetAccount(ctx, signerAcc)
		// 	signerAccs[i] = signerAcc
		// }

		// // cache the signer accounts in the context
		// // LL: WithSigners is only used for testing!! (see x/auth/context_test.go)
		// ctx = WithSigners(ctx, signerAccs)

		// // TODO: tx tags (?)

		// return ctx, sdk.Result{}, false // continue...



	}
}


// verify the signature and increment the sequence.
// if the account doesn't have a pubkey, set it.
func processSig(
	ctx sdk.Context, um UTXOMapper,
	addr sdk.Address, sig sdk.StdSignature, signBytes []byte) {

	// // Get the account.
	// acc = am.GetAccount(ctx, addr)
	// if acc == nil {
	// 	return nil, sdk.ErrUnknownAddress(addr.String()).Result()
	// }

	// // Check and increment sequence number.
	// seq := acc.GetSequence()
	// if seq != sig.Sequence {
	// 	return nil, sdk.ErrInvalidSequence(
	// 		fmt.Sprintf("Invalid sequence. Got %d, expected %d", sig.Sequence, seq)).Result()
	// }
	// acc.SetSequence(seq + 1)

	// If pubkey is not known for account,
	// set it from the StdSignature.

	pubKey = sig.PubKey
	if pubKey.Empty() {
		return nil, sdk.ErrInvalidPubKey("PubKey not found").Result()
	}

	if !bytes.Equal(pubKey.Address(), addr) {
		return nil, sdk.ErrInvalidPubKey(
			fmt.Sprintf("PubKey does not match Signer address %v", addr)).Result()
	}
	err := acc.SetPubKey(pubKey)
	if err != nil {
		return nil, sdk.ErrInternal("setting PubKey on signer's account").Result()
	}


	// Check sig.
	if !pubKey.VerifyBytes(signBytes, sig.Signature) {
		return nil, sdk.ErrUnauthorized("signature verification failed").Result()
	}

	return
}





