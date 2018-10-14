package app

import (
	"fmt"
	"os"
	"testing"

	"crypto/ecdsa"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	secp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	types "github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	rlp "github.com/ethereum/go-ethereum/rlp"
)

/*
	Note: Check() has been taken out from testing
	at the moment because it increments txIndex

*/

func newChildChain() *ChildChain {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return NewChildChain(logger, db, nil)
}

// Creates a deposit of value 100 for each address in input
func InitTestChain(cc *ChildChain, addrs ...common.Address) {
	var genUTXOs []GenesisUTXO
	for i, addr := range addrs {
		genUTXOs = append(genUTXOs, NewGenesisUTXO(addr.Hex(), "100", [4]string{"0", "0", "0", fmt.Sprintf("%d", i+1)}))
	}

	pubKey := secp256k1.GenPrivKey().PubKey()

	genState := GenesisState{
		Validator: pubKey,
		UTXOs: genUTXOs,
	}

	appStateBytes, err := cc.cdc.MarshalJSON(genState)
	if err != nil {
		panic(err)
	}

	initRequest := abci.RequestInitChain{AppStateBytes: appStateBytes}
	cc.InitChain(initRequest)
}

func GenerateSimpleMsg(Owner0, NewOwner0 common.Address, position [4]uint64, amount0 uint64, fee uint64) types.SpendMsg {
	confirmSigs := [2]types.Signature{types.Signature{}, types.Signature{}}
	return types.SpendMsg{
		Blknum0:      position[0],
		Txindex0:     uint16(position[1]),
		Oindex0:      uint8(position[2]),
		DepositNum0:  position[3],
		Owner0:       Owner0,
		ConfirmSigs0: confirmSigs,
		Blknum1:      0,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       common.Address{},
		ConfirmSigs1: confirmSigs,
		Newowner0:    NewOwner0,
		Amount0:      amount0,
		Newowner1:    common.Address{},
		Amount1:      0,
		FeeAmount:    fee,
	}
}

// Returns a confirmsig array signed by privKey0 and privKey1
func CreateConfirmSig(position types.PlasmaPosition, privKey0, privKey1 *ecdsa.PrivateKey, two_inputs bool) (confirmSigs [2]types.Signature) {
	confirmBytes := position.GetSignBytes()
	hash := ethcrypto.Keccak256(confirmBytes)
	confirmSig, _ := ethcrypto.Sign(hash, privKey0)

	var confirmSig1 []byte
	if two_inputs {
		confirmSig1, _ = ethcrypto.Sign(hash, privKey1)
	}
	confirmSigs = [2]types.Signature{types.Signature{confirmSig}, types.Signature{confirmSig1}}
	return confirmSigs
}

// helper for constructing single or double input tx
func GetTx(msg types.SpendMsg, privKeyA, privKeyB *ecdsa.PrivateKey, two_sigs bool) (tx types.BaseTx) {
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ := ethcrypto.Sign(hash, privKeyA)

	tx = types.NewBaseTx(msg, []types.Signature{{
		Sig: sig,
	}})

	if two_sigs {
		sig1, _ := ethcrypto.Sign(hash, privKeyB)
		tx.Signatures = append(tx.Signatures, types.Signature{sig1})
	}

	return tx
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
		Sig: sig,
	}})

	txBytes, err := rlp.EncodeToBytes(tx)

	require.NoError(t, err)

	// Must Commit to set checkState
	cc.BeginBlock(abci.RequestBeginBlock{})
	cc.EndBlock(abci.RequestEndBlock{})
	cc.Commit()

	// Run a check
	cres := cc.CheckTx(txBytes)
	require.Equal(t, sdk.CodeType(6),
		sdk.CodeType(cres.Code), cres.Log)

	// Simulate a Block
	cc.BeginBlock(abci.RequestBeginBlock{})
	dres := cc.DeliverTx(txBytes)
	require.Equal(t, sdk.CodeType(6), sdk.CodeType(dres.Code), dres.Log)

}

