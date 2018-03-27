package app

import (
	"os"
	"testing"
	//"fmt" //for debugging
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
	crypto "github.com/tendermint/go-crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/abci/types"
	types "plasma-mvp-sidechain/types"
	sdk "github.com/cosmos/cosmos-sdk/types" 

)

func newChildChain() *ChildChain {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return NewChildChain(logger, db)
}

func TestSpendMsg(t *testing.T) {
	cc := newChildChain()
	
	// Construct a SpendMsg
	var msg = types.SpendMsg{
		Blknum1: 0,
		Txindex1: 0,
		Oindex1: 0,
		Owner1:crypto.Address([]byte("origin")),
		Blknum2: 0,
		Txindex2:0,
		Oindex2:0,
		Owner2: crypto.Address([]byte("")),
		Newowner1: crypto.Address([]byte("recipient")),
		Denom1: 1000,
		Newowner2:crypto.Address([]byte("")),
		Denom2:0,
		Fee: 1,
	}

	priv := crypto.GenPrivKeyEd25519()
	sig := priv.Sign(msg.GetSignBytes())
	tx := types.NewBaseTx(msg, []sdk.StdSignature {{
			PubKey: 	priv.PubKey(),
			Signature:	sig,
		}})
	//Change to RLP once implemented
	cdc := MakeCodec()
	txBytes, err:= cdc.MarshalBinary(tx)
	require.NoError(t, err)

	// Run a check 
	cres := cc.CheckTx(txBytes)
	assert.Equal(t, sdk.CodeType(101),
				sdk.CodeType(cres.Code), cres.Log)

	// Simulate a Block
	cc.BeginBlock(abci.RequestBeginBlock{})
	dres := cc.DeliverTx(txBytes)
	assert.Equal(t, sdk.CodeType(101), sdk.CodeType(dres.Code), dres.Log)

}