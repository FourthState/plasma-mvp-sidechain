package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

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

func TestAnteChecks(t *testing.T) {
	// setup
	ctx, utxoStore, plasmaStore := setup()
	handler := NewAnteHandler(utxoStore, plasmaStore, conn{})

	// bad keys to check against the deposit
	badPrivKey, _ := crypto.GenerateKey()

	type validationCase struct {
		reason string
		msgs.SpendMsg
	}

	// cases to check for. cases with signature checks will get set subsequent to this step
	// array of pointers because we are setting signatures after using `range`
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
					Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, big.NewInt(2)), [65]byte{}, nil),
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
					Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, big.NewInt(2)), [65]byte{}, nil),
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
					Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, big.NewInt(2)), [65]byte{}, nil),
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
					Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, big.NewInt(2)), [65]byte{}, nil),
					Output0: plasma.NewOutput(addr, big.NewInt(10)),
					Output1: plasma.NewOutput(addr, big.NewInt(10)),
					Fee:     utils.Big1,
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

// TODO: set handlers with position inputs
