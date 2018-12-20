package app

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"crypto/ecdsa"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	secp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	types "github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	rlp "github.com/ethereum/go-ethereum/rlp"
)

const (
	privkey            = "9cd69f009ac86203e54ec50e3686de95ff6126d3b30a19f926a0fe9323c17181"
	nodeURL            = "ws://127.0.0.1:8545"
	plasmaContractAddr = "5cae340fb2c2bb0a2f194a95cda8a1ffdc9d2f85"
)

/* Note: Since the headers only contain information about the height
 *       of a block, the updated header is not being set in the context
 *		 and therefore the block hash set in end blocker will be nil.
 *	 	 This is only true for these tests, it works as expected when
 *		 you run a full node via plasmad.
 */

func newChildChain() *ChildChain {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	privkeyFile, _ := ioutil.TempFile("", "privateKey")
	privkeyFile.Write([]byte(privkey))
	defer os.Remove(privkeyFile.Name())
	return NewChildChain(logger, db, nil, SetEthConfig(true, privkeyFile.Name(), plasmaContractAddr, nodeURL, "0", "0"))
}

// Adds a initial utxo at the specified position
// Note: The input keys, txHash, and blockHash are not accurate
// Use only for spending
func storeInitUTXO(cc *ChildChain, position types.PlasmaPosition, addr common.Address) {
	// Add blockhash to plasmaStore
	blknumKey := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(blknumKey, uint64(position.Blknum))
	key := append(utils.RootHashPrefix, blknumKey...)
	blockHash := []byte("merkle root")

	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: int64(position.Blknum) + 1}})
	ctx := cc.NewContext(false, abci.Header{})
	cc.plasmaStore.Set(ctx, key, blockHash)

	// Creates an input key that maps to itself
	var inputKeys [][]byte
	inKey := cc.utxoMapper.ConstructKey(addr.Bytes(), position)
	inputKeys = append(inputKeys, inKey)

	txhash := []byte("txhash")
	input := utxo.NewUTXOwithInputs(addr.Bytes(), 100, "Ether", position, txhash, inputKeys)
	cc.utxoMapper.ReceiveUTXO(ctx, input)
}

// Creates a deposit of value 100 for each address in input
func InitTestChain(cc *ChildChain, valAddr common.Address, addrs ...common.Address) {
	var genUTXOs []GenesisUTXO
	for i, addr := range addrs {
		genUTXOs = append(genUTXOs, NewGenesisUTXO(addr.Hex(), "100", [4]string{"0", "0", "0", fmt.Sprintf("%d", i+1)}))
	}

	pubKey := secp256k1.GenPrivKey().PubKey()

	genValidator := GenesisValidator{
		ConsPubKey: pubKey,
		Address:    valAddr.String(),
	}

	genState := GenesisState{
		Validator: genValidator,
		UTXOs:     genUTXOs,
	}

	appStateBytes, err := cc.cdc.MarshalJSON(genState)
	if err != nil {
		panic(err)
	}

	initRequest := abci.RequestInitChain{AppStateBytes: appStateBytes}
	cc.InitChain(initRequest)
}

func GenerateSimpleMsg(Owner0, NewOwner0 common.Address, position [4]uint64, amount0 uint64, fee uint64) types.SpendMsg {
	var confirmSigs [][65]byte
	return types.SpendMsg{
		Blknum0:           position[0],
		Txindex0:          uint16(position[1]),
		Oindex0:           uint8(position[2]),
		DepositNum0:       position[3],
		Owner0:            Owner0,
		Input0ConfirmSigs: confirmSigs,
		Blknum1:           0,
		Txindex1:          0,
		Oindex1:           0,
		DepositNum1:       0,
		Owner1:            common.Address{},
		Input1ConfirmSigs: confirmSigs,
		Newowner0:         NewOwner0,
		Amount0:           amount0,
		Newowner1:         common.Address{},
		Amount1:           0,
		FeeAmount:         fee,
	}
}

// Returns a confirmsig array signed by privKey0 and privKey1
func CreateConfirmSig(hash []byte, privKey0, privKey1 *ecdsa.PrivateKey, two_inputs bool) (confirmSigs [][65]byte) {

	var confirmSig0 [65]byte
	signHash := utils.SignHash(hash)
	confirmSig0Slice, _ := ethcrypto.Sign(signHash, privKey0)
	copy(confirmSig0[:], confirmSig0Slice)
	confirmSigs = append(confirmSigs, confirmSig0)

	var confirmSig1 [65]byte
	if two_inputs {
		confirmSig1Slice, _ := ethcrypto.Sign(signHash, privKey1)
		copy(confirmSig1[:], confirmSig1Slice)
		confirmSigs = append(confirmSigs, confirmSig1)
	}
	return confirmSigs
}

