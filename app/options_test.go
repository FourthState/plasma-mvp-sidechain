package app

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	privkey = "713bd18559e878e0fa3ee32c8ff3ef4393b82ff9f272a3d7de707882f9a3f7d7"
)

func TestSetPrivKey(t *testing.T) {
	// create a private key file
	privkey_file, err := ioutil.TempFile("", "private_key")
	require.NoError(t, err)

	defer os.Remove(privkey_file.Name())

	n, err := privkey_file.Write([]byte(privkey))
	require.NoError(t, err)
	require.Equal(t, n, len(privkey))

	db := dbm.NewMemDB()
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
	cc := NewChildChain(logger, db, nil,
		SetEthPrivKey(privkey_file.Name(), true),
	)

	private_key, _ := crypto.LoadECDSA(privkey_file.Name())
	privkey_file.Close()
	require.EqualValues(t, private_key, cc.validatorPrivKey)
	require.Equal(t, true, cc.isValidator)
}