func TestSpendDeposit(t *testing.T) {
	cc := newChildChain()

	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()
	addrA := utils.PrivKeyToAddress(privKeyA)
	addrB := utils.PrivKeyToAddress(privKeyB)

	InitTestChain(cc, addrA)

	msg := GenerateSimpleMsg(addrA, addrB, [4]uint64{0, 0, 0, 1}, 100, 0)

	// Set confirm signatures
	msg.ConfirmSigs0 = CreateConfirmSig(types.NewPlasmaPosition(0, 0, 0, 1), privKeyA, &ecdsa.PrivateKey{}, false)

	// Signs the hash of the transaction
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ := ethcrypto.Sign(hash, privKeyA)
	tx := types.NewBaseTx(msg, []types.Signature{{
		Sig: sig,
	}})

	// Must commit for checkState to be set correctly. Should be fixed in next version of SDK
	cc.BeginBlock(abci.RequestBeginBlock{})
	cc.EndBlock(abci.RequestEndBlock{})
	cc.Commit()

	// Simulate a block
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})

	// Deliver tx, updates states
	dres := cc.Deliver(tx)

	require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

	// Create context
	ctx := cc.NewContext(false, abci.Header{})

	// Retrieve UTXO from context
	position := types.NewPlasmaPosition(1, 0, 0, 0)
	utxo := cc.utxoMapper.GetUTXO(ctx, addrB.Bytes(), position)
	expected := types.NewBaseUTXO(addrB, [2]common.Address{addrA, common.Address{}}, 100, "", position)

	require.Equal(t, expected, utxo, "UTXO did not get added to store correctly")

}

func TestSpendTx(t *testing.T) {
	cc := newChildChain()

	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()
	addrA := utils.PrivKeyToAddress(privKeyA)
	addrB := utils.PrivKeyToAddress(privKeyB)

	InitTestChain(cc, addrA)
	cc.Commit()

	msg := GenerateSimpleMsg(addrA, addrB, [4]uint64{0, 0, 0, 1}, 100, 0)

	// Set confirm signatures
	msg.ConfirmSigs0 = CreateConfirmSig(types.NewPlasmaPosition(0, 0, 0, 1), privKeyA, &ecdsa.PrivateKey{}, false)

	// Signs the hash of the transaction
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ := ethcrypto.Sign(hash, privKeyA)
	tx := types.NewBaseTx(msg, []types.Signature{{
		Sig: sig,
	}})

	// Simulate a block
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})

	// Deliver tx, updates states
	res := cc.Deliver(tx)

	require.True(t, res.IsOK(), res.Log)

	cc.EndBlock(abci.RequestEndBlock{})
	cc.Commit()

	// Test that spending from a non-deposit/non-genesis UTXO works

	// generate simple msg
	msg = GenerateSimpleMsg(addrB, addrA, [4]uint64{1, 0, 0, 0}, 100, 0)

	// Set confirm signatures
	msg.ConfirmSigs0 = CreateConfirmSig(types.NewPlasmaPosition(1, 0, 0, 0), privKeyA, &ecdsa.PrivateKey{}, false)

	// Signs the hash of the transaction
	hash = ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ = ethcrypto.Sign(hash, privKeyB)
	tx = types.NewBaseTx(msg, []types.Signature{{
		Sig: sig,
	}})

	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 5}})

	dres := cc.Deliver(tx)

	require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

	// Create context
	ctx := cc.NewContext(false, abci.Header{})

	// Retrieve UTXO from context
	position := types.NewPlasmaPosition(5, 0, 0, 0)
	utxo := cc.utxoMapper.GetUTXO(ctx, addrA.Bytes(), position)
	expected := types.NewBaseUTXO(addrA, [2]common.Address{addrB, common.Address{}}, 100, "", position)

	require.Equal(t, expected, utxo, "UTXO did not get added to store correctly")

}

// helper struct for readability
type Input struct {
	owner_index  int64
	addr         common.Address
	position     types.PlasmaPosition
	input_index0 int64
	input_index1 int64
}

