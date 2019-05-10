package handlers

import (
	"bytes"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// the reason for an interface is to allow the connection object
// to be cooked when testing the ante handler
type plasmaConn interface {
	GetDeposit(*big.Int, *big.Int) (plasma.Deposit, *big.Int, bool)
	HasTxBeenExited(*big.Int, plasma.Position) bool
}

func NewAnteHandler(txStore store.TxStore, depositStore store.DepositStore, blockStore store.BlockStore, client plasmaConn) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {
		msg := tx.GetMsgs()[0] // tx should only have one msg
		switch mtype := msg.Type(); mtype {
		case "include_deposit":
			depositMsg := msg.(msgs.IncludeDepositMsg)
			return includeDepositAnteHandler(ctx, depositStore, blockStore, depositMsg, client)
		case "spend_utxo":
			spendMsg := msg.(msgs.SpendMsg)
			return spendMsgAnteHandler(ctx, spendMsg, txStore, depositStore, blockStore, client)
		default:
			return ctx, msgs.ErrInvalidTransaction(DefaultCodespace, "msg is not of type SpendMsg or IncludeDepositMsg").Result(), true
		}
	}
}

func spendMsgAnteHandler(ctx sdk.Context, spendMsg msgs.SpendMsg, txStore store.TxStore, depositStore store.DepositStore, blockStore store.BlockStore, client plasmaConn) (newCtx sdk.Context, res sdk.Result, abort bool) {
	var totalInputAmt, totalOutputAmt *big.Int

	// attempt to recover signers
	signers := spendMsg.GetSigners()
	if len(signers) == 0 {
		return ctx, msgs.ErrInvalidTransaction(DefaultCodespace, "failed recovering signers").Result(), true
	}

	/* validate the first input */
	amt, res := validateInput(ctx, spendMsg.Input0, common.BytesToAddress(signers[0]), txStore, blockStore, client)
	if !res.IsOK() {
		return ctx, res, true
	}

	// must cover the fee
	if amt.Cmp(spendMsg.Fee) < 0 {
		return ctx, ErrInsufficientFee(DefaultCodespace, "first input has an insufficient amount to pay the fee").Result(), true
	}

	totalInputAmt = amt

	/* validate second input if applicable */
	if spendMsg.HasSecondInput() {
		if len(signers) != 2 {
			return ctx, msgs.ErrInvalidTransaction(DefaultCodespace, "second signature not present").Result(), true
		}
		amt, res = validateInput(ctx, spendMsg.Input1, common.BytesToAddress(signers[1]), txStore, blockStore, client)
		if !res.IsOK() {
			return ctx, res, true
		}

		totalInputAmt = totalInputAmt.Add(totalInputAmt, amt)
	}

	// input0 + input1 (totalInputAmt) == output0 + output1 + Fee (totalOutputAmt)
	totalOutputAmt = spendMsg.Output0.Amount
	if spendMsg.HasSecondOutput() {
		totalOutputAmt = new(big.Int).Add(new(big.Int).Add(totalOutputAmt, spendMsg.Output1.Amount), spendMsg.Fee)
	} else {
		totalOutputAmt = new(big.Int).Add(totalOutputAmt, spendMsg.Fee)
	}

	if totalInputAmt.Cmp(totalOutputAmt) != 0 {
		return ctx, msgs.ErrInvalidTransaction(DefaultCodespace, "inputs do not equal Outputs (+ fee)").Result(), true
	}

	return ctx, sdk.Result{}, false
}

