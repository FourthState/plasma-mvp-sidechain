package handlers

import (
	"crypto/sha256"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

var (
	privKey, _ = crypto.GenerateKey()
	addr       = crypto.PubkeyToAddress(privKey.PublicKey)
	// bad keys to check against the deposit
	badPrivKey, _ = crypto.GenerateKey()
)

type inputUTXO struct {
	BlockNum     *big.Int
	TxIndex      uint16
	OIndex       uint8
	DepositNonce *big.Int
	Address      common.Address
	Spent        bool
}

// cook the plasma connection
type conn struct{}

// all deposits should be in an amount of 10eth owner by addr(defined above)
func (p conn) GetDeposit(nonce *big.Int) (plasma.Deposit, *big.Int, bool) {
	return plasma.Deposit{addr, big.NewInt(10), utils.Big0}, big.NewInt(-2), true
}
func (p conn) HasTxBeenExited(pos plasma.Position) bool { return false }

var _ plasmaConn = conn{}

// cook up different plasma connection that will always claim input exitted
type exitConn struct{}

// all deposits should be in an amount of 10eth owner by addr(defined above)
func (p exitConn) GetDeposit(nonce *big.Int) (plasma.Deposit, *big.Int, bool) {
	return plasma.Deposit{addr, big.NewInt(10), utils.Big0}, big.NewInt(-2), true
}
func (p exitConn) HasTxBeenExited(pos plasma.Position) bool { return true }

func TestAnteChecks(t *testing.T) {
	// setup
	ctx, utxoStore, plasmaStore, presenceClaimStore := setup()
	handler := NewAnteHandler(utxoStore, plasmaStore, presenceClaimStore, conn{})

	// cook up some input UTXOs to start in UTXO store
	inputs := []inputUTXO{
		{nil, 0, 0, utils.Big1, addr, false},
		{nil, 0, 0, utils.Big2, addr, false},
		{nil, 0, 0, big.NewInt(3), addr, true},
	}
	setupInputs(ctx, utxoStore, inputs...)

	type validationCase struct {
		reason string
		msgs.SpendMsg
	}

	// cases to check for. cases with signature checks will get set subsequent to this step
	// array of pointers because we are setting signatures after using `range`
	// since InputKeys not set, confirm signatures will simply be 0 bytes
	invalidCases := []*validationCase{
		&validationCase{
			reason: "incorrect first signature",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil),
					Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big0), [65]byte{}, nil),
					Output0: plasma.NewOutput(addr, big.NewInt(10)),
					Output1: plasma.NewOutput(common.Address{}, utils.Big0),
					Fee:     utils.Big0,
				},
			},
		},
		&validationCase{
			reason: "incorrect second signature",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil),
					Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big2), [65]byte{}, nil),
					Output0: plasma.NewOutput(addr, big.NewInt(20)),
					Output1: plasma.NewOutput(common.Address{}, utils.Big0),
					Fee:     utils.Big0,
				},
			},
		},
		&validationCase{
			reason: "no signatures",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil),
					Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big2), [65]byte{}, nil),
					Output0: plasma.NewOutput(addr, big.NewInt(20)),
					Output1: plasma.NewOutput(common.Address{}, utils.Big0),
					Fee:     utils.Big0,
				},
			},
		},
		&validationCase{
			reason: "invalid fee",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil),
					Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big2), [65]byte{}, nil),
					Output0: plasma.NewOutput(addr, big.NewInt(20)),
					Output1: plasma.NewOutput(common.Address{}, utils.Big0),
					Fee:     big.NewInt(20),
				},
			},
		},
		&validationCase{
			reason: "unbalanced transaction",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil),
					Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big2), [65]byte{}, nil),
					Output0: plasma.NewOutput(addr, big.NewInt(10)),
					Output1: plasma.NewOutput(addr, big.NewInt(10)),
					Fee:     utils.Big1,
				},
			},
		},
		&validationCase{
			reason: "input deposit utxo does not exist",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil),
					Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, big.NewInt(4)), [65]byte{}, nil),
					Output0: plasma.NewOutput(addr, big.NewInt(10)),
					Output1: plasma.NewOutput(addr, big.NewInt(10)),
					Fee:     utils.Big0,
				},
			},
		},
		&validationCase{
			reason: "input transaction utxo does not exist",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil),
					Input1:  plasma.NewInput(plasma.NewPosition(utils.Big1, 3, 1, nil), [65]byte{}, nil),
					Output0: plasma.NewOutput(addr, big.NewInt(10)),
					Output1: plasma.NewOutput(addr, big.NewInt(10)),
					Fee:     utils.Big0,
				},
			},
		},
		&validationCase{
			reason: "input utxo already spent",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil),
					Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, big.NewInt(3)), [65]byte{}, nil),
					Output0: plasma.NewOutput(addr, big.NewInt(10)),
					Output1: plasma.NewOutput(addr, big.NewInt(10)),
					Fee:     utils.Big0,
				},
			},
		},
	}

	// set invalid first signature
	txHash := utils.ToEthSignedMessageHash(invalidCases[0].SpendMsg.TxHash())
	sig, _ := crypto.Sign(txHash, badPrivKey)
	copy(invalidCases[0].SpendMsg.Input0.Signature[:], sig)

	// set invalid second signature but correct first signature
	txHash = utils.ToEthSignedMessageHash(invalidCases[1].SpendMsg.TxHash())
	sig, _ = crypto.Sign(txHash, badPrivKey)
	copy(invalidCases[1].SpendMsg.Input1.Signature[:], sig)
	sig, _ = crypto.Sign(txHash, privKey)
	copy(invalidCases[1].SpendMsg.Input0.Signature[:], sig)

	// set valid signatures for remaining cases
	for _, txCase := range invalidCases[3:] {
		txHash = utils.ToEthSignedMessageHash(txCase.SpendMsg.TxHash())
		sig, _ = crypto.Sign(txHash, privKey)
		copy(txCase.SpendMsg.Input0.Signature[:], sig[:])
		copy(txCase.SpendMsg.Input1.Signature[:], sig[:])
	}

	for _, txCase := range invalidCases {
		_, res, abort := handler(ctx, txCase.SpendMsg, false)
		require.False(t, res.IsOK(), txCase.reason)
		require.True(t, abort, txCase.reason)
	}
}