func getInputKeys(mapper utxo.Mapper, inputs ...Input) (res [][]byte) {
	for _, in := range inputs {
		if !reflect.DeepEqual(in.addr, common.Address{}) {
			res = append(res, mapper.ConstructKey(in.addr.Bytes(), in.position))
		}
	}
	return res
}

// helper for constructing single or double input tx
func GetTx(msg types.SpendMsg, privKeyA, privKeyB *ecdsa.PrivateKey, two_sigs bool) (tx types.BaseTx) {
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	signHash := utils.SignHash(hash)
	var sigs [2][65]byte
	sig, _ := ethcrypto.Sign(signHash, privKeyA)
	copy(sigs[0][:], sig)

	if two_sigs {
		sig1, _ := ethcrypto.Sign(signHash, privKeyB)
		copy(sigs[1][:], sig1)
	}

	tx = types.NewBaseTx(msg, sigs)
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
	var sigs [2][65]byte
	sig, _ := ethcrypto.Sign(hash, privKeyA)
	copy(sigs[0][:], sig)
	tx := types.NewBaseTx(msg, sigs)

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

func TestSpendTx(t *testing.T) {
	cc := newChildChain()

	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()
	addrA := utils.PrivKeyToAddress(privKeyA)
	addrB := utils.PrivKeyToAddress(privKeyB)

	InitTestChain(cc, utils.GenerateAddress(), addrA)
	cc.Commit()

	// Add a UTXO into the utxoMapper
	position := types.NewPlasmaPosition(1, 0, 0, 0)
	storeInitUTXO(cc, position, addrB)
	cc.Commit()

	txHash := []byte("txhash")
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 5}})

	// Create context
	ctx := cc.NewContext(false, abci.Header{})
	blknumKey := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(blknumKey, uint64(1))
	key := append(utils.RootHashPrefix, blknumKey...)
	blockHash := cc.plasmaStore.Get(ctx, key)

	// Test that spending from a non-deposit/non-genesis UTXO works

	// generate simple msg
	msg := GenerateSimpleMsg(addrB, addrA, [4]uint64{1, 0, 0, 0}, 100, 0)

	hash := tmhash.Sum(append(txHash, blockHash...))
	// Set confirm signatures
	msg.Input0ConfirmSigs = CreateConfirmSig(hash, privKeyB, &ecdsa.PrivateKey{}, false)

	// Signs the hash of the transaction
	tx := GetTx(msg, privKeyB, nil, false)
	txBytes, _ := rlp.EncodeToBytes(tx)
	txHash = tmhash.Sum(txBytes)

	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 5}})

	dres := cc.DeliverTx(txBytes)

	require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

	// Retrieve UTXO from context
	position = types.NewPlasmaPosition(5, 0, 0, 0)
	actual := cc.utxoMapper.GetUTXO(ctx, addrA.Bytes(), position)

	inputKey := cc.utxoMapper.ConstructKey(addrB.Bytes(), types.NewPlasmaPosition(1, 0, 0, 0))
	expected := utxo.NewUTXOwithInputs(addrA.Bytes(), 100, "Ether", position, txHash, [][]byte{inputKey})

	require.Equal(t, expected, actual, "UTXO did not get added to store correctly")

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

	InitTestChain(cc, utils.GenerateAddress(), addrs...)

	// Add inital utxo
	position := types.NewPlasmaPosition(6, 0, 0, 0)
	storeInitUTXO(cc, position, addrs[0])

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
		// Tx spends the init tx and creates 2 new ouputs for addr[1] and addr[2]
		{
			Input{0, addrs[0], types.NewPlasmaPosition(6, 0, 0, 0), 0, -1},
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

		// Create context
		ctx := cc.NewContext(false, abci.Header{Height: 7 + int64(index)})

		msg := types.SpendMsg{
			Blknum0:     tc.input0.position.Blknum,
			Txindex0:    tc.input0.position.TxIndex,
			Oindex0:     tc.input0.position.Oindex,
			DepositNum0: tc.input0.position.DepositNum,
			Owner0:      tc.input0.addr,
			Blknum1:     tc.input1.position.Blknum,
			Txindex1:    tc.input1.position.TxIndex,
			Oindex1:     tc.input1.position.Oindex,
			DepositNum1: tc.input1.position.DepositNum,
			Owner1:      tc.input1.addr,
			Newowner0:   tc.newowner0,
			Amount0:     tc.amount0,
			Newowner1:   tc.newowner1,
			Amount1:     tc.amount1,
			FeeAmount:   0,
		}

		if tc.input0.position.DepositNum == 0 && tc.input0.position.Blknum != 0 {
			// note: all cases currently have inputs belonging to the previous tx
			// and therefore we only need to grab the first txhash from the inptus
			input_utxo := cc.utxoMapper.GetUTXO(ctx, tc.input0.addr.Bytes(), tc.input0.position)
			blknumKey := make([]byte, binary.MaxVarintLen64)
			binary.PutUvarint(blknumKey, uint64(7+int64(index-1)))
			key := append(utils.RootHashPrefix, blknumKey...)
			blockhash := cc.plasmaStore.Get(ctx, key)
			hash := tmhash.Sum(append(input_utxo.TxHash, blockhash...))

			msg.Input0ConfirmSigs = CreateConfirmSig(hash, keys[tc.input0.input_index0], keys[input0_index1], tc.input0.input_index1 != -1)
			msg.Input1ConfirmSigs = CreateConfirmSig(hash, keys[input1_index0], keys[input1_index1], tc.input1.input_index1 != -1)
		}

		tx := GetTx(msg, keys[tc.input0.owner_index], keys[tc.input1.owner_index], !utils.ZeroAddress(msg.Owner1))
		txBytes, _ := rlp.EncodeToBytes(tx)

		dres := cc.DeliverTx(txBytes)

		require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

		// Retrieve utxo from context
		position := types.NewPlasmaPosition(uint64(index)+7, 0, 0, 0)
		utxo0 := cc.utxoMapper.GetUTXO(ctx, tc.newowner0.Bytes(), position)

		inputKeys := getInputKeys(cc.utxoMapper, tc.input0, tc.input1)
		txHash := tmhash.Sum(txBytes)
		expected := utxo.NewUTXOwithInputs(tc.newowner0.Bytes(), tc.amount0, "Ether", position, txHash, inputKeys)

		require.Equal(t, expected, utxo0, fmt.Sprintf("First UTXO did not get added to the utxo store correctly. Failed on test case: %d", index))

		if !utils.ZeroAddress(msg.Newowner1) {
			position = types.NewPlasmaPosition(uint64(index)+7, 0, 1, 0)
			utxo1 := cc.utxoMapper.GetUTXO(ctx, tc.newowner1.Bytes(), position)

			expected = utxo.NewUTXOwithInputs(tc.newowner1.Bytes(), tc.amount1, "Ether", position, txHash, inputKeys)

			require.Equal(t, expected, utxo1, fmt.Sprintf("Second UTXO did not get added to the utxo store correctly. Failed on test case: %d", index))
		}

		// Check that inputs were removed
		recovered := cc.utxoMapper.GetUTXO(ctx, msg.Owner0.Bytes(), tc.input0.position)
		require.False(t, recovered.Valid, fmt.Sprintf("first input was not removed from the utxo store. Failed on test case: %d", index))

		if !utils.ZeroAddress(msg.Owner1) {
			recovered = cc.utxoMapper.GetUTXO(ctx, msg.Owner1.Bytes(), tc.input1.position)
			require.False(t, recovered.Valid, fmt.Sprintf("second input was not removed from the utxo store. Failed on test case: %d", index))
		}

		cc.EndBlock(abci.RequestEndBlock{Height: 7 + int64(index)})
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

	InitTestChain(cc, utils.GenerateAddress(), addrs...)
	cc.Commit()

	for i := uint64(0); i < N; i++ {
		cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: int64(i + 1)}})
		position := types.NewPlasmaPosition(i+1, 0, 0, 0)
		storeInitUTXO(cc, position, addrs[i])
		cc.Commit()
		cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: N + 1}})

		msgs[i] = GenerateSimpleMsg(addrs[i], addrs[i], [4]uint64{i + 1, 0, 0, 0}, 100, 0)
		hash := tmhash.Sum(append([]byte("txhash"), []byte("merkle root")...))
		msgs[i].Input0ConfirmSigs = CreateConfirmSig(hash, keys[i], &ecdsa.PrivateKey{}, false)

		txs[i] = GetTx(msgs[i], keys[i], &ecdsa.PrivateKey{}, false)
		txBytes, _ := rlp.EncodeToBytes(txs[i])

		dres := cc.DeliverTx(txBytes)
		require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

	}
	cc.EndBlock(abci.RequestEndBlock{})
	cc.Commit()
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: N + 2}})
	ctx := cc.NewContext(false, abci.Header{Height: N + 1})

	// Retrieve and check UTXO from context
	for i := uint16(0); i < N; i++ {
		txBytes, _ := rlp.EncodeToBytes(txs[i])
		position := types.NewPlasmaPosition(N+1, i, 0, 0)
		actual := cc.utxoMapper.GetUTXO(ctx, addrs[i].Bytes(), position)

		inputKey := cc.utxoMapper.ConstructKey(addrs[i].Bytes(), types.NewPlasmaPosition(uint64(i+1), 0, 0, 0))
		expected := utxo.NewUTXOwithInputs(addrs[i].Bytes(), 100, "Ether", position, tmhash.Sum(txBytes), [][]byte{inputKey})
		expected.TxHash = tmhash.Sum(txBytes)

		require.Equal(t, expected, actual, fmt.Sprintf("UTXO %d did not get added to store correctly", i+1))

		position = types.NewPlasmaPosition(uint64(i)+1, 0, 0, 0)
		deposit := cc.utxoMapper.GetUTXO(ctx, addrs[i].Bytes(), position)
		require.False(t, deposit.Valid, fmt.Sprintf("utxo %d did not get removed correctly from the utxo store", i+1))
	}

	// send to different address
	for i := uint16(0); i < N; i++ {
		msgs[i].Blknum0 = N + 1
		msgs[i].Txindex0 = i
		msgs[i].DepositNum0 = 0

		blknumKey := make([]byte, binary.MaxVarintLen64)
		binary.PutUvarint(blknumKey, uint64(N+1))
		key := append(utils.RootHashPrefix, blknumKey...)
		blockhash := cc.plasmaStore.Get(ctx, key)

		txBytes, _ := rlp.EncodeToBytes(txs[i])
		txHash := tmhash.Sum(txBytes)

		hash := tmhash.Sum(append(txHash, blockhash...))
		msgs[i].Input0ConfirmSigs = CreateConfirmSig(hash, keys[i], &ecdsa.PrivateKey{}, false)

		msgs[i].Newowner0 = addrs[(i+1)%N]
		txs[i] = GetTx(msgs[i], keys[i], &ecdsa.PrivateKey{}, false)
		txBytes, _ = rlp.EncodeToBytes(txs[i])

		dres := cc.DeliverTx(txBytes)
		require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)
	}

	ctx = cc.NewContext(false, abci.Header{})

	// Retrieve and check UTXO from context
	for i := uint16(0); i < N; i++ {
		txBytes, _ := rlp.EncodeToBytes(txs[i])
		actual := cc.utxoMapper.GetUTXO(ctx, addrs[(i+1)%N].Bytes(), types.NewPlasmaPosition(N+2, i, 0, 0))

		inputKey := cc.utxoMapper.ConstructKey(addrs[i].Bytes(), types.NewPlasmaPosition(N+1, i, 0, 0))
		expected := utxo.NewUTXOwithInputs(addrs[(i+1)%N].Bytes(), 100, "Ether", types.NewPlasmaPosition(N+2, i, 0, 0), tmhash.Sum(txBytes), [][]byte{inputKey})

		require.Equal(t, expected, actual, fmt.Sprintf("UTXO %d did not get added to store correctly", i+1))

		input := cc.utxoMapper.GetUTXO(ctx, addrs[i].Bytes(), types.NewPlasmaPosition(N+1, i, 0, 0))
		require.False(t, input.Valid, fmt.Sprintf("UTXO %d  did not get removed from the utxo store correctly", i))
	}

}

