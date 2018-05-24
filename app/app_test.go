package app

import (
	"os"
	"testing"
	//"fmt" //for debugging
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	db "plasma-mvp-sidechain/db"
	utils "plasma-mvp-sidechain/utils"
	//rlp "github.com/ethereum/go-ethereum/rlp"
)

func newChildChain() *ChildChain {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return NewChildChain(logger, db)
}

func TestDepositMsg(t *testing.T) {
	cc := newChildChain()

	confirmSigs := [2]crypto.Signature{crypto.SignatureSecp256k1{}, crypto.SignatureSecp256k1{}}
	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()

	// Construct a SpendMsg
	var msg = db.SpendMsg{
		Blknum1:      0,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       utils.EthPrivKeyToSDKAddress(privKeyA),
		ConfirmSigs1: confirmSigs,
		Blknum2:      0,
		Txindex2:     0,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       crypto.Address([]byte("")),
		ConfirmSigs2: confirmSigs,
		Newowner1:    utils.EthPrivKeyToSDKAddress(privKeyB),
		Denom1:       1000,
		Newowner2:    crypto.Address([]byte("")),
		Denom2:       0,
		Fee:          1,
	}

	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ := ethcrypto.Sign(hash, privKeyA)
	tx := db.NewBaseTx(msg, []sdk.StdSignature{{
		PubKey:    nil,
		Signature: crypto.SignatureSecp256k1(sig),
	}})

	cdc := MakeCodec()
	txBytes, err := cdc.MarshalBinary(tx)

	require.NoError(t, err)

	// Run a check
	cres := cc.CheckTx(txBytes)
	assert.Equal(t, sdk.CodeType(6),
		sdk.CodeType(cres.Code), cres.Log)

	// Simulate a Block
	cc.BeginBlock(abci.RequestBeginBlock{})
	dres := cc.DeliverTx(txBytes)
	assert.Equal(t, sdk.CodeType(6), sdk.CodeType(dres.Code), dres.Log)

}
