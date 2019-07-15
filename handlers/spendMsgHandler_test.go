package handlers

import (
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

var nextTxIndex = func() uint16 { return 0 }
var feeUpdater = func(num *big.Int) sdk.Error { return nil }

func TestSpend(t *testing.T) {
	// plasmaStore is at next block height 1
	ctx, utxoStore, plasmaStore, _ := setup()
	privKey, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(privKey.PublicKey)

	spendHandler := NewSpendHandler(utxoStore, plasmaStore, nextTxIndex, feeUpdater)

	// store inputs to be spent
	pos := plasma.NewPosition(utils.Big0, 0, 0, utils.Big1)
	utxo := store.UTXO{
		Output: plasma.NewOutput(addr, big.NewInt(20)),
		// position (0,0,0,1)
		Position: pos,
		Spent:    false,
	}
	utxoStore.StoreUTXO(ctx, utxo)

	// create a msg that spends the first input and creates two outputs
	newOwner := common.HexToAddress("1")
	msg := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Input0:  plasma.NewInput(pos, [65]byte{}, nil),
			Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, nil), [65]byte{}, nil),
			Output0: plasma.NewOutput(newOwner, big.NewInt(10)),
			Output1: plasma.NewOutput(newOwner, big.NewInt(10)),
			Fee:     utils.Big0,
		},
	}
	// fill in the signature
	sig, err := crypto.Sign(utils.ToEthSignedMessageHash(msg.TxHash()), privKey)
	copy(msg.Input0.Signature[:], sig)
	err = msg.ValidateBasic()
	require.NoError(t, err)

	res := spendHandler(ctx, msg)
	require.Truef(t, res.IsOK(), "failed to handle spend: %s", res)

	// check that the utxo store reflects the spend
	utxo, ok := utxoStore.GetUTXO(ctx, addr, pos)
	require.True(t, ok, "input to the msg does not exist in the store")
	require.True(t, utxo.Spent, "input not marked as spent after the handler")

	// new first output was created at BlockHeight 1 and txIndex 0 and outputIndex 0
	pos = plasma.NewPosition(utils.Big1, 0, 0, nil)
	utxo, ok = utxoStore.GetUTXO(ctx, newOwner, pos)
	require.True(t, ok, "new output was not created")
	require.False(t, utxo.Spent, "new output marked as spent")
	require.Equal(t, utxo.Output.Amount, big.NewInt(10), "new output has incorrect amount")

	// new second output was created at BlockHeight 0 and txIndex 0 and outputIndex 1
	pos = plasma.NewPosition(utils.Big1, 0, 1, nil)
	utxo, ok = utxoStore.GetUTXO(ctx, newOwner, pos)
	require.True(t, ok, "new output was not created")
	require.False(t, utxo.Spent, "new output marked as spent")
	require.Equal(t, utxo.Output.Amount, big.NewInt(10), "new output has incorrect amount")
}