func TestAnteExitedInput(t *testing.T) {
	// setup
	ctx, utxoStore, plasmaStore, presenceClaimStore := setup()
	handler := NewAnteHandler(utxoStore, plasmaStore, presenceClaimStore, exitConn{})

	// place input in store
	input := inputUTXO{
		BlockNum:     utils.Big1,
		TxIndex:      0,
		OIndex:       0,
		DepositNonce: nil,
		Spent:        false,
		Address:      addr,
	}
	setupInputs(ctx, utxoStore, input)

	// create msg
	spendMsg := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Input0:  plasma.NewInput(plasma.NewPosition(utils.Big1, 0, 0, nil), [65]byte{}, nil),
			Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, nil), [65]byte{}, nil),
			Output0: plasma.NewOutput(addr, big.NewInt(10)),
			Output1: plasma.NewOutput(addr, big.NewInt(9)),
			Fee:     utils.Big1,
		},
	}

	// set signature
	txHash := utils.ToEthSignedMessageHash(spendMsg.TxHash())
	sig, _ := crypto.Sign(txHash, privKey)
	copy(spendMsg.Input0.Signature[:], sig[:])

	_, res, abort := handler(ctx, spendMsg, false)
	require.False(t, res.IsOK(), "Result OK even though input exitted")
	require.True(t, abort, "Did not abort tx even though input exitted")

	// TODO: test case where grandparent exitted but parent didn't
}

