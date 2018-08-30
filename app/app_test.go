package app

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"crypto/ecdsa"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	crypto "github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	types "github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	rlp "github.com/ethereum/go-ethereum/rlp"
)

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

	genState := GenesisState{
		UTXOs: genUTXOs,
	}

	appStateBytes, err := json.Marshal(genState)
	if err != nil {
		panic(err)
	}

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

// Returns a confirmsig array signed by privKey1 and privKey2
func CreateConfirmSig(position1 types.Position, privKey1, privKey2 *ecdsa.PrivateKey, two_inputs bool) (confirmSigs [2]types.Signature) {
	confirmBytes := position1.GetSignBytes()
	hash := ethcrypto.Keccak256(confirmBytes)
	confirmSig, _ := ethcrypto.Sign(hash, privKey1)

	var confirmSig2 []byte
	if two_inputs {
		confirmSig2, _ = ethcrypto.Sign(hash, privKey2)
	}
	confirmSigs = [2]types.Signature{types.Signature{confirmSig}, types.Signature{confirmSig2}}
	return confirmSigs
}

// helper for constructing single or double input tx
func GetTx(msg types.SpendMsg, privKeyA, privKeyB *ecdsa.PrivateKey, two_sigs bool) (tx types.BaseTx) {
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig1, _ := ethcrypto.Sign(hash, privKeyA)

	tx = types.NewBaseTx(msg, []types.Signature{{
		Sig: sig1,
	}})

	if two_sigs {
		sig2, _ := ethcrypto.Sign(hash, privKeyB)
		tx.Signatures = append(tx.Signatures, types.Signature{sig2})
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
		Sig: crypto.SignatureSecp256k1(sig),
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
	msg.ConfirmSigs1 = CreateConfirmSig(types.NewPosition(0, 0, 0, 1), privKeyA, &ecdsa.PrivateKey{}, false)

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

	// Run a check
	cres := cc.Check(tx)
	require.Equal(t, sdk.CodeType(0),
		sdk.CodeType(cres.Code), cres.Log)

	// Deliver tx, updates states
	dres := cc.Deliver(tx)

	require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

	// Create context
	ctx := cc.NewContext(false, abci.Header{})

	// Retrieve UTXO from context
	utxo := cc.utxoMapper.GetUTXO(ctx, addrB, types.NewPosition(1, 0, 0, 0))
	expected := types.NewBaseUTXO(addrB, [2]common.Address{addrA, common.Address{}}, 100, types.NewPosition(1, 0, 0, 0))

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
	msg.ConfirmSigs1 = CreateConfirmSig(types.NewPosition(0, 0, 0, 1), privKeyA, &ecdsa.PrivateKey{}, false)

	// Signs the hash of the transaction
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ := ethcrypto.Sign(hash, privKeyA)
	tx := types.NewBaseTx(msg, []types.Signature{{
		Sig: crypto.SignatureSecp256k1(sig),
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
	msg.ConfirmSigs1 = CreateConfirmSig(types.NewPosition(1, 0, 0, 0), privKeyA, &ecdsa.PrivateKey{}, false)

	// Signs the hash of the transaction
	hash = ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ = ethcrypto.Sign(hash, privKeyB)
	tx = types.NewBaseTx(msg, []types.Signature{{
		Sig: crypto.SignatureSecp256k1(sig),
	}})

	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 5}})

	// Run a check
	cres := cc.Check(tx)
	require.Equal(t, sdk.CodeType(0),
		sdk.CodeType(cres.Code), cres.Log)

	dres := cc.Deliver(tx)

	require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

	// Create context
	ctx := cc.NewContext(false, abci.Header{})

	// Retrieve UTXO from context
	utxo := cc.utxoMapper.GetUTXO(ctx, addrA, types.NewPosition(5, 0, 0, 0))
	expected := types.NewBaseUTXO(addrA, [2]common.Address{addrB, common.Address{}}, 100, types.NewPosition(5, 0, 0, 0))

	require.Equal(t, expected, utxo, "UTXO did not get added to store correctly")

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
		owner_index1   uint64
		position1      types.Position
		input_index1_1 uint64
		input_index1_2 uint64
		owner_index2   uint64
		addr2          common.Address
		position2      types.Position
		input_index2_1 uint64
		input_index2_2 uint64
		newowner1      common.Address
		denom1         uint64
		newowner2      common.Address
		denom2         uint64
	}{
		// Test Case 0: 1 input 2 output
		// Tx spends the genesis deposit and creates 2 new ouputs for addr[1] and addr[2]
		{0, types.NewPosition(0, 0, 0, 1), 0, 0, 0, common.Address{}, types.Position{}, 0, 0, addrs[1], 20, addrs[2], 80},

		// Test Case 1: 2 different inputs, 1 output
		// Tx spends outputs from test case 0 and creates 1 output for addr[3]
		{1, types.NewPosition(7, 0, 0, 0), 0, 0, 2, addrs[2], types.NewPosition(7, 0, 1, 0), 0, 0, addrs[3], 100, common.Address{}, 0},

		// Test Case 2: 1 input 2 ouput
		// Tx spends output from test case 1 and creates 2 new outputs for addr[3] and addr[4]
		{3, types.NewPosition(8, 0, 0, 0), 1, 2, 0, common.Address{}, types.Position{}, 0, 0, addrs[3], 75, addrs[4], 25},

		// Test Case 3: 2 different inputs 2 outputs
		// Tx spends outputs from test case 2 and creates 2 new outputs both for addr[3]
		{3, types.NewPosition(9, 0, 0, 0), 3, 0, 4, addrs[4], types.NewPosition(9, 0, 1, 0), 3, 0, addrs[3], 70, addrs[3], 30},

		// Test Case 4: 2 same inputs, 1 output (merge)
		// Tx spends outputs from test case 3 and creates 1 new output for addr[3]
		{3, types.NewPosition(10, 0, 0, 0), 3, 4, 3, addrs[3], types.NewPosition(10, 0, 1, 0), 3, 4, addrs[3], 100, common.Address{}, 0},
	}

	for index, tc := range cases {
		cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 7 + int64(index)}})

		msg := types.SpendMsg{
			Blknum1:      tc.position1.Blknum,
			Txindex1:     tc.position1.TxIndex,
			Oindex1:      tc.position1.Oindex,
			DepositNum1:  tc.position1.DepositNum,
			Owner1:       addrs[tc.owner_index1],
			ConfirmSigs1: CreateConfirmSig(tc.position1, keys[tc.input_index1_1], keys[tc.input_index1_2], true),
			Blknum2:      tc.position2.Blknum,
			Txindex2:     tc.position2.TxIndex,
			Oindex2:      tc.position2.Oindex,
			DepositNum2:  tc.position2.DepositNum,
			Owner2:       tc.addr2,
			ConfirmSigs2: CreateConfirmSig(tc.position2, keys[tc.input_index2_1], keys[tc.input_index2_2], true),
			Newowner1:    tc.newowner1,
			Denom1:       tc.denom1,
			Newowner2:    tc.newowner2,
			Denom2:       tc.denom2,
			Fee:          0,
		}

		tx := GetTx(msg, keys[tc.owner_index1], keys[tc.owner_index2], !utils.ZeroAddress(msg.Owner2))

		// Run a check
		cres := cc.Check(tx)
		require.Equal(t, sdk.CodeType(0), sdk.CodeType(cres.Code), cres.Log)

		dres := cc.Deliver(tx)

		require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

		// Create context
		ctx := cc.NewContext(false, abci.Header{})

		// Retrieve utxo from context
		utxo := cc.utxoMapper.GetUTXO(ctx, tc.newowner1, types.NewPosition(uint64(index)+7, 0, 0, 0))
		expected := types.NewBaseUTXO(tc.newowner1, [2]common.Address{msg.Owner1, msg.Owner2}, tc.denom1, types.NewPosition(uint64(index)+7, 0, 0, 0))
		require.Equal(t, expected, utxo, fmt.Sprintf("First UTXO did not get added to the utxo store correctly. Failed on test case: %d", index))

		if !utils.ZeroAddress(msg.Newowner2) {
			utxo = cc.utxoMapper.GetUTXO(ctx, tc.newowner2, types.NewPosition(uint64(index)+7, 0, 1, 0))
			expected = types.NewBaseUTXO(tc.newowner2, [2]common.Address{msg.Owner1, msg.Owner2}, tc.denom2, types.NewPosition(uint64(index)+7, 0, 1, 0))
			require.Equal(t, expected, utxo, fmt.Sprintf("Second UTXO did not get added to the utxo store correctly. Failed on test case: %d", index))
		}

		// Check that inputs were removed
		utxo = cc.utxoMapper.GetUTXO(ctx, msg.Owner1, tc.position1)
		require.Nil(t, utxo, fmt.Sprintf("first input was not removed from the utxo store. Failed on test case: %d", index))

		if !utils.ZeroAddress(msg.Owner2) {
			utxo = cc.utxoMapper.GetUTXO(ctx, msg.Owner2, tc.position2)
			require.Nil(t, utxo, fmt.Sprintf("second input was not removed from the utxo store. Failed on test case: %d", index))
		}

		cc.EndBlock(abci.RequestEndBlock{})
		cc.Commit()
	}
}

