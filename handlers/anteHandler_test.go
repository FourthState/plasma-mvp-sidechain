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

type Tx struct {
	Transaction plasma.Transaction
	ConfirmationHash []byte
	Spent        []bool
	Spenders [][]byte
	Position plasma.Position
}

type Deposit struct {
	Owner common.Address
	Nonce *big.Int
	EthBlockNum *big.Int
	Amount *big.Int
	Spent	bool
	Spender	[]byte
}

// cook the plasma connection
type conn struct{}

// all deposits should be in an amount of 10eth owner by addr(defined above)
func (p conn) GetDeposit(tmBlock *big.Int, nonce *big.Int) (plasma.Deposit, *big.Int, bool) {
	return plasma.Deposit{addr, big.NewInt(10), utils.Big0}, big.NewInt(-2), true
}
func (p conn) HasTxBeenExited(tmBlock *big.Int, pos plasma.Position) bool { return false }

var _ plasmaConn = conn{}

// cook up different plasma connection that will always claim Inputs exitted
type exitConn struct{}

// all deposits should be in an amount of 10eth owner by addr(defined above)
func (p exitConn) GetDeposit(tmBlock *big.Int, nonce *big.Int) (plasma.Deposit, *big.Int, bool) {
	return plasma.Deposit{addr, big.NewInt(10), utils.Big0}, big.NewInt(-2), true
}
func (p exitConn) HasTxBeenExited(tmBlock *big.Int, pos plasma.Position) bool { return true }

func TestAnteChecks(t *testing.T) {
	// setup
	ctx, txStore, depositStore, blockStore := setup()
	handler := NewAnteHandler(txStore, depositStore, blockStore, conn{})

	// cook up some input deposits
	inputs := []Deposit{
		{addr, big.NewInt(1), big.NewInt(100), big.NewInt(10), false, []byte{}},
		{addr, big.NewInt(2), big.NewInt(101), big.NewInt(10), false, []byte{}},
		{addr, big.NewInt(3), big.NewInt(102), big.NewInt(10), false, []byte{}},
	}
	setupDeposits(ctx, depositStore, inputs...)

	type validationCase struct {
		reason string
		msgs.SpendMsg
	}

	// cases to check for. cases with signature checks will get set subsequent to this step
	// array of pointers because we are setting signatures after using `range`
	// since InputsKeys not set, confirm signatures will simply be 0 bytes
	invalidCases := []*validationCase{
		&validationCase{
			reason: "incorrect first signature",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil)},
					Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(10))},
					Fee:     utils.Big0,
				},
			},
		},
		&validationCase{
			reason: "incorrect second signature",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil), plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big2), [65]byte{}, nil)},
					Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(20))},
					Fee:     utils.Big0,
				},
			},
		},
		&validationCase{
			reason: "no signatures",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil), plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big2), [65]byte{}, nil)},
					Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(20))},
					Fee:     utils.Big0,
				},
			},
		},
		&validationCase{
			reason: "invalid fee",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil), plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big2), [65]byte{}, nil)},
					Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(20))},
					Fee:     big.NewInt(20),
				},
			},
		},
		&validationCase{
			reason: "unbalanced transaction",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil), plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big2), [65]byte{}, nil)},
					Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(10)), plasma.NewOutput(addr, big.NewInt(10))},
					Fee:     utils.Big1,
				},
			},
		},
		&validationCase{
			reason: "Inputs deposit utxo does not exist",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil), plasma.NewInput(plasma.NewPosition(nil, 0, 0, big.NewInt(4)), [65]byte{}, nil)},
					Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(10)), plasma.NewOutput(addr, big.NewInt(10))},
					Fee:     utils.Big0,
				},
			},
		},
		&validationCase{
			reason: "Inputs transaction utxo does not exist",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil), plasma.NewInput(plasma.NewPosition(utils.Big1, 3, 1, nil), [65]byte{}, nil)},
					Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(10)), plasma.NewOutput(addr, big.NewInt(10))},
					Fee:     utils.Big0,
				},
			},
		},
		&validationCase{
			reason: "Inputs utxo already spent",
			SpendMsg: msgs.SpendMsg{
				Transaction: plasma.Transaction{
					Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil), plasma.NewInput(plasma.NewPosition(nil, 0, 0, big.NewInt(3)), [65]byte{}, nil)},
					Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(10)), plasma.NewOutput(addr, big.NewInt(10))},
					Fee:     utils.Big0,
				},
			},
		},
	}

	// set invalid first signature
	txHash := utils.ToEthSignedMessageHash(invalidCases[0].SpendMsg.TxHash())
	sig, _ := crypto.Sign(txHash, badPrivKey)
	copy(invalidCases[0].SpendMsg.Inputs[0].Signature[:], sig)

	// set invalid second signature but correct first signature
	txHash = utils.ToEthSignedMessageHash(invalidCases[1].SpendMsg.TxHash())
	sig, _ = crypto.Sign(txHash, badPrivKey)
	copy(invalidCases[1].SpendMsg.Inputs[1].Signature[:], sig)
	sig, _ = crypto.Sign(txHash, privKey)
	copy(invalidCases[1].SpendMsg.Inputs[0].Signature[:], sig)

	// set valid signatures for remaining cases
	for _, txCase := range invalidCases[3:] {
		txHash = utils.ToEthSignedMessageHash(txCase.SpendMsg.TxHash())
		sig, _ = crypto.Sign(txHash, privKey)
		copy(txCase.SpendMsg.Inputs[0].Signature[:], sig[:])
		copy(txCase.SpendMsg.Inputs[1].Signature[:], sig[:])
	}

	for _, txCase := range invalidCases {
		_, res, abort := handler(ctx, txCase.SpendMsg, false)
		require.False(t, res.IsOK(), txCase.reason)
		require.True(t, abort, txCase.reason)
	}
}