func TestAnteInvalidConfirmSig(t *testing.T) {
	// setup
	ctx, utxoStore, plasmaStore, presenceClaimStore := setup()
	handler := NewAnteHandler(utxoStore, plasmaStore, presenceClaimStore, conn{})

	// place input in store
	inputs := []inputUTXO{
		{nil, 0, 0, utils.Big1, addr, false},
		{nil, 0, 0, utils.Big2, addr, true},
	}
	setupInputs(ctx, utxoStore, inputs...)

	parentTx := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big2), [65]byte{}, nil),
			Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, nil), [65]byte{}, nil),
			Output0: plasma.NewOutput(addr, big.NewInt(10)),
			Output1: plasma.NewOutput(common.Address{}, nil),
			Fee:     utils.Big0,
		},
	}

	// set regular transaction utxo in store
	// parent input was 0.0.0.2
	// must create input key and confirmation hash
	// also need confirm sig of parent in order to spend
	inputKey := store.GetUTXOStoreKey(addr, plasma.NewPosition(nil, 0, 0, utils.Big2))
	confBytes := sha256.Sum256(append(parentTx.MerkleHash(), ctx.BlockHeader().DataHash...))
	confHash := utils.ToEthSignedMessageHash(confBytes[:])
	badConfSig, _ := crypto.Sign(confHash, badPrivKey)
	inputUTXO := store.UTXO{
		InputKeys:        [][]byte{inputKey},
		ConfirmationHash: confBytes[:],
		Output: plasma.Output{
			Owner:  addr,
			Amount: big.NewInt(10),
		},
		Position: plasma.NewPosition(utils.Big1, 0, 0, nil),
	}
	utxoStore.StoreUTXO(ctx, inputUTXO)

	// store confirm sig into correct format
	var invalidConfirmSig [65]byte
	copy(invalidConfirmSig[:], badConfSig)

	// create msg
	spendMsg := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Input0:  plasma.NewInput(plasma.NewPosition(utils.Big1, 0, 0, nil), [65]byte{}, [][65]byte{invalidConfirmSig}),
			Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil),
			Output0: plasma.NewOutput(addr, big.NewInt(10)),
			Output1: plasma.NewOutput(addr, big.NewInt(9)),
			Fee:     utils.Big1,
		},
	}

	// set signature
	txHash := utils.ToEthSignedMessageHash(spendMsg.TxHash())
	sig, _ := crypto.Sign(txHash, privKey)
	copy(spendMsg.Input0.Signature[:], sig[:])
	copy(spendMsg.Input1.Signature[:], sig[:])

	_, res, abort := handler(ctx, spendMsg, false)
	require.False(t, res.IsOK(), "tx OK with invalid parent confirm sig")
	require.True(t, abort, "tx with invalid parent confirm sig did not abort")

}

func TestAnteValidTx(t *testing.T) {
	// setup
	ctx, utxoStore, plasmaStore, presenceClaimStore := setup()
	handler := NewAnteHandler(utxoStore, plasmaStore, presenceClaimStore, conn{})

	// place input in store
	inputs := []inputUTXO{
		{nil, 0, 0, utils.Big1, addr, false},
		{nil, 0, 0, utils.Big2, addr, true},
	}
	setupInputs(ctx, utxoStore, inputs...)

	parentTx := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big2), [65]byte{}, nil),
			Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, nil), [65]byte{}, nil),
			Output0: plasma.NewOutput(addr, big.NewInt(10)),
			Output1: plasma.NewOutput(common.Address{}, nil),
			Fee:     utils.Big0,
		},
	}

	// set regular transaction utxo in store
	// parent input was 0.0.0.2
	// must create input key and confirmation hash
	// also need confirm sig of parent in order to spend
	inputKey := store.GetUTXOStoreKey(addr, plasma.NewPosition(nil, 0, 0, utils.Big2))
	confBytes := sha256.Sum256(append(parentTx.MerkleHash(), ctx.BlockHeader().DataHash...))
	confHash := utils.ToEthSignedMessageHash(confBytes[:])
	confSig, _ := crypto.Sign(confHash, privKey)
	inputUTXO := store.UTXO{
		InputKeys:        [][]byte{inputKey},
		ConfirmationHash: confBytes[:],
		Output: plasma.Output{
			Owner:  addr,
			Amount: big.NewInt(10),
		},
		Position: plasma.NewPosition(utils.Big1, 0, 0, nil),
	}
	utxoStore.StoreUTXO(ctx, inputUTXO)

	// store confirm sig into correct format
	var confirmSig [65]byte
	copy(confirmSig[:], confSig)

	// create msg
	spendMsg := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Input0:  plasma.NewInput(plasma.NewPosition(utils.Big1, 0, 0, nil), [65]byte{}, [][65]byte{confirmSig}),
			Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil),
			Output0: plasma.NewOutput(addr, big.NewInt(10)),
			Output1: plasma.NewOutput(addr, big.NewInt(9)),
			Fee:     utils.Big1,
		},
	}

	// set signature
	txHash := utils.ToEthSignedMessageHash(spendMsg.TxHash())
	sig, _ := crypto.Sign(txHash, privKey)
	copy(spendMsg.Input0.Signature[:], sig[:])
	copy(spendMsg.Input1.Signature[:], sig[:])

	_, res, abort := handler(ctx, spendMsg, false)
	require.True(t, res.IsOK(), "Valid tx does not have OK result")
	require.False(t, abort, "Valid tx aborted")

}

