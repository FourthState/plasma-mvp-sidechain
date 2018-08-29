package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// Basic tests checking methods for BaseUTXO
func TestUTXO(t *testing.T) {
	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()
	addrA := utils.PrivKeyToAddress(privKeyA)
	addrB := utils.PrivKeyToAddress(privKeyB)
	utxo := NewBaseUTXO(common.Address{}, [2]common.Address{common.Address{}, common.Address{}}, 100, NewPosition(0, 0, 0, 0))

	// try to set address to another blank address
	err := utxo.SetAddress(common.Address{})
	require.EqualError(t, err, "address provided is nil")

	// set address to addrB
	err = utxo.SetAddress(addrB)
	require.NoError(t, err)

	// try to set address to addrA (currently set to addrB)
	err = utxo.SetAddress(addrA)
	require.EqualError(t, err, "cannot override BaseUTXO Address")

	// check get method
	addr := utxo.GetAddress()
	require.Equal(t, addr, addrB, fmt.Sprintf("BaseUTXO GetAddress() method returned the wrong address: %s", addr))

	// try to set input address to blank addresses
	err = utxo.SetInputAddresses([2]common.Address{common.Address{}, common.Address{}})
	require.EqualError(t, err, "address provided is nil")

	// set input addresses to addrA, addrA
	err = utxo.SetInputAddresses([2]common.Address{addrA, addrA})
	require.NoError(t, err)

	// try to set input address to addrB
	err = utxo.SetInputAddresses([2]common.Address{addrB, common.Address{}})
	require.EqualError(t, err, "cannot override BaseUTXO Address")

	// check get method
	addrs := utxo.GetInputAddresses()
	require.Equal(t, addrs, [2]common.Address{addrA, addrA})

	// try to set denom when it already has a value
	err = utxo.SetDenom(100000000)
	require.EqualError(t, err, "cannot override BaseUTXO denomination")

	// check get method
	amount := utxo.GetDenom()
	require.Equal(t, amount, uint64(100), "the wrong denomination was returned by GetDenom()")

	// try to set position to incorrect position 0, 0, 0, 0
	err = utxo.SetPosition(0, uint16(0), uint8(0), 0)
	require.EqualError(t, err, "position cannot be set to 0, 0, 0, 0")

	// set position
	err = utxo.SetPosition(5, 12, 1, 0)
	require.NoError(t, err)

	// try to set to different position
	err = utxo.SetPosition(1, 23, 1, 0)
	require.EqualError(t, err, "cannot override BaseUTXO Position")

	// check get method
	position := utxo.GetPosition()
	require.Equal(t, position, NewPosition(5, 12, 1, 0))

}