func TestAnteExitedInputs(t *testing.T) {
	// setup
	ctx, txStore, depositStore, blockStore := setup()
	handler := NewAnteHandler(txStore, depositStore, blockStore, exitConn{})

	// place Inputs in store
	inputs := Tx{
		Transaction:     ,
		Address:      addr,
		ConfirmationHash: ,
		Spent: []bool{false},
		Spenders: [][]byte{},
		Position: []plasma.Position[]
	}
	setupTxs(ctx, txStore, inputs)

	// create msg
	spendMsg := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(utils.Big1, 0, 0, nil), [65]byte{}, nil), plasma.NewInput(plasma.NewPosition(nil, 0, 0, nil), [65]byte{}, nil)},
			Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(10)), plasma.NewOutput(addr, big.NewInt(9))},
			Fee:     utils.Big1,
		},
	}

	// set signature
	txHash := utils.ToEthSignedMessageHash(spendMsg.TxHash())
	sig, _ := crypto.Sign(txHash, privKey)
	copy(spendMsg.Inputs[0].Signature[:], sig[:])

	_, res, abort := handler(ctx, spendMsg, false)
	require.False(t, res.IsOK(), "Result OK even though inputs exitted")
	require.True(t, abort, "Did not abort tx even though inputs exitted")

	// TODO: test case where grandparent exitted but parent didn't
}

func TestAnteInvalidConfirmSig(t *testing.T) {
	// setup
	ctx, txStore, depositStore, blockStore := setup()
	handler := NewAnteHandler(txStore, depositStore, blockStore, conn{})

	// place inputs in store
	inputs := []Deposit{
		{addr, utils.Big1, big.NewInt(50), big.NewInt(10), false, []byte{}},
		{addr, utils.Big2, big.NewInt(55), big.NewInt(10), false, []byte{}},
	}
	setupDeposits(ctx, depositStore, inputs...)

	parentTx := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big2), [65]byte{}, nil)},
			Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(10))},
			Fee:     utils.Big0,
		},
	}

	// set regular transaction utxo in store
	// parent input was 0.0.0.2
	// must create confirmation hash
	// also need confirm sig of parent in order to spend
	InputsKey := store.GetUTXOStoreKey(addr, plasma.NewPosition(nil, 0, 0, utils.Big2))
	confBytes := sha256.Sum256(append(parentTx.MerkleHash(), ctx.BlockHeader().DataHash...))
	confHash := utils.ToEthSignedMessageHash(confBytes[:])
	badConfSig, _ := crypto.Sign(confHash, badPrivKey)
	inputUTXO := store.UTXO{
		InputsKeys:       [][]byte{InputsKey},
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
			Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(utils.Big1, 0, 0, nil), [65]byte{}, [][65]byte{invalidConfirmSig}), plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil)},
			Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(10)), plasma.NewOutput(addr, big.NewInt(9))},
			Fee:     utils.Big1,
		},
	}

	// set signature
	txHash := utils.ToEthSignedMessageHash(spendMsg.TxHash())
	sig, _ := crypto.Sign(txHash, privKey)
	copy(spendMsg.Inputs[0].Signature[:], sig[:])
	copy(spendMsg.Inputs[1].Signature[:], sig[:])

	_, res, abort := handler(ctx, spendMsg, false)
	require.False(t, res.IsOK(), "tx OK with invalid parent confirm sig")
	require.True(t, abort, "tx with invalid parent confirm sig did not abort")

}