/*=====================================================================================================================================*/
// Deposit Antehandler tests

func TestAnteDeposit(t *testing.T) {
	// setup
	ctx, utxoStore, plasmaStore, presenceClaimStore := setup()
	handler := NewAnteHandler(utxoStore, plasmaStore, presenceClaimStore, conn{})

	// place input in store
	inputs := []inputUTXO{
		{nil, 0, 0, utils.Big1, addr, false},
		{nil, 0, 0, utils.Big2, addr, true},
	}
	setupInputs(ctx, utxoStore, inputs...)

	msg := msgs.IncludeDepositMsg{
		DepositNonce: big.NewInt(3),
		Owner:        addr,
	}

	_, res, abort := handler(ctx, msg, false)

	require.True(t, res.IsOK(), "Valid IncludeDepositMsg has erroneous result")
	require.False(t, abort, "Valid IncludeDepositMsg aborted")

	// try to include Deposit that already exists
	msg.DepositNonce = big.NewInt(1)

	_, res, abort = handler(ctx, msg, false)

	require.False(t, res.IsOK(), "Allowed to re-include deposit")
	require.True(t, abort, "Redundant IncludeDepositMsg did not abort")
}

type unfinalConn struct{}

func (u unfinalConn) GetDeposit(nonce *big.Int) (plasma.Deposit, *big.Int, bool) {
	dep := plasma.Deposit{
		Owner:       addr,
		Amount:      big.NewInt(10),
		EthBlockNum: big.NewInt(50),
	}
	return dep, big.NewInt(10), false
}

func (u unfinalConn) HasTxBeenExited(pos plasma.Position) bool { return false }

type dneConn struct{}

func (d dneConn) GetDeposit(nonce *big.Int) (plasma.Deposit, *big.Int, bool) {
	return plasma.Deposit{}, nil, false
}

func (d dneConn) HasTxBeenExited(pos plasma.Position) bool { return false }

func TestAnteDepositUnfinal(t *testing.T) {
	// setup
	ctx, utxoStore, plasmaStore, presenceClaimStore := setup()
	// connection always returns unfinalized deposits
	handler := NewAnteHandler(utxoStore, plasmaStore, presenceClaimStore, unfinalConn{})

	msg := msgs.IncludeDepositMsg{
		DepositNonce: big.NewInt(3),
		Owner:        addr,
	}

	_, res, abort := handler(ctx, msg, false)

	require.False(t, res.IsOK(), "Unfinalized deposit inclusion did not error")
	require.True(t, abort, "Unfinalized deposit inclusion did not abort")

}

func TestAnteDepositExitted(t *testing.T) {
	// setup
	ctx, utxoStore, plasmaStore, presenceClaimStore := setup()
	// connection always returns exitted deposits
	handler := NewAnteHandler(utxoStore, plasmaStore, presenceClaimStore, exitConn{})

	msg := msgs.IncludeDepositMsg{
		DepositNonce: big.NewInt(3),
		Owner:        addr,
	}

	_, res, abort := handler(ctx, msg, false)

	require.False(t, res.IsOK(), "Exitted deposit inclusion did not error")
	require.True(t, abort, "Exitted deposit inclusion did not abort")

}

func TestAnteDepositDNE(t *testing.T) {
	// setup
	ctx, utxoStore, plasmaStore, presenceClaimStore := setup()
	// connection always returns exitted deposits
	handler := NewAnteHandler(utxoStore, plasmaStore, presenceClaimStore, dneConn{})

	msg := msgs.IncludeDepositMsg{
		DepositNonce: big.NewInt(3),
		Owner:        addr,
	}

	_, res, abort := handler(ctx, msg, false)

	require.False(t, res.IsOK(), "Nonexistent deposit inclusion did not error")
	require.True(t, abort, "Nonexistent deposit inclusion did not abort")

}

func setupInputs(ctx sdk.Context, utxoStore store.UTXOStore, inputs ...inputUTXO) {
	for _, i := range inputs {
		utxo := store.UTXO{
			Output: plasma.Output{
				Owner:  i.Address,
				Amount: big.NewInt(10),
			},
			Spent:    i.Spent,
			Position: plasma.NewPosition(i.BlockNum, i.TxIndex, i.OIndex, i.DepositNonce),
		}
		utxoStore.StoreUTXO(ctx, utxo)
	}
}
