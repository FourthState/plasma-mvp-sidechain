package app

import (
	"os"
	"fmt"
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	types "github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	rlp "github.com/ethereum/go-ethereum/rlp"
)

func newChildChain() *ChildChain {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return NewChildChain(logger, db)
}

func InitTestChain(addr common.Address, cc *ChildChain) {
	// Currently only initialize chain with one deposited UTXO
	genState := GenesisUTXO{
		Address: addr.Hex(),
		Denom: 100,
		Position: [4]uint64{0, 0, 0, 1},
	}
	genBytes, err := json.Marshal(genState)
	if err != nil {
		panic(err)
	}
	appStateBytes := []byte(fmt.Sprintf("{\"UTXOs\": [%s]}", string(genBytes)))

	initRequest := abci.RequestInitChain{AppStateBytes: appStateBytes}
	cc.InitChain(initRequest)
}

func GenerateSimpleMsg(Owner1, NewOwner1 common.Address, position [4]uint64, denom1 uint64, fee uint64) types.SpendMsg {
	confirmSigs := [2]types.Signature{types.Signature{}, types.Signature{}}
	return types.SpendMsg{
		Blknum1:      position[0],
		Txindex1:     uint16(position[1]),
		Oindex1:      uint8(position[2]),
		DepositNum1:  position[3],
		Owner1:       Owner1,
		ConfirmSigs1: confirmSigs,
		Blknum2:      0,
		Txindex2:     0,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       common.Address{},
		ConfirmSigs2: confirmSigs,
		Newowner1:    NewOwner1,
		Denom1:       denom1,
		Newowner2:    common.Address{},
		Denom2:       0,
		Fee:          fee,
	}
}

// Attempts to spend a non-existent utxo
// without depositing first.
func TestBadSpendMsg(t *testing.T) {
	cc := newChildChain()

	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()

	// Construct a SpendMsg
	msg := GenerateSimpleMsg(utils.PrivKeyToAddress(privKeyA), utils.PrivKeyToAddress(privKeyB),
	                        	[4]uint64{1, 0, 0, 0}, 1000, 1)

	// Signs the hash of the transaction
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ := ethcrypto.Sign(hash, privKeyA)
	tx := types.NewBaseTx(msg, []types.Signature{{
		Sig: crypto.SignatureSecp256k1(sig),
	}})

	txBytes, err := rlp.EncodeToBytes(tx)

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

func TestSpendDeposit(t *testing.T) {
	cc := newChildChain()

	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()
	addrA := utils.PrivKeyToAddress(privKeyA)
	addrB := utils.PrivKeyToAddress(privKeyB)

	InitTestChain(addrA, cc)

	msg := GenerateSimpleMsg(addrA, addrB, [4]uint64{0, 0, 0, 1}, 100, 0)

	// Set confirm signatures
	confirmBytes := types.NewPosition(0, 0, 0, 1).GetSignBytes()
	hash := ethcrypto.Keccak256(confirmBytes)
	confirmSig, _ := ethcrypto.Sign(hash, privKeyA)
	msg.ConfirmSigs1 = [2]types.Signature{types.Signature{confirmSig}, types.Signature{confirmSig}}

	// Signs the hash of the transaction
	hash = ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ := ethcrypto.Sign(hash, privKeyA)
	tx := types.NewBaseTx(msg, []types.Signature{{
		Sig: crypto.SignatureSecp256k1(sig),
	}})

	// Must commit for checkState to be set correctly. Should be fixed in next version of SDK
	cc.BeginBlock(abci.RequestBeginBlock{})
	cc.EndBlock(abci.RequestEndBlock{})
	cc.Commit()

	// Simulate a block
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})

	// Run a check
	cres := cc.Check(tx)
	assert.Equal(t, sdk.CodeType(0),
		sdk.CodeType(cres.Code), cres.Log)

	// Deliver tx, updates states
	dres := cc.Deliver(tx)

	assert.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

	// Create context
	ctx := cc.NewContext(false, abci.Header{})

	// Retrieve UTXO from context
	utxo := cc.utxoMapper.GetUTXO(ctx, types.NewPosition(1, 0, 0, 0))
	expected := types.NewBaseUTXO(addrB, [2]common.Address{addrA, common.Address{}}, 100, types.NewPosition(1, 0, 0, 0))

	assert.Equal(t, expected, utxo, "UTXO did not get added to store correctly")

}

func TestSpendTx(t *testing.T) {
	cc := newChildChain()

	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()
	addrA := utils.PrivKeyToAddress(privKeyA)
	addrB := utils.PrivKeyToAddress(privKeyB)

	InitTestChain(addrA, cc)
	cc.Commit()

	msg := GenerateSimpleMsg(addrA, addrB, [4]uint64{0, 0, 0, 1}, 100, 0)

	// Set confirm signatures
	confirmBytes := types.NewPosition(0, 0, 0, 1).GetSignBytes()
	hash := ethcrypto.Keccak256(confirmBytes)
	confirmSig, _ := ethcrypto.Sign(hash, privKeyA)
	msg.ConfirmSigs1 = [2]types.Signature{types.Signature{confirmSig}, types.Signature{confirmSig}}

	// Signs the hash of the transaction
	hash = ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ := ethcrypto.Sign(hash, privKeyA)
	tx := types.NewBaseTx(msg, []types.Signature{{
		Sig: crypto.SignatureSecp256k1(sig),
	}})

	// Simulate a block
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})

	// Deliver tx, updates states
	cc.Deliver(tx)

	cc.EndBlock(abci.RequestEndBlock{})
	cc.Commit()

	// Test that spending from a non-deposit/non-genesis UTXO works

	// generate simple msg
	msg = GenerateSimpleMsg(addrB, addrA, [4]uint64{1, 0, 0, 0}, 100, 0)

	// Set confirm signatures
	confirmBytes = types.NewPosition(1, 0, 0, 0).GetSignBytes()
	hash = ethcrypto.Keccak256(confirmBytes)
	confirmSig, _ = ethcrypto.Sign(hash, privKeyA)
	msg.ConfirmSigs1 = [2]types.Signature{types.Signature{confirmSig}, types.Signature{confirmSig}}

	// Signs the hash of the transaction
	hash = ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ = ethcrypto.Sign(hash, privKeyB)
	tx = types.NewBaseTx(msg, []types.Signature{{
		Sig: crypto.SignatureSecp256k1(sig),
	}})

	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 5}})

	// Run a check
	cres := cc.Check(tx)
	assert.Equal(t, sdk.CodeType(0),
		sdk.CodeType(cres.Code), cres.Log)

	dres := cc.Deliver(tx)

	assert.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

	// Create context
	ctx := cc.NewContext(false, abci.Header{})

	// Retrieve UTXO from context
	utxo := cc.utxoMapper.GetUTXO(ctx, types.NewPosition(5, 0, 0, 0))
	expected := types.NewBaseUTXO(addrA, [2]common.Address{addrB, common.Address{}}, 100, types.NewPosition(5, 0, 0, 0))

	assert.Equal(t, expected, utxo, "UTXO did not get added to store correctly")

}