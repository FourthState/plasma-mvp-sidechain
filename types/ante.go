package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	//abci "github.com/tendermint/abci/types"
	//"github.com/spf13/viper"
)

// NewAnteHandler returns an AnteHandler that checks signatures,
// and deducts fees from the first signer.
func NewAnteHandler(utxoMapper UTXOMapper) sdk.AnteHandler {
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

		_, ok := tx.(BaseTx)
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


		spendMsg, ok := msg.(SpendMsg)
		if !ok {
			return ctx, sdk.ErrInternal("Msg must be of type SpendMsg").Result(), true
		}
		signBytes := spendMsg.GetSignBytes() // should this be just the inputs

		// Check for empty sigs

		position1 := Position{spendMsg.Blknum1, spendMsg.Txindex1, spendMsg.Oindex1}
		res := processSig(ctx, utxoMapper, position1, sigs[0], signBytes)
		if !res.IsOK() {
			return ctx, res, true
		}
		
		//Change signBytes to correct value
		res = processConfirmSig(ctx, utxoMapper, position1, spendMsg.ConfirmSigs1,signBytes)
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

		position2 := Position{spendMsg.Blknum2, spendMsg.Txindex2, spendMsg.Oindex2}
		res = processSig(ctx, utxoMapper, position2, sigs[1], signBytes)
		if !res.IsOK() {
			return ctx, res, true
		}

		res = processConfirmSig(ctx, utxoMapper, position2, spendMsg.ConfirmSigs2,signBytes)
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
		if !ctx.IsCheckTx() {
			header := ctx.BlockHeader()
			feeTxIndex := uint16(header.GetNumTxs()) - 1
			feePosition := Position{uint64(ctx.BlockHeight()) * 1000, feeTxIndex, 0} //adjust based on where feeutxo is
			feeUTXO := utxoMapper.GetUTXO(ctx, feePosition)
			if GetTxIndex(ctx) == feeTxIndex { //is fee msg
				if feeUTXO.GetDenom() != spendMsg.Denom1 {
					return ctx, sdk.ErrUnauthorized("Fees collected does not match fees reported").Result(), true
				}
			} else {
				// Is not fee Msg
				fee := spendMsg.Fee
				// first transaction in a block
				if feeUTXO != nil {
					fee = fee + feeUTXO.GetDenom()
					utxoMapper.DeleteUTXO(ctx, feePosition)
				}
				
				feeUTXO = NewBaseUTXO(crypto.Address([]byte("")),[2]crypto.Address{crypto.Address([]byte("")),
				crypto.Address([]byte(""))}, crypto.PubKeySecp256k1{}, [2]crypto.PubKey{crypto.PubKeySecp256k1{}, crypto.PubKeySecp256k1{}}, fee, feePosition)
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
				// if the account doesn't have a pubkey, set it.}
func processSig(
	ctx sdk.Context, um UTXOMapper,
	position Position, sig sdk.StdSignature, signBytes []byte) (
	res sdk.Result) {

	// Get the utxo.
	utxo := um.GetUTXO(ctx, position)
	if utxo == nil {
		return sdk.ErrUnknownRequest("UTXO trying to be spent, does not exist").Result()
	}


	// If pubkey is not known for account,
	// set it from the StdSignature.

	pubKey := utxo.GetPubKey()
	if pubKey == nil{
		pubKey = sig.PubKey
		if pubKey == nil {
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
	position Position, sig [2]crypto.Signature, signBytes []byte) (
	res sdk.Result) {
	
	utxo := utxoMapper.GetUTXO(ctx, position)
	if utxo == nil {
		return sdk.ErrUnknownRequest("UTXO trying to be spent, does not exist").Result()
	}

	csPubKeys := utxo.GetCSPubKey()
	pubKey1 := csPubKeys[0]
	pubKey2 := csPubKeys[1]
	if pubKey1 == nil || pubKey2 == nil {
		return sdk.ErrInvalidPubKey("PubKey not found").Result()
	}

	// Check sig.
	if !pubKey1.VerifyBytes(signBytes, sig[0]) {
		return sdk.ErrUnauthorized("signature verification failed").Result()
	}

	if !pubKey2.VerifyBytes(signBytes, sig[1]) {
		return sdk.ErrUnauthorized("signature verification failed").Result()
	}

	return sdk.Result{}
}