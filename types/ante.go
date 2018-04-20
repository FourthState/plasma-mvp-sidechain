package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
)

// NewAnteHandler returns an AnteHandler that checks signatures,
// and deducts fees from the first signer.
func NewAnteHandler(utxoMapper UTXOMapper, spentDeposits DepositMapper) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx,
	) (_ sdk.Context, _ sdk.Result, abort bool) {
		sigs := tx.GetSignatures()
		if len(sigs) == 0 {
			return ctx,
				sdk.ErrUnauthorized("no signers").Result(),
				true
		}


		msg := tx.GetMsg()

		baseTx, ok := tx.(BaseTx)
		if !ok {
			return ctx, sdk.ErrInternal("tx must be in form of BaseTx").Result(), true
		}

		// Assert that number of signatures is correct.
		// GetSigners returns list of crypto.Address
		var signerAddrs = msg.GetSigners()
		if len(sigs) != len(signerAddrs) {
			return ctx,
				sdk.ErrUnauthorized("wrong number of signers").Result(),
				true
		}

		// 1. Add in ValidateBasic inputs = outputs + fee //spendmsg
		//
		// 2. are input positions in utxoMapper? // ante
		//
		// 3. If not, is deposit? Not in spentDeposits && in rootchain at blknum1 // ante
		//
		// 4. Signature of inputs is valid from sender // ante
		//
		// 5. Confirm Sigs belong to previous owner. // ante
		// 		A -> B -> C. C is UTXO, spending of B (past UTXO) required
		// 		con of A, B -> C destroys UTXO A. B and C still in system. 
		//
		// Ways of dealing with fee:
		// 1. isCheckTx() don't worry about fee
		// 2. if txIndex == 2^16 - 1
		// 3. Since we know poisition of fee utxo, accumlate fees to that position
		// 4. Then check that created msg matches fee


		spendMsg, ok: = msg.(SpendMsg)
		if !ok {
			return ctx, sdk.ErrInternal("Msg must be of type SpendMsg").Result(), true
		}
		signBytes := spendMsg.GetSignBytes() // should this be just the inputs

		// Check for empty sigs

		position1 := [3]uint{spendMsg.Blknum1, spendMsg.Txindex1, spendMsg.Oindex1}
		res := processSig(ctx, utxoMapper, position1, sigs[0], signBytes)
		if !res.IsOK() {
			return ctx, res, true
		}

		utxo1 := utxoMapper.GetUTXO(ctx, position1) 
		if utxo1 == nil {
			// May be a deposit Msg
			//if spentDeposits.GetDeposit(ctx, position1) { //true -> spent
				//return sdk.ErrUnauthorized("Deposit has already been spent").Result(), true
			//}
			// Ping rootContract

			// TODO: return error
		}

		// check for empty sig

		position2 := [3]uint{spendMsg.Blknum2, spendMsg.Txindex2, spendMsg.Oindex2}
		res = processSig(ctx, utxoMapper, position2, sigs[1], signBytes)
		if !res.IsOK() {
			return ctx, res, true
		}


		utxo2 := utxoMapper.GetUTXO(ctx, position2)
		if utxo2 == nil {
			// May be a deposit Msg
			//if spentDeposits.GetDeposit(ctx, position2) { //true -> spent
				//return sdk.ErrUnauthorized("Deposit has already been spent").Result(), true
			//}
			// Ping RootContract

			// TODO: return error
		}

		// If DeliverTx() update fee
		// Rough outline of dealing with fees
		if !ctx.isCheckTx() {
			feePosition := [3]uint{ctx.BlockHeight * 1000, 65535, 0} //adjust based on where feeutxo is
			feeUTXO := utxoMapper.GetUTXO(ctx, feePosition)
			if txIndex == 65535 { //is fee msg
				if feeUTXO.Denom1 != spendMsg.Denom1 {
					return ErrUnauthorized("Fees collected does not match fees reported").Result(), true
				}
			} else {
				// Is not fee Msg
				fee := spendMsg.Fee
				// first transaction in a block
				if feeUTXO != nil {
					fee = fee + feeUTXO.Denom1
					utxoMapper.DeleteUTXO(ctx, feePosition, fee)
				}
				feeUTXO = NewBaseUTXO(nil, nil, nil, nil, fee, feePosition)
				utxoMapper.AddUTXO(ctx, feeUTXO)
			}
		}	

		// LL: not sure what ChainID and viper is?
		// CA: If we want to do cross chain transactions (in the future) 
		// 	   we will need to sign with a chain-id and add to root contract?
		//		
		//chainID := ctx.ChainID()
		//if chainID == "" {
			//chainID = viper.GetString("chain-id")
		//}

		// signBytes := sdk.StdSignBytes(chain-id, fee, msg)

		// // cache the signer accounts in the context
		// // LL: WithSigners is only used for testing!! (see x/auth/context_test.go)
		// ctx = WithSigners(ctx, signerAccs)

		// // TODO: tx tags (?)

		return ctx, sdk.Result{}, false // continue...
	}
}


// verify the signature
// if the account doesn't have a pubkey, set it.
func processSig(
	ctx sdk.Context, um UTXOMapper,
	position [3]uint, sig sdk.StdSignature, signBytes []byte) (
	res sdk.Result) {

	// Get the utxo.
	utxo = um.GetUTXO(ctx, position)
	if utxo == nil {
		return sdk.ErrUnknownRequest("UTXO trying to be spent, does not exist").Result()
	}


	// If pubkey is not known for account,
	// set it from the StdSignature.

	pubKey := utxo.GetPubKey()
	if pubKey.Empty() {
		pubKey = sig.PubKey
		if pubKey.Empty() {
			return sdk.ErrInvalidPubKey("PubKey not found").Result()
		}
		if !bytes.Equal(pubKey.Address(), utxo.GetAddress()) {
			return sdk.ErrInvalidPubKey("PubKey does not match signer address").Result()
		}

		err := utxo.SetPubKey(pubKey)
		if err != nil {
			return sdk.ErrInternal("setting PubKey on signer's account").Result()
		}
	}

	// Check sig.
	if !pubKey.VerifyBytes(signBytes, sig.Signature) {
		return sdk.ErrUnauthorized("signature verification failed").Result()
	}

	return sdk.Result{}
}

func processConfirmSig(
	ctx sdk.Context, utxoMapper UTXOMapper,
	position [3]uint, sig sdk.StdSignature, signBytes []byte) (
	res sdk.Result) {
	
	utxo = um.GetUTXO(ctx, position)
	if utxo == nil {
		return sdk.ErrUnknownRequest("UTXO trying to be spent, does not exist").Result()
	}

	pubKey := utxo.GetCSPubKey()
	if pubKey.Empty() {
		pubKey = sig.PubKey
		if pubKey.Empty() {
			return sdk.ErrInvalidPubKey("PubKey not found").Result()
		}

		if !bytes.Equal(pubKey.Address(), utxo.GetCSAddress()) {
			return sdk.ErrInvalidPubKey("PubKey does not match signer address").Result()
		}

		err := utxo.SetCSPubKey(pubKey)
		if err != nil {
			return sdk.ErrInternal("setting PubKey on signer's account").Result()
		}
	}

	// Check sig.
	if !pubKey.VerifyBytes(signBytes, sig.Signature) {
		return sdk.ErrUnauthorized("signature verification failed").Result()
	}

	return sdk.Result{}
}