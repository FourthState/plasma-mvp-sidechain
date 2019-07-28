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
	HasTxExited(*big.Int, plasma.Position) (bool, error)
}

// NewAnteHandler returns an ante handler capable of handling include_deposit
// and spend_utxo Msgs.
func NewAnteHandler(ds store.DataStore, client plasmaConn) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {
		msg := tx.GetMsgs()[0] // tx should only have one msg
		switch mtype := msg.Type(); mtype {
		case "include_deposit":
			depositMsg := msg.(msgs.IncludeDepositMsg)
			return includeDepositAnteHandler(ctx, depositMsg, ds, client)
		case "spend_utxo":
			spendMsg := msg.(msgs.SpendMsg)
			return spendMsgAnteHandler(ctx, spendMsg, ds, client)
		default:
			return ctx, ErrInvalidTransaction("msg is not of type SpendMsg or IncludeDepositMsg").Result(), true
		}
	}
}

func spendMsgAnteHandler(ctx sdk.Context, spendMsg msgs.SpendMsg, ds store.DataStore, client plasmaConn) (newCtx sdk.Context, res sdk.Result, abort bool) {
	var totalInputAmt, totalOutputAmt *big.Int
	totalInputAmt = big.NewInt(0)
	totalOutputAmt = big.NewInt(0)

	// attempt to recover signers
	signers := spendMsg.GetSigners()
	if len(signers) == 0 {
		return ctx, ErrInvalidTransaction("failed recovering signers").Result(), true
	}

	if len(signers) != len(spendMsg.Inputs) {
		return ctx, ErrInvalidSignature("number of signers does not equal number of signatures").Result(), true
	}

	/* validate inputs */
	for i, signer := range signers {
		amt, res := validateInput(ctx, spendMsg.Inputs[i], common.BytesToAddress(signer), ds, client)
		if !res.IsOK() {
			return ctx, res, true
		}

		// first input must cover the fee
		if i == 0 && amt.Cmp(spendMsg.Fee) < 0 {
			return ctx, ErrInsufficientFee("first input has an insufficient amount to pay the fee").Result(), true
		}

		totalInputAmt = totalInputAmt.Add(totalInputAmt, amt)
	}

	// input0 + input1 (totalInputAmt) == output0 + output1 + Fee (totalOutputAmt)
	for _, output := range spendMsg.Outputs {
		totalOutputAmt = new(big.Int).Add(totalOutputAmt, output.Amount)
	}
	totalOutputAmt = new(big.Int).Add(totalOutputAmt, spendMsg.Fee)

	if totalInputAmt.Cmp(totalOutputAmt) != 0 {
		return ctx, ErrInvalidTransaction("inputs do not equal Outputs (+ fee)").Result(), true
	}

	return ctx, sdk.Result{}, false
}

// validates the inputs against the output store and returns the amount of the respective input
func validateInput(ctx sdk.Context, input plasma.Input, signer common.Address, ds store.DataStore, client plasmaConn) (*big.Int, sdk.Result) {
	var amt *big.Int

	// inputUTXO must be owned by the signer due to the prefix so we do not need to
	// check the owner of the position
	inputUTXO, ok := ds.GetOutput(ctx, input.Position)
	if !ok {
		return nil, ErrInvalidInput("input, %v, does not exist", input.Position).Result()
	}
	if !bytes.Equal(inputUTXO.Output.Owner[:], signer[:]) {
		return nil, ErrSignatureVerificationFailure(fmt.Sprintf("transaction was not signed by correct address. Got: 0x%x. Expected: 0x%x", signer, inputUTXO.Output.Owner)).Result()
	}
	if inputUTXO.Spent {
		return nil, ErrInvalidInput("input, %v, already spent", input.Position).Result()
	}
	exited, err := client.HasTxExited(ds.PlasmaBlockHeight(ctx), input.Position)
	if err != nil {
		return nil, ErrInvalidInput("failed to retrieve exit information on input, %v", input.Position).Result()
	} else if exited {
		return nil, ErrExitedInput("input, %v, utxo has exitted", input.Position).Result()
	}

	// validate inputs/confirmation signatures if not a fee utxo or deposit utxo
	if !input.IsDeposit() && !input.IsFee() {
		tx, ok := ds.GetTxWithPosition(ctx, input.Position)
		if !ok {
			return nil, sdk.ErrInternal(fmt.Sprintf("failed to retrieve the transaction that input with position %s belongs to", input.Position)).Result()
		}

		res := validateConfirmSignatures(ctx, input, tx, ds)
		if !res.IsOK() {
			return nil, res
		}

		// check if the parent utxo has exited
		for _, in := range tx.Transaction.Inputs {
			exited, err = client.HasTxExited(ds.PlasmaBlockHeight(ctx), in.Position)
			if err != nil {
				return nil, ErrInvalidInput(fmt.Sprintf("failed to retrieve exit information on input, %v", in.Position)).Result()
			} else if exited {
				return nil, ErrExitedInput(fmt.Sprintf("a parent of the input has exited. Position: %v", in.Position)).Result()
			}
		}
	}

	amt = inputUTXO.Output.Amount

	return amt, sdk.Result{}
}

// validates the input's confirm signatures
func validateConfirmSignatures(ctx sdk.Context, input plasma.Input, inputTx store.Transaction, ds store.DataStore) sdk.Result {
	if len(input.ConfirmSignatures) > 0 && len(inputTx.Transaction.Inputs) != len(input.ConfirmSignatures) {
		return ErrInvalidTransaction("incorrect number of confirm signatures").Result()
	}
	confirmationHash := utils.ToEthSignedMessageHash(inputTx.ConfirmationHash[:])[:]
	for i, in := range inputTx.Transaction.Inputs {
		utxo, ok := ds.GetOutput(ctx, in.Position)
		if !ok {
			return ErrInvalidInput(fmt.Sprintf("failed to retrieve input utxo with position %s", in.Position)).Result()
		}

		pubKey, _ := crypto.SigToPub(confirmationHash, input.ConfirmSignatures[i][:])
		signer := crypto.PubkeyToAddress(*pubKey)
		if !bytes.Equal(signer[:], utxo.Output.Owner[:]) {
			return ErrSignatureVerificationFailure(fmt.Sprintf("confirm signature not signed by the correct address. Got: %x. Expected: %x", signer, utxo.Output.Owner)).Result()
		}
	}

	return sdk.Result{}
}

func includeDepositAnteHandler(ctx sdk.Context, msg msgs.IncludeDepositMsg, ds store.DataStore, client plasmaConn) (newCtx sdk.Context, res sdk.Result, abort bool) {
	if ds.HasDeposit(ctx, msg.DepositNonce) {
		return ctx, ErrInvalidTransaction("deposit, %s, already exists in store", msg.DepositNonce.String()).Result(), true
	}
	deposit, threshold, ok := client.GetDeposit(ds.PlasmaBlockHeight(ctx), msg.DepositNonce)
	if !ok && threshold == nil {
		return ctx, ErrInvalidTransaction("deposit, %s, does not exist.", msg.DepositNonce.String()).Result(), true
	}
	if !ok {
		return ctx, ErrInvalidTransaction("deposit, %s, has not finalized yet. Please wait at least %d blocks before resubmitting", msg.DepositNonce.String(), threshold.Int64()).Result(), true
	}
	if !bytes.Equal(deposit.Owner[:], msg.Owner[:]) {
		return ctx, ErrInvalidTransaction("deposit, %s, does not equal the owner specified in the include-deposit Msg", msg.DepositNonce).Result(), true
	}

	depositPosition := plasma.NewPosition(big.NewInt(0), 0, 0, msg.DepositNonce)
	exited, err := client.HasTxExited(ds.PlasmaBlockHeight(ctx), depositPosition)
	if err != nil {
		return ctx, ErrInvalidTransaction("failed to retrieve deposit information for deposit, %s", msg.DepositNonce.String()).Result(), true
	} else if exited {
		return ctx, ErrInvalidTransaction("deposit, %s, has already exitted from rootchain", msg.DepositNonce.String()).Result(), true
	}
	if !bytes.Equal(msg.Owner.Bytes(), deposit.Owner.Bytes()) {
		return ctx, ErrInvalidTransaction(fmt.Sprintf("msg has the wrong owner field for given deposit. Resubmit with correct deposit owner: %s", deposit.Owner.String())).Result(), true
	}
	return ctx, sdk.Result{}, false
}