func TestFee(t *testing.T) {
	cc := newChildChain()

	valPrivKey, _ := ethcrypto.GenerateKey()
	valAddr := utils.PrivKeyToAddress(valPrivKey)
	privKeys := make([]*ecdsa.PrivateKey, 2)
	addrs := make([]common.Address, 2)

	for i, _ := range privKeys {
		privKeys[i], _ = ethcrypto.GenerateKey()
		addrs[i] = utils.PrivKeyToAddress(privKeys[i])
	}

	InitTestChain(cc, valAddr, addrs...)
	position := types.NewPlasmaPosition(1, 0, 0, 0)
	storeInitUTXO(cc, position, addrs[0])
	position = types.NewPlasmaPosition(2, 0, 0, 0)
	storeInitUTXO(cc, position, addrs[1])

	cc.Commit()

	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 3}})
	// Create tx's with fees and deliver them in block 3
	msg1 := GenerateSimpleMsg(addrs[0], addrs[1], [4]uint64{1, 0, 0, 0}, 90, 10)
	msg2 := GenerateSimpleMsg(addrs[1], addrs[0], [4]uint64{2, 0, 0, 0}, 90, 10)

	blockHash := []byte("merkle root")
	txHash := []byte("txhash")
	hash := tmhash.Sum(append(txHash, blockHash...))
	msg1.Input0ConfirmSigs = CreateConfirmSig(hash, privKeys[0], &ecdsa.PrivateKey{}, false)
	msg2.Input0ConfirmSigs = CreateConfirmSig(hash, privKeys[1], &ecdsa.PrivateKey{}, false)

	tx1 := GetTx(msg1, privKeys[0], nil, false)
	tx2 := GetTx(msg2, privKeys[1], nil, false)

	check1 := cc.Check(tx1)
	check2 := cc.Check(tx2)

	res1 := cc.Deliver(tx1)
	res2 := cc.Deliver(tx2)

	// Assert checks pass
	require.Equal(t, sdk.CodeOK, sdk.CodeType(check1.Code), check1.Log)
	require.Equal(t, sdk.CodeOK, sdk.CodeType(check2.Code), check2.Log)

	// Assert delivering tx passes
	require.Equal(t, sdk.CodeOK, sdk.CodeType(res1.Code), res1.Log)
	require.Equal(t, sdk.CodeOK, sdk.CodeType(res2.Code), res2.Log)

	cc.EndBlock(abci.RequestEndBlock{})
	cc.Commit()

	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 4}})

	expectedPosition1 := types.NewPlasmaPosition(3, uint16(0), uint8(0), 0)
	expectedPosition2 := types.NewPlasmaPosition(3, uint16(1), uint8(0), 0)

	expectedValPosition := types.NewPlasmaPosition(3, uint16(1<<16-1), uint8(0), 0)

	ctx := cc.NewContext(false, abci.Header{Height: 4})

	utxo1 := cc.utxoMapper.GetUTXO(ctx, addrs[1].Bytes(), expectedPosition1)
	utxo2 := cc.utxoMapper.GetUTXO(ctx, addrs[0].Bytes(), expectedPosition2)

	valUTXO := cc.utxoMapper.GetUTXO(ctx, valAddr.Bytes(), expectedValPosition)

	// Check that users and validators have expected UTXO's
	require.Equal(t, uint64(90), utxo1.Amount, "UTXO1 does not have expected amount")
	require.Equal(t, uint64(90), utxo2.Amount, "UTXO2 does not have expected amount")
	require.Equal(t, uint64(20), valUTXO.Amount, "Validator fees did not get collected into UTXO correctly")

	// Check that validator can spend his fees as if they were a regular UTXO on sidechain
	valMsg := GenerateSimpleMsg(valAddr, addrs[0], [4]uint64{3, 1<<16 - 1, 0, 0}, 10, 10)

	valTx := GetTx(valMsg, valPrivKey, nil, false)

	valRes := cc.Deliver(valTx)

	require.Equal(t, sdk.CodeOK, sdk.CodeType(valRes.Code), valRes.Log)

	cc.EndBlock(abci.RequestEndBlock{})
	cc.Commit()

	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 5}})

	ctx = cc.NewContext(false, abci.Header{Height: 3})

	// Check that fee Amount gets reset between blocks. feeAmount for block 2 is 10 not 30.
	feeUTXO2 := cc.utxoMapper.GetUTXO(ctx, valAddr.Bytes(), types.NewPlasmaPosition(4, 1<<16-1, 0, 0))
	require.Equal(t, uint64(10), feeUTXO2.Amount, "Fee Amount on second block is incorrect")
}
