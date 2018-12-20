package app

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"

	"github.com/FourthState/plasma-mvp-sidechain/utils"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func TestSetEthConfig(t *testing.T) {
	rootchain := utils.GenerateAddress()

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
		SetEthConfig(true, privkey_file.Name(), rootchain.String(), nodeURL, "200", "16"),
	)

	private_key, _ := crypto.LoadECDSA(privkey_file.Name())
	privkey_file.Close()
	require.EqualValues(t, private_key, cc.validatorPrivKey)
	require.Equal(t, true, cc.isValidator)

	var empty ethcmn.Address
	require.NotEqual(t, empty, cc.rootchain)
	require.Equal(t, rootchain, cc.rootchain)

	require.Equal(t, uint64(200), cc.minFees)

	require.Equal(t, uint64(16), cc.blockFinality)
}