func TestAnteValidTx(t *testing.T) {
	// setup
	ctx, txStore, depositStore, blockStore := setup()
	handler := NewAnteHandler(txStore, depositStore, blockStore, conn{})

	// place inputs in store
	inputs := []InputUTXO{
		{nil, 0, 0, utils.Big1, addr, false},
		{nil, 0, 0, utils.Big2, addr, true},
	}
	setupInputs(ctx, txStore, inputs...)

	parentTx := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big2), [65]byte{}, nil)},
			Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(10))},
			Fee:     utils.Big0,
		},
	}

	// set regular transaction utxo in store
	// parent input was 0.0.0.2
	// must create input key and confirmation hash
	// also need confirm sig of parent in order to spend
	InputKey := store.GetUTXOStoreKey(addr, plasma.NewPosition(nil, 0, 0, utils.Big2))
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
			Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(utils.Big1, 0, 0, nil), [65]byte{}, [][65]byte{confirmSig}), plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), [65]byte{}, nil)},
			Outputs: []plasma.Output{plasma.NewOutput(addr, big.NewInt(10)), plasma.NewOutput(addr, big.NewInt(9))},
			Fee:     utils.Big1,
		},
	}

	// set signature
	txHash := utils.ToEthSignedMessageHash(spendMsg.TxHash())
	sig, _ := crypto.Sign(txHash, privKey)
	copy(spendMsg.Inputs[0].Signature[:], sig[:])
	copy(spendMsg.Inputs[1].Signature[:], sig[:])

	_, res, abort := handler(ctx, spendMsg, false)
	require.True(t, res.IsOK(), "Valid tx does not have OK result")
	require.False(t, abort, "Valid tx aborted")

}

/*=====================================================================================================================================*/
// Deposit Antehandler tests

func TestAnteDeposit(t *testing.T) {
	// setup
	ctx, txStore, depositStore, blockStore := setup()
	handler := NewAnteHandler(txStore, depositStore, blockStore, conn{})

	// place input in store
	inputs := []InputUTXO{
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

func (u unfinalConn) GetDeposit(tmBlock *big.Int, nonce *big.Int) (plasma.Deposit, *big.Int, bool) {
	dep := plasma.Deposit{
		Owner:       addr,
		Amount:      big.NewInt(10),
		EthBlockNum: big.NewInt(50),
	}
	return dep, big.NewInt(10), false
}

func (u unfinalConn) HasTxBeenExited(tmBlock *big.Int, pos plasma.Position) bool { return false }

type dneConn struct{}

func (d dneConn) GetDeposit(tmBlock *big.Int, nonce *big.Int) (plasma.Deposit, *big.Int, bool) {
	return plasma.Deposit{}, nil, false
}

func (d dneConn) HasTxBeenExited(tmBlock *big.Int, pos plasma.Position) bool { return false }

func TestAnteDepositUnfinal(t *testing.T) {
	// setup
	ctx, txStore, depositStore, blockStore := setup()
	// connection always returns unfinalized deposits
	handler := NewAnteHandler(txStore, depositStore, blockStore, unfinalConn{})

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
	ctx, txStore, depositStore, blockStore := setup()
	// connection always returns exitted deposits
	handler := NewAnteHandler(utxoStore, depositStore, blockStore, exitConn{})

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
	ctx, txStore, depositStore, blockStore := setup()
	// connection always returns exitted deposits
	handler := NewAnteHandler(txStore, depositStore, blockStore, dneConn{})

	msg := msgs.IncludeDepositMsg{
		DepositNonce: big.NewInt(3),
		Owner:        addr,
	}

	_, res, abort := handler(ctx, msg, false)

	require.False(t, res.IsOK(), "Nonexistent deposit inclusion did not error")
	require.True(t, abort, "Nonexistent deposit inclusion did not abort")

}

func setupDeposits(ctx sdk.Context, txStore store.TxStore, depositStore store.DepositStore, inputs ...InputUTXO) {
	for _, i := range inputs {
			deposit := store.Deposit{
				Deposit: plasma.Deposit{
					Owner: i.Owner,
					Amount: i.Amount,
					EthBlockNum: i.EthBlockNum
				},
				Spent: i.Spent,
				Spender: i.Spender,
			}
			depositStore.StoreDeposit(ctx, i.Nonce, deposit)
			txStore.StoreDepositWithAccount(ctx, i.Nonce, deposit)
	}
}

func setupTxs(ctx sdk.Context, txStore store.TxStore, inputs ...InputUTXO) {
	for _, i := range inputs {
		pos := plasma.NewPosition(i.BlockNum, i.TxIndex, i.OIndex, i.DepositNonce)
		if pos.IsDeposit() {
					} else {
			tx := store.Transaction{
				Output: plasma.Output{
					Owner:  i.Address,
					Amount: big.NewInt(10),
				},
				Spent:    i.Spent,
				Position: pos,
			}
		}
		
	}
}
