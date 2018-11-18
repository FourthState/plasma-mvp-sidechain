package app

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	secp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"os"
	"testing"

	"github.com/FourthState/plasma-mvp-sidechain/utils"
)

func TestGenesisState(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	app := NewChildChain(logger, db, nil)

	addrs := []common.Address{utils.GenerateAddress(), utils.GenerateAddress()}

	var genUTXOs []GenesisUTXO
	for i, addr := range addrs {
		genUTXOs = append(genUTXOs, NewGenesisUTXO(addr.Hex(), "100", [4]string{"0", "0", "0", fmt.Sprintf("%d", i+1)}))
	}

	pubKey := secp256k1.GenPrivKey().PubKey()
	valAddr := utils.GenerateAddress()

	genValidator := GenesisValidator{
		ConsPubKey: pubKey,
		Address:    valAddr.String(),
	}

	genState := GenesisState{
		Validator: genValidator,
		UTXOs:     genUTXOs,
	}

	appBytes, err := app.cdc.MarshalJSON(genState)
	assert.Nil(t, err)
	var genState2 GenesisState
	err = app.cdc.UnmarshalJSON(appBytes, &genState2)
	assert.Nil(t, err)

	assert.Equal(t, genState, genState2)

	res := app.InitChain(abci.RequestInitChain{AppStateBytes: appBytes})
	expected := abci.ResponseInitChain{
		Validators: []abci.ValidatorUpdate{abci.ValidatorUpdate{
			PubKey: tmtypes.TM2PB.PubKey(pubKey),
			Power:  1,
		}},
	}
	assert.Equal(t, expected, res)
}