// validates the inputs against the utxo store and returns the amount of the respective input
func validateInput(ctx sdk.Context, input plasma.Input, signer common.Address, txStore store.TxStore, blockStore store.BlockStore, client plasmaConn) (*big.Int, sdk.Result) {
	var amt *big.Int

	// inputUTXO must be owned by the signer due to the prefix so we do not need to
	// check the owner of the position
	inputUTXO, ok := utxoStore.GetUTXO(ctx, signer, input.Position)
	if !ok {
		return nil, msgs.ErrInvalidTransaction(DefaultCodespace, "input, %v, does not exist by owner %x", input.Position, signer).Result()
	}
	if inputUTXO.Spent {
		return nil, msgs.ErrInvalidTransaction(DefaultCodespace, "input, %v, already spent", input.Position).Result()
	}
	if client.HasTxBeenExited(plasmaStore.CurrentPlasmaBlockNum(ctx), input.Position) {
		return nil, ErrExitedInput(DefaultCodespace, "input, %v, utxo has exitted", input.Position).Result()
	}

	// validate confirm signatures if not a fee utxo or deposit utxo
	if input.TxIndex < 1<<16-1 && input.DepositNonce.Sign() == 0 {
		res := validateConfirmSignatures(ctx, input, inputUTXO)
		if !res.IsOK() {
			return nil, res
		}
	}

	// check if the parent utxo has exited
	for _, key := range inputUTXO.InputKeys {
		utxo, _ := utxoStore.GetUTXOWithKey(ctx, key)
		if client.HasTxBeenExited(plasmaStore.CurrentPlasmaBlockNum(ctx), utxo.Position) {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("a parent of the input has exited. Owner: %x, Position: %v", utxo.Output.Owner, utxo.Position)).Result()
		}
	}

	amt = inputUTXO.Output.Amount

	return amt, sdk.Result{}
}

// validates the input's confirm signatures
func validateConfirmSignatures(ctx sdk.Context, input plasma.Input, inputTx store.Transaction) sdk.Result {
	if len(input.ConfirmSignatures) > 0 && (inputTx.Transaction.HasSecondInput) {
		return msgs.ErrInvalidTransaction(DefaultCodespace, "incorrect number of confirm signatures").Result()
	}

	confirmationHash := utils.ToEthSignedMessageHash(inputTx.ConfirmationHash[:])[:]
	for i, key := range inputUTXO.InputKeys {
		address := key[:common.AddressLength]

		pubKey, _ := crypto.SigToPub(confirmationHash, input.ConfirmSignatures[i][:])
		signer := crypto.PubkeyToAddress(*pubKey)
		if !bytes.Equal(signer[:], address) {
			return sdk.ErrUnauthorized(fmt.Sprintf("confirm signature not signed by the correct address. Got: %x. Expected %x", signer, address)).Result()
		}
	}

	return sdk.Result{}
}

func includeDepositAnteHandler(ctx sdk.Context, depositStore store.DepositStore, blockStore store.BlockStore, msg msgs.IncludeDepositMsg, client plasmaConn) (newCtx sdk.Context, res sdk.Result, abort bool) {
	if depositStore.HasDeposit(ctx, msg.Owner, msg.DepositNonce) {
		return ctx, msgs.ErrInvalidTransaction(DefaultCodespace, "deposit, %s, already exists in store", msg.DepositNonce.String()).Result(), true
	}
	_, threshold, ok := client.GetDeposit(blockStore.CurrentPlasmaBlockNum(ctx), msg.DepositNonce)
	if !ok && threshold == nil {
		return ctx, msgs.ErrInvalidTransaction(DefaultCodespace, "deposit, %s, does not exist.", msg.DepositNonce.String()).Result(), true
	}
	if !ok {
		return ctx, msgs.ErrInvalidTransaction(DefaultCodespace, "deposit, %s, has not finalized yet. Please wait at least %d blocks before resubmitting", msg.DepositNonce.String(), threshold.Int64()).Result(), true
	}
	exited := client.HasTxBeenExited(blockStore.CurrentPlasmaBlockNum(ctx), depositPosition)
	if exited {
		return ctx, msgs.ErrInvalidTransaction(DefaultCodespace, "deposit, %s, has already exitted from rootchain", msg.DepositNonce.String()).Result(), true
	}
	return ctx, sdk.Result{}, false
}