// Test that several txs can go into a block and that txindex increments correctly
// Change value of N to increase or decrease txs in the block
func TestMultiTxBlocks(t *testing.T) {
	const N = 20
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
		msgs[i].ConfirmSigs1 = CreateConfirmSig(types.NewPosition(0, 0, 0, i+1), keys[i], &ecdsa.PrivateKey{}, false)
		txs[i] = GetTx(msgs[i], keys[i], &ecdsa.PrivateKey{}, false)

		cres := cc.Check(txs[i])
		require.Equal(t, sdk.CodeType(0), sdk.CodeType(cres.Code), cres.Log)

		dres := cc.Deliver(txs[i])
		require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

	}
	ctx := cc.NewContext(false, abci.Header{})

	// Retrieve and check UTXO from context
	for i := uint16(0); i < N; i++ {
		utxo := cc.utxoMapper.GetUTXO(ctx, addrs[i], types.NewPosition(1, i, 0, 0))
		expected := types.NewBaseUTXO(addrs[i], [2]common.Address{addrs[i], common.Address{}}, 100, types.NewPosition(1, i, 0, 0))

		require.Equal(t, expected, utxo, fmt.Sprintf("UTXO %d did not get added to store correctly", i+1))

		utxo = cc.utxoMapper.GetUTXO(ctx, addrs[i], types.NewPosition(0, 0, 0, uint64(i)+1))
		require.Nil(t, utxo, fmt.Sprintf("deposit %d did not get removed correctly from the utxo store", i+1))
	}

	cc.EndBlock(abci.RequestEndBlock{})
	cc.Commit()
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})

	// send to different address
	for i := uint16(0); i < N; i++ {
		msgs[i].Blknum1 = 1
		msgs[i].Txindex1 = i
		msgs[i].DepositNum1 = 0
		msgs[i].ConfirmSigs1 = CreateConfirmSig(types.NewPosition(1, i, 0, 0), keys[i], &ecdsa.PrivateKey{}, false)
		msgs[i].Newowner1 = addrs[(i+1)%N]
		txs[i] = GetTx(msgs[i], keys[i], &ecdsa.PrivateKey{}, false)

		cres := cc.Check(txs[i])
		require.Equal(t, sdk.CodeType(0), sdk.CodeType(cres.Code), cres.Log)

		dres := cc.Deliver(txs[i])
		require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)
	}

	ctx = cc.NewContext(false, abci.Header{})

	// Retrieve and check UTXO from context
	for i := uint16(0); i < N; i++ {
		utxo := cc.utxoMapper.GetUTXO(ctx, addrs[(i+1)%N], types.NewPosition(2, i, 0, 0))
		expected := types.NewBaseUTXO(addrs[(i+1)%N], [2]common.Address{addrs[i], common.Address{}}, 100, types.NewPosition(2, i, 0, 0))

		require.Equal(t, expected, utxo, fmt.Sprintf("UTXO %d did not get added to store correctly", i+1))

		utxo = cc.utxoMapper.GetUTXO(ctx, addrs[i], types.NewPosition(1, i, 0, 0))
		require.Nil(t, utxo, fmt.Sprintf("UTXO %d  did not get removed from the utxo store correctly", i))
	}

}
