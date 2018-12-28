package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

var feeUpdater FeeUpdater = func(amt *big.Int) sdk.Error {
	return nil
}

var (
	privKey, _ = crypto.GenerateKey()
	addr       = crypto.PubkeyToAddress(privKey.PublicKey)
)

// cook the plasma connection
type conn struct{}

// all deposits should be in an amount of 10eth owner by addr(defined above)
func (p conn) GetDeposit(nonce *big.Int) (plasma.Deposit, bool) {
	return plasma.Deposit{addr, big.NewInt(10), utils.Big0}, true
}
func (p conn) HasTxBeenExited(pos plasma.Position) bool { return false }

var _ plasmaConn = conn{}

func TestIncorrectFirstSignature(t *testing.T) {
	ctx, utxoStore, plasmaStore := setup()
	msg := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), addr, [65]byte{}, [][65]byte{}),
			Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, nil), common.Address{}, [65]byte{}, [][65]byte{}),
			Output0: plasma.NewOutput(addr, big.NewInt(10)),
			Output1: plasma.NewOutput(common.Address{}, utils.Big0),
			Fee:     utils.Big0,
		},
	}

	// Input0's signature signed by the wrong address
	badKey, _ := crypto.GenerateKey()
	txHash := utils.ToEthSignedMessageHash(msg.TxHash())
	var sig [65]byte
	s, _ := crypto.Sign(txHash[:], badKey)
	copy(sig[:], s)
	msg.Input0.Signature = sig

	handler := NewAnteHandler(utxoStore, plasmaStore, feeUpdater, conn{})
	_, res, abort := handler(ctx, msg, false)

	require.True(t, abort, "handler did not abort with no signatures")
	require.Equal(t, sdk.CodeUnauthorized, res.Code, "handler did not catch signature authorization error")
}

func TestIncorrectSecondSignature(t *testing.T) {
	ctx, utxoStore, plasmaStore := setup()
	// two deposits
	msg := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), addr, [65]byte{}, [][65]byte{}),
			Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, big.NewInt(1)), addr, [65]byte{}, [][65]byte{}),
			Output0: plasma.NewOutput(addr, big.NewInt(10)),
			Output1: plasma.NewOutput(common.Address{}, utils.Big0),
			Fee:     utils.Big0,
		},
	}

	txHash := utils.ToEthSignedMessageHash(msg.TxHash())

	// first signature will be correct but second will be incorrect
	var sig0 [65]byte
	s0, _ := crypto.Sign(txHash[:], privKey)
	copy(sig0[:], s0)
	msg.Input0.Signature = sig0

	// second signature will be corrupt
	var sig1 [65]byte
	badKey, _ := crypto.GenerateKey()
	s1, _ := crypto.Sign(txHash[:], badKey)
	copy(sig1[:], s1)
	msg.Input1.Signature = sig1

	handler := NewAnteHandler(utxoStore, plasmaStore, feeUpdater, conn{})
	_, res, abort := handler(ctx, msg, false)

	require.True(t, abort, "handler did not abort in the abscence of a second signature")
	require.Equal(t, sdk.CodeUnauthorized, res.Code, "handler did catch signature authorization error")
}

func TestInvalidFee(t *testing.T) {
	ctx, utxoStore, plasmaStore := setup()
	msg := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Input0:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), addr, [65]byte{}, [][65]byte{}),
			Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, nil), common.Address{}, [65]byte{}, [][65]byte{}),
			Output0: plasma.NewOutput(addr, big.NewInt(10)),
			Output1: plasma.NewOutput(common.Address{}, utils.Big0),
			Fee:     big.NewInt(11), // larger than the first input
		},
	}
	txHash := utils.ToEthSignedMessageHash(msg.TxHash())

	var sig0 [65]byte
	s0, _ := crypto.Sign(txHash[:], privKey)
	copy(sig0[:], s0)
	msg.Input0.Signature = sig0

	handler := NewAnteHandler(utxoStore, plasmaStore, feeUpdater, conn{})
	_, res, abort := handler(ctx, msg, false)

	require.True(t, abort, "handler did not catch a fee amount larger than the first input")
	require.Equal(t, CodeInsufficientFee, res.Code, "handler did not catch insufficient fee error")
}

func TestUnbalancedTransaction(t *testing.T) {
	ctx, utxoStore, plasmaStore := setup()
	msg := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			// 20eth of inputs
			Input0: plasma.NewInput(plasma.NewPosition(nil, 0, 0, utils.Big1), addr, [65]byte{}, [][65]byte{}),
			Input1: plasma.NewInput(plasma.NewPosition(nil, 0, 0, big.NewInt(2)), addr, [65]byte{}, [][65]byte{}),
			// 20eth of outputd
			Output0: plasma.NewOutput(addr, big.NewInt(20)),
			Output1: plasma.NewOutput(common.Address{}, utils.Big0),
			Fee:     big.NewInt(1), // creates an unbalanced equation
		},
	}
	txHash := utils.ToEthSignedMessageHash(msg.TxHash())

	var sig0 [65]byte
	s0, _ := crypto.Sign(txHash[:], privKey)
	copy(sig0[:], s0)
	msg.Input0.Signature = sig0
	msg.Input1.Signature = sig0

	handler := NewAnteHandler(utxoStore, plasmaStore, feeUpdater, conn{})
	_, res, abort := handler(ctx, msg, false)

	require.True(t, abort, "handler did not catch unbalanced inputs and outputs")
	require.Equal(t, msgs.CodeInvalidTransaction, res.Code, "handler did not catch invalid transaction due to unbalanced inputs and outputs")
}

// TODO: set handlers with position inputs
