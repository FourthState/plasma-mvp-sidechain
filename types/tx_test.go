package types

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"

	"github.com/FourthState/plasma-mvp-sidechain/utils"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	capKey := sdk.NewKVStoreKey("capkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, capKey
}

func GenBasicSpendMsg() SpendMsg {
	// Creates Basic Spend Msg with no owners or recipients
	var confirmSigs [2][65]byte
	return SpendMsg{
		Blknum0:      1,
		Txindex0:     0,
		Oindex0:      0,
		DepositNum0:  0,
		Owner0:       common.Address{},
		ConfirmSigs0: confirmSigs,
		Blknum1:      1,
		Txindex1:     1,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       common.Address{},
		ConfirmSigs1: confirmSigs,
		Newowner0:    common.Address{},
		Amount0:      150,
		Newowner1:    common.Address{},
		Amount1:      50,
		FeeAmount:    0,
	}
}

func GenSpendMsgWithAddresses() SpendMsg {
	// Creates Basic Spend Msg with owners and recipients
	var confirmSigs [2][65]byte
	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()

	return SpendMsg{
		Blknum0:      1,
		Txindex0:     0,
		Oindex0:      0,
		DepositNum0:  0,
		Owner0:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs0: confirmSigs,
		Blknum1:      1,
		Txindex1:     1,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs1: confirmSigs,
		Newowner0:    utils.PrivKeyToAddress(privKeyB),
		Amount0:      150,
		Newowner1:    utils.PrivKeyToAddress(privKeyB),
		Amount1:      50,
		FeeAmount:    0,
	}
}

// Creates a transaction with no owners
func TestNoOwners(t *testing.T) {
	var msg = GenBasicSpendMsg()
	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(201),
		err.Code(), err.Error())
}

// Creates a transaction with no recipients
func TestNoRecipients(t *testing.T) {
	privKeyA, _ := ethcrypto.GenerateKey()
	var msg = GenBasicSpendMsg()
	msg.Owner0 = utils.PrivKeyToAddress(privKeyA)
	msg.Owner1 = utils.PrivKeyToAddress(privKeyA)
	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(201),
		err.Code(), err.Error())
}

// The oindex is neither 0 or 1
func TestIncorrectOIndex(t *testing.T) {
	var msg = GenSpendMsgWithAddresses()
	msg.Oindex0 = 2

	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(202),
		err.Code(), err.Error())

	msg.Oindex0 = 0
	msg.Oindex1 = 2
	err = msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(202),
		err.Code(), err.Error())

}

// Creates an invalid transaction referencing utxo and deposit
func TestInvalidSpendDeposit(t *testing.T) {
	var msg = GenSpendMsgWithAddresses()
	msg.DepositNum0 = 5

	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(204), err.Code(), err.Error())

	msg.DepositNum0 = 0
	msg.DepositNum1 = 123
	err = msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(204), err.Code(), err.Error())
}

// Creates an invalid transaction spending same position twice
func TestInvalidPosition(t *testing.T) {
	var msg = GenSpendMsgWithAddresses()
	// Set second position equal to first position
	msg.Txindex1 = 0

	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(204), err.Code(), err.Error())
}

// Try to spend with 0 denomination for first output
func TestInvalidDenomination(t *testing.T) {
	var msg = GenSpendMsgWithAddresses()
	msg.Amount0 = 0

	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(203), err.Code(), err.Error())
}

// Tests GetSigners method
func TestGetSigners(t *testing.T) {
	msg := GenBasicSpendMsg()
	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()

	msg.Owner0 = utils.PrivKeyToAddress(privKeyA)
	addrs := []sdk.AccAddress{sdk.AccAddress(msg.Owner0.Bytes())}
	signers := msg.GetSigners() // GetSigners() returns []sdk.AccAddress by interface constraint
	require.Equal(t, addrs, signers, "signer Address do not match")

	msg.Owner1 = utils.PrivKeyToAddress(privKeyB)
	addrs = []sdk.AccAddress{sdk.AccAddress(msg.Owner0.Bytes()), sdk.AccAddress(msg.Owner1.Bytes())}
	signers = msg.GetSigners()
	require.Equal(t, addrs, signers, "signer Addresses do not match")
}
