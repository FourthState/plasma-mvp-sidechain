package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// return a base utxo with nothing set, along with two addresses
func GetBareUTXO() (utxo *BaseUTXO, addrA, addrB common.Address) {
	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()
	addrA = utils.PrivKeyToAddress(privKeyA)
	addrB = utils.PrivKeyToAddress(privKeyB)
	return &BaseUTXO{}, addrA, addrB
}

// Basic tests checking methods for BaseUTXO
func TestGetSetAddress(t *testing.T) {
	utxo, addrA, addrB := GetBareUTXO()

	// try to set address to another blank address
	err := utxo.SetAddress(common.Address{}.Bytes())
	require.Error(t, err)

	// set address to addrB
	err = utxo.SetAddress(addrB.Bytes())
	require.NoError(t, err)

	// try to set address to addrA (currently set to addrB)
	err = utxo.SetAddress(addrA.Bytes())
	require.Error(t, err)

	// check get method
	addr := utxo.GetAddress()
	require.Equal(t, addr, addrB.Bytes(), fmt.Sprintf("BaseUTXO GetAddress() method returned the wrong address: %s", addr))
}

// Test GetInputAddresses() and SetInputAddresses
func TestInputAddresses(t *testing.T) {
	utxo, addrA, addrB := GetBareUTXO()

	// try to set input address to blank addresses
	err := utxo.SetInputAddresses([2]common.Address{common.Address{}, common.Address{}})
	require.Error(t, err)

	// set input addresses to addrA, addrA
	err = utxo.SetInputAddresses([2]common.Address{addrA, addrA})
	require.NoError(t, err)

	// try to set input address to addrB
	err = utxo.SetInputAddresses([2]common.Address{addrB, common.Address{}})
	require.Error(t, err)

	// check get method
	addrs := utxo.GetInputAddresses()
	require.Equal(t, addrs, [2]common.Address{addrA, addrA})
}

// Test GetAmount() and SetAmount()
func TestAmount(t *testing.T) {
	utxo := NewBaseUTXO(common.Address{}, [2]common.Address{common.Address{}, common.Address{}}, 100, "ether", PlasmaPosition{})

	// try to set denom when it already has a value
	err := utxo.SetAmount(100000000)
	require.Error(t, err)

	// check get method
	amount := utxo.GetAmount()
	require.Equal(t, amount, uint64(100), "the wrong amount was returned by GetAmount()")
}

// Test GetPosition() and SetPosition()
func TestPosition(t *testing.T) {
	utxo := BaseUTXO{}
	position := NewPlasmaPosition(0, uint16(0), uint8(0), 0)

	// try to set position to incorrect position 0, 0, 0, 0
	err := utxo.SetPosition(position)
	require.Error(t, err)

	// set position
	position = NewPlasmaPosition(5, uint16(12), uint8(1), 0)
	err = utxo.SetPosition(position)
	require.NoError(t, err)

	// try to set to different position
	position = NewPlasmaPosition(1, uint16(23), uint8(1), 0)
	err = utxo.SetPosition(position)
	require.Error(t, err)

	// check get method
	require.Equal(t, utxo.GetPosition(), NewPlasmaPosition(5, uint16(12), uint8(1), 0), "the wrong position was returned")
}
