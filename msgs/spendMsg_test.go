package msgs

import (
	"crypto/ecdsa"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

// generate a sample spend with no owners
func genSpendMsg() SpendMsg {
	spendMsg := SpendMsg{}
	spendMsg.Input0 = plasma.NewInput(big.NewInt(1), 0, 0, big.NewInt(0), common.Address{}, [][65]byte{})
	spendMsg.Input1 = plasma.NewInput(big.NewInt(1), 1, 0, big.NewInt(0), common.Address{}, [][65]byte{})
	spendMsg.Output0 = plasma.NewOutput(common.Address{}, big.NewInt(150))
	spendMsg.Output1 = plasma.NewOutput(common.Address{}, big.NewInt(50))
	spendMsg.Fee = big.NewInt(0)

	return spendMsg
}

func generateKeys() (common.Address, *ecdsa.PrivateKey) {
	privKey, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(privKey.PublicKey)

	return addr, privKey
}

func TestNoOwners(t *testing.T) {
	msg := genSpendMsg()
	err := msg.ValidateBasic()

	require.Equal(t, sdk.CodeType(201), err.Code(), err.Error())
}

func TestNoRecipients(t *testing.T) {
	addr, _ := generateKeys()
	msg := genSpendMsg()
	msg.Input0.Owner = addr
	msg.Input1.Owner = addr

	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(201), err.Code(), err.Error())
}

func TestIncorrectOutputIndex(t *testing.T) {
	addr, _ := generateKeys()
	msg := genSpendMsg()
	msg.Input0.Owner = addr
	msg.Input0.OutputIndex = 2
	msg.Output0.Owner = addr

	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(202), err.Code(), err.Error())

	msg.Input0.OutputIndex = 0
	msg.Input1.OutputIndex = 2
	msg.Input1.Owner = addr
	err = msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(202), err.Code(), err.Error())
}

func TestInvalidDepositSpend(t *testing.T) {
	addr, _ := generateKeys()
	msg := genSpendMsg()
	msg.Output0.Owner = addr
	// utxo position has been specified
	msg.Input0.Owner = addr
	msg.Input0.DepositNonce = big.NewInt(10)

	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(204), err.Code(), err.Error())

	msg.Input0.DepositNonce = big.NewInt(0)
	msg.Input1.Owner = addr
	msg.Input1.DepositNonce = big.NewInt(10)
	err = msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(204), err.Code(), err.Error())
}

func TestInvalidPosition(t *testing.T) {
	addr, _ := generateKeys()
	msg := genSpendMsg()
	//set second position to be equal to the first
	msg.Input1.Owner = addr
	msg.Input1.TxIndex = 0
	msg.Output0.Owner = addr

	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(204), err.Code(), err.Error())
}

func TestInvalidDenomination(t *testing.T) {
	addr, _ := generateKeys()
	msg := genSpendMsg()
	msg.Input0.Owner = addr
	msg.Output0.Owner = addr
	msg.Output0.Amount = big.NewInt(0)

	err := msg.ValidateBasic()
	require.Equal(t, sdk.CodeType(203), err.Code(), err.Error())
}

func TestGetSigners(t *testing.T) {
	addr, _ := generateKeys()
	msg := genSpendMsg()
	msg.Input0.Owner = addr
	msg.Output0.Owner = addr

	addrs := []sdk.AccAddress{sdk.AccAddress(addr.Bytes())}
	signers := msg.GetSigners()
	require.Equal(t, addrs, signers, "signer addresses do not match")

	msg.Input1.Owner = addr
	addrs = append(addrs, sdk.AccAddress(addr.Bytes()))
	signers = msg.GetSigners()
	require.Equal(t, addrs, signers, "signer addresses do not match")
}
