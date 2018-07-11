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
	confirmSigs := [2]Signature{Signature{}, Signature{}}
	return SpendMsg{
		Blknum1:      1000,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       common.Address{},
		ConfirmSigs1: confirmSigs,
		Blknum2:      1000,
		Txindex2:     1,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       common.Address{},
		ConfirmSigs2: confirmSigs,
		Newowner1:    common.Address{},
		Denom1:       150,
		Newowner2:    common.Address{},
		Denom2:       50,
		Fee:          0,
	}
}

func GenSpendMsgWithAddresses() SpendMsg {
	// Creates Basic Spend Msg with owners and recipients
	confirmSigs := [2]Signature{Signature{}, Signature{}}
	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()

	return SpendMsg{
		Blknum1:      1000,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs1: confirmSigs,
		Blknum2:      1000,
		Txindex2:     1,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs2: confirmSigs,
		Newowner1:    utils.PrivKeyToAddress(privKeyB),
		Denom1:       150,
		Newowner2:    utils.PrivKeyToAddress(privKeyB),
		Denom2:       50,
		Fee:          0,
	}
}

// Creates a transaction with no owners
func TestNoOwners(t *testing.T) {
	var msg = GenBasicSpendMsg()
	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(101),
		err.Code(), err.Error())
}

// Creates a transaction with no recipients
func TestNoRecipients(t *testing.T) {
	privKeyA, _ := ethcrypto.GenerateKey()
	var msg = GenBasicSpendMsg()
	msg.Owner1 = utils.PrivKeyToAddress(privKeyA)
	msg.Owner2 = utils.PrivKeyToAddress(privKeyA)
	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(101),
		err.Code(), err.Error())
}

// The oindex is neither 0 or 1
func TestIncorrectOIndex(t *testing.T) {
	var msg1 = GenSpendMsgWithAddresses()
	msg1.Oindex1 = 2
	var msg2 = GenSpendMsgWithAddresses()
	msg2.Oindex2 = 2

	err1 := msg1.ValidateBasic()
	require.Equal(t, sdk.CodeType(102),
		err1.Code(), err1.Error())

	err2 := msg2.ValidateBasic()
	require.Equal(t, sdk.CodeType(102),
		err2.Code(), err2.Error())

}

// Creates an invalid transaction referencing utxo and deposit
func TestInvalidSpendDeposit(t *testing.T) {
	var msg1 = GenSpendMsgWithAddresses()
	msg1.DepositNum1 = 5

	err := msg1.ValidateBasic()
	require.Equal(t, sdk.CodeType(106), err.Code(), err.Error())
}

// Creates an invalid transaction spending same position twice
func TestInvalidPosition(t *testing.T) {
	var msg1 = GenSpendMsgWithAddresses()
	// Set second position equal to first position
	msg1.Txindex2 = 0

	err := msg1.ValidateBasic()
	require.Equal(t, sdk.CodeType(106), err.Code(), err.Error())
}

// Tests GetSigners method
func TestGetSigners(t *testing.T) {
	msg := GenBasicSpendMsg()
	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()

	msg.Owner1 = utils.PrivKeyToAddress(privKeyA)
	addrs := []sdk.AccAddress{sdk.AccAddress(msg.Owner1.Bytes())}
	signers := msg.GetSigners() // GetSigners() returns []sdk.AccAddress by interface constraint
	require.Equal(t, addrs, signers, "signer Address do not match")

	msg.Owner2 = utils.PrivKeyToAddress(privKeyB)
	addrs = []sdk.AccAddress{sdk.AccAddress(msg.Owner1.Bytes()), sdk.AccAddress(msg.Owner2.Bytes())}
	signers = msg.GetSigners()
	require.Equal(t, addrs, signers, "signer Addresses do not match")
}