// Tests 1 input 2 ouput, 2 input (different addresses) 1 output,
// 2 input (different addresses) 2 ouputs, and 2 input (same address) 1 output
func TestDifferentTxForms(t *testing.T) {
	// Initialize child chain with deposit
	cc := newChildChain()
	var keys [6]*ecdsa.PrivateKey
	var addrs []common.Address

	for i := 0; i < 6; i++ {
		keys[i], _ = ethcrypto.GenerateKey()
		addrs = append(addrs, utils.PrivKeyToAddress(keys[i]))
	}

	InitTestChain(cc, addrs...)
	cc.Commit()

	cases := []struct {
		input0    Input
		input1    Input
		newowner0 common.Address
		amount0   uint64
		newowner1 common.Address
		amount1   uint64
	}{
		// Test Case 0: 1 input 2 output
		// Tx spends the genesis deposit and creates 2 new ouputs for addr[1] and addr[2]
		{
			Input{0, addrs[0], types.NewPlasmaPosition(0, 0, 0, 1), 0, -1},
			Input{0, common.Address{}, types.PlasmaPosition{}, -1, -1},
			addrs[1], 20,
			addrs[2], 80,
		},

		// Test Case 1: 2 different inputs, 1 output
		// Tx spends outputs from test case 0 and creates 1 output for addr[3]
		{
			Input{1, addrs[1], types.NewPlasmaPosition(7, 0, 0, 0), 0, -1},
			Input{2, addrs[2], types.NewPlasmaPosition(7, 0, 1, 0), 0, -1},
			addrs[3], 100,
			common.Address{}, 0,
		},

		// Test Case 2: 1 input 2 ouput
		// Tx spends output from test case 1 and creates 2 new outputs for addr[3] and addr[4]
		{
			Input{3, addrs[3], types.NewPlasmaPosition(8, 0, 0, 0), 1, 2},
			Input{0, common.Address{}, types.PlasmaPosition{}, -1, -1},
			addrs[3], 75,
			addrs[4], 25,
		},

		// Test Case 3: 2 different inputs 2 outputs
		// Tx spends outputs from test case 2 and creates 2 new outputs both for addr[3]
		{
			Input{3, addrs[3], types.NewPlasmaPosition(9, 0, 0, 0), 3, -1},
			Input{4, addrs[4], types.NewPlasmaPosition(9, 0, 1, 0), 3, -1},
			addrs[3], 70,
			addrs[3], 30,
		},

		// Test Case 4: 2 same inputs, 1 output (merge)
		// Tx spends outputs from test case 3 and creates 1 new output for addr[3]
		{
			Input{3, addrs[3], types.NewPlasmaPosition(10, 0, 0, 0), 3, 4},
			Input{3, addrs[3], types.NewPlasmaPosition(10, 0, 1, 0), 3, 4},
			addrs[3], 100,
			common.Address{}, 0,
		},
	}

	for index, tc := range cases {
		cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 7 + int64(index)}})

		input0_index1 := utils.GetIndex(tc.input0.input_index1)
		input1_index0 := utils.GetIndex(tc.input1.input_index0)
		input1_index1 := utils.GetIndex(tc.input1.input_index1)
		msg := types.SpendMsg{
			Blknum0:      tc.input0.position.Blknum,
			Txindex0:     tc.input0.position.TxIndex,
			Oindex0:      tc.input0.position.Oindex,
			DepositNum0:  tc.input0.position.DepositNum,
			Owner0:       tc.input0.addr,
			ConfirmSigs0: CreateConfirmSig(tc.input0.position, keys[tc.input0.input_index0], keys[input0_index1], tc.input0.input_index1 != -1),
			Blknum1:      tc.input1.position.Blknum,
			Txindex1:     tc.input1.position.TxIndex,
			Oindex1:      tc.input1.position.Oindex,
			DepositNum1:  tc.input1.position.DepositNum,
			Owner1:       tc.input1.addr,
			ConfirmSigs1: CreateConfirmSig(tc.input1.position, keys[input1_index0], keys[input1_index1], tc.input1.input_index1 != -1),
			Newowner0:    tc.newowner0,
			Amount0:      tc.amount0,
			Newowner1:    tc.newowner1,
			Amount1:      tc.amount1,
			FeeAmount:    0,
		}

		tx := GetTx(msg, keys[tc.input0.owner_index], keys[tc.input1.owner_index], !utils.ZeroAddress(msg.Owner1))

		dres := cc.Deliver(tx)

		require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

		// Create context
		ctx := cc.NewContext(false, abci.Header{})

		// Retrieve utxo from context
		position := types.NewPlasmaPosition(uint64(index)+7, 0, 0, 0)
		utxo := cc.utxoMapper.GetUTXO(ctx, tc.newowner0.Bytes(), position)
		expected := types.NewBaseUTXO(tc.newowner0, [2]common.Address{msg.Owner0, msg.Owner1}, tc.amount0, "", position)
		require.Equal(t, expected, utxo, fmt.Sprintf("First UTXO did not get added to the utxo store correctly. Failed on test case: %d", index))

		if !utils.ZeroAddress(msg.Newowner1) {
			position = types.NewPlasmaPosition(uint64(index)+7, 0, 1, 0)
			utxo = cc.utxoMapper.GetUTXO(ctx, tc.newowner1.Bytes(), position)
			expected = types.NewBaseUTXO(tc.newowner1, [2]common.Address{msg.Owner0, msg.Owner1}, tc.amount1, "", position)
			require.Equal(t, expected, utxo, fmt.Sprintf("Second UTXO did not get added to the utxo store correctly. Failed on test case: %d", index))
		}

		// Check that inputs were removed
		utxo = cc.utxoMapper.GetUTXO(ctx, msg.Owner0.Bytes(), tc.input0.position)
		require.Nil(t, utxo, fmt.Sprintf("first input was not removed from the utxo store. Failed on test case: %d", index))

		if !utils.ZeroAddress(msg.Owner1) {
			utxo = cc.utxoMapper.GetUTXO(ctx, msg.Owner1.Bytes(), tc.input1.position)
			require.Nil(t, utxo, fmt.Sprintf("second input was not removed from the utxo store. Failed on test case: %d", index))
		}

		cc.EndBlock(abci.RequestEndBlock{})
		cc.Commit()
	}
}

