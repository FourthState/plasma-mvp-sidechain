package store

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"testing"
)

func TestSavSig(t *testing.T) {
	// setup testing env
	os.Mkdir("testing", os.ModePerm)
	viper.Set(DirFlag, os.ExpandEnv("./testing"))

	// cleanup
	defer func() {
		viper.Reset()
		os.RemoveAll("testing")

	}()

	cases := [][]int64{
		// must save for both output indicies
		{5, 0, 0, 0},
		// deposit
		{0, 0, 0, 10},
		// different txindex
		{5, 1, 0, 0},
		// different blk num
		{6, 1, 0, 0},
	}

	for i, p := range cases {
		key, _ := crypto.GenerateKey()
		txHash := crypto.Keccak256([]byte("txhash"))

		expected, err := crypto.Sign(txHash, key)
		pos := plasma.NewPosition(big.NewInt(p[0]), uint16(p[1]), uint8(p[2]), big.NewInt(p[3]))

		_, err = GetSig(pos)
		require.Errorf(t, err, "case %d: did not error when getting non existent signature for position %s", i, pos)

		err = SaveSig(pos, expected, true)
		require.NoError(t, err, "case %d: failed to save signature for position %s", i, pos)

		actual, err := GetSig(pos)
		require.NoError(t, err, "case %d: failed when getting signature for position %s", i, pos)
		require.Equal(t, expected, actual, "case %d: actual signature was not equal to expected signature for position %s", i, pos)

		if !pos.IsDeposit() {
			// changing output index should not effect
			// retrieval of signature
			pos.OutputIndex = uint8(1)
			actual, err = GetSig(pos)

			require.NoError(t, err, "case %d: failed when getting signature for position %s", i, pos)
			require.Equal(t, expected, actual, "case %d: actual signature was not equal to expected signature for position %s", i, pos)
		}
	}
}

// pass in a signature without length 65 bytes
func TestBadSigs(t *testing.T) {
	// setup testing env
	os.Mkdir("testing", os.ModePerm)
	viper.Set(DirFlag, os.ExpandEnv("./testing"))

	// cleanup
	defer func() {
		viper.Reset()
		os.RemoveAll("testing")

	}()

	key, _ := crypto.GenerateKey()
	txHash := crypto.Keccak256([]byte("transaction hash"))

	sig, _ := crypto.Sign(txHash, key)
	pos := plasma.NewPosition(big.NewInt(10), uint16(5), uint8(0), big.NewInt(0))

	err := SaveSig(pos, sig[:60], true)
	require.Error(t, err, "did not reject a signature with a length not equal to 65")
}

// Save two confirmation signatures
func TestMultiConfirmSig(t *testing.T) {
	// setup testing env
	os.Mkdir("testing", os.ModePerm)
	viper.Set(DirFlag, os.ExpandEnv("./testing"))

	// cleanup
	defer func() {
		viper.Reset()
		os.RemoveAll("testing")

	}()

	key, _ := crypto.GenerateKey()
	txHash0 := crypto.Keccak256([]byte("hash one"))
	txHash1 := crypto.Keccak256([]byte("the second hash"))

	sig0, _ := crypto.Sign(txHash0, key)
	sig1, _ := crypto.Sign(txHash1, key)

	pos := plasma.NewPosition(big.NewInt(1000), uint16(256), uint8(0), big.NewInt(0))

	// save sig1 first
	err := SaveSig(pos, sig1, false)
	require.NoError(t, err, "failed to save first confirm signature")

	// prepend sig0
	err = SaveSig(pos, sig0, true)
	require.NoError(t, err, "failed to save second confirm signature")

	sigs, err := GetSig(pos)
	require.NoError(t, err, "failed to retrieve confirm signatures")
	require.Equal(t, append(sig0, sig1...), sigs, "retrieved signatures do not match expected signatures")
}
