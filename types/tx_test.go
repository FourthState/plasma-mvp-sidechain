package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func GenBasicSpendMsg() SpendMsg {
	// Creates Basic Spend Msg with no owners or recipients
	confirmSigs := [2]crypto.Signature{crypto.SignatureSecp256k1{}, crypto.SignatureSecp256k1{}}
	return SpendMsg{
		Blknum1:      1000,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       crypto.Address([]byte("")),
		ConfirmSigs1: confirmSigs,
		Blknum2:      1000,
		Txindex2:     1,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       crypto.Address([]byte("")),
		ConfirmSigs2: confirmSigs,
		Newowner1:    crypto.Address([]byte("")),
		Denom1:       150,
		Newowner2:    crypto.Address([]byte("")),
		Denom2:       50,
		Fee:          0,
	}
}

func GenSpendMsgWithAddresses() SpendMsg {
	// Creates Basic Spend Msg with owners and recipients
	confirmSigs := [2]crypto.Signature{crypto.SignatureSecp256k1{}, crypto.SignatureSecp256k1{}}
	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()
	
	return SpendMsg{
		Blknum1:      1000,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       EthPrivKeyToSDKAddress(privKeyA),
		ConfirmSigs1: confirmSigs,
		Blknum2:      1000,
		Txindex2:     1,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       EthPrivKeyToSDKAddress(privKeyA),
		ConfirmSigs2: confirmSigs,
		Newowner1:    EthPrivKeyToSDKAddress(privKeyB),
		Denom1:       150,
		Newowner2:    EthPrivKeyToSDKAddress(privKeyB),
		Denom2:       50,
		Fee:          0,
	}
}

func TestNoOwners(t *testing.T) {
	var msg = GenBasicSpendMsg()
	err := msg.ValidateBasic()
	assert.Equal(t, sdk.CodeType(101),
		err.Code(), err.Error())
}

func TestNoRecipients(t *testing.T) {
	privKeyA := crypto.GenerateKey()
	var msg = GenBasicSpendMsg()
	msg.Owner1 = EthPrivKeyToSDKAddress(privKeyA)
	msg.Owner2 = EthPrivKeyToSDKAddress(privKeyA)
	err := msg.ValidateBasic()
	assert.Equal(t, sdk.CodeType(101),
		err.Code(), err.Error())
}

func TestIncorrectOIndex(t *testing.T) {
	var msg1 = GenSpendMsgWithAddresses()
	msg1.Oindex1 = 2
	var msg2 = GenSpendMsgWithAddresses()
	msg2.Oindex2 = 2

	err1 := msg1.ValidateBasic()
	assert.Equal(t, sdk.CodeType(102),
		err1.Code(), err1.Error())

	err2 := msg2.ValidateBasic()
	assert.Equal(t, sdk.CodeType(102),
		err2.Code(), err2.Error())

}

func TestInvalidSpendDeposit(t *testing.T) {
	var msg1 = GenSpendMsgWithAddresses()
	msg1.DepositNum1 = 5

	err := msg1.ValidateBasic()
	assert.Equal(t, sdk.CodeType(106), err.Code(), err.Error())
}

func TestGetSigners(t *testing.T) {
	msg := GenBasicSpendMsg()
	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()

	msg.Owner1 = EthPrivKeyToSDKAddress(privKeyA)
	addrs := []crypto.Address{msg.Owner1}
	signers := msg.GetSigners()
	assert.Equal(t, addrs, signers, "Signer Address do not match")

	msg.Owner2 = EthPrivKeyToSDKAddress(privKeyB)
	addrs = []crypto.Address{msg.Owner1, msg.Owner2}
	signers = msg.GetSigners()
	assert.Equal(t, addrs, signers, "Signer Addresses do not match")
}