// Test that several txs can go into a block and that txindex increments correctly
// Change value of N to increase or decrease txs in the block
func TestMultiTxBlocks(t *testing.T) {
	const N = 5
	// Initialize child chain with deposit
	cc := newChildChain()
	var keys [N]*ecdsa.PrivateKey
	var addrs []common.Address
	var msgs [N]types.SpendMsg
	var txs [N]sdk.Tx

	for i := 0; i < N; i++ {
		keys[i], _ = ethcrypto.GenerateKey()
		addrs = append(addrs, utils.PrivKeyToAddress(keys[i]))
	}

	InitTestChain(cc, addrs...)
	cc.Commit()
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})

	for i := uint64(0); i < N; i++ {
		msgs[i] = GenerateSimpleMsg(addrs[i], addrs[i], [4]uint64{0, 0, 0, i + 1}, 100, 0)
		msgs[i].ConfirmSigs0 = CreateConfirmSig(types.NewPlasmaPosition(0, 0, 0, i+1), keys[i], &ecdsa.PrivateKey{}, false)
		txs[i] = GetTx(msgs[i], keys[i], &ecdsa.PrivateKey{}, false)

		dres := cc.Deliver(txs[i])
		require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

	}
	ctx := cc.NewContext(false, abci.Header{})

	// Retrieve and check UTXO from context
	for i := uint16(0); i < N; i++ {
		position := types.NewPlasmaPosition(1, i, 0, 0)
		utxo := cc.utxoMapper.GetUTXO(ctx, addrs[i].Bytes(), position)
		expected := types.NewBaseUTXO(addrs[i], [2]common.Address{addrs[i], common.Address{}}, 100, "", position)

		require.Equal(t, expected, utxo, fmt.Sprintf("UTXO %d did not get added to store correctly", i+1))

		position = types.NewPlasmaPosition(0, 0, 0, uint64(i)+1)
		utxo = cc.utxoMapper.GetUTXO(ctx, addrs[i].Bytes(), position)
		require.Nil(t, utxo, fmt.Sprintf("deposit %d did not get removed correctly from the utxo store", i+1))
	}

	cc.EndBlock(abci.RequestEndBlock{})
	cc.Commit()
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})

	// send to different address
	for i := uint16(0); i < N; i++ {
		msgs[i].Blknum0 = 1
		msgs[i].Txindex0 = i
		msgs[i].DepositNum0 = 0
		msgs[i].ConfirmSigs0 = CreateConfirmSig(types.NewPlasmaPosition(1, i, 0, 0), keys[i], &ecdsa.PrivateKey{}, false)
		msgs[i].Newowner0 = addrs[(i+1)%N]
		txs[i] = GetTx(msgs[i], keys[i], &ecdsa.PrivateKey{}, false)

		dres := cc.Deliver(txs[i])
		require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)
	}

	ctx = cc.NewContext(false, abci.Header{})

	// Retrieve and check UTXO from context
	for i := uint16(0); i < N; i++ {
		utxo := cc.utxoMapper.GetUTXO(ctx, addrs[(i+1)%N].Bytes(), types.NewPlasmaPosition(2, i, 0, 0))
		expected := types.NewBaseUTXO(addrs[(i+1)%N], [2]common.Address{addrs[i], common.Address{}}, 100, "", types.NewPlasmaPosition(2, i, 0, 0))

		require.Equal(t, expected, utxo, fmt.Sprintf("UTXO %d did not get added to store correctly", i+1))

		utxo = cc.utxoMapper.GetUTXO(ctx, addrs[i].Bytes(), types.NewPlasmaPosition(1, i, 0, 0))
		require.Nil(t, utxo, fmt.Sprintf("UTXO %d  did not get removed from the utxo store correctly", i))
	}

}
