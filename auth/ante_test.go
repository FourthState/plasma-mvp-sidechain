package auth

import (
	"crypto/ecdsa"
	db "github.com/FourthState/plasma-mvp-sidechain/db"
	types "github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/log"
	"testing"
)

func setup() (sdk.Context, types.UTXOMapper, *uint16, *uint64) {
	ms, capKey := db.SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil, log.NewNopLogger())
	mapper := db.NewUTXOMapper(capKey, db.MakeCodec())

	return ctx, mapper, new(uint16), new(uint64)

}

/// @param privA confirmSig Address
/// @param privB owner address
/// single input
func newSingleInputUTXO(privA *ecdsa.PrivateKey, privB *ecdsa.PrivateKey, position types.Position) types.UTXO {
	addrA := utils.PrivKeyToAddress(privA)
	addrB := utils.PrivKeyToAddress(privB)
	confirmAddr := [2]common.Address{addrA, common.Address{}}
	return types.NewBaseUTXO(addrB, confirmAddr, 100, position)
}

/// @param privA confirmSig Address
/// @param privB owner address
/// two inputs
func newUTXO(privA *ecdsa.PrivateKey, privB *ecdsa.PrivateKey, position types.Position) types.UTXO {
	addrA := utils.PrivKeyToAddress(privA)
	addrB := utils.PrivKeyToAddress(privB)
	confirmAddr := [2]common.Address{addrA, addrA}
	return types.NewBaseUTXO(addrB, confirmAddr, 100, position)
}

func GenBasicSpendMsg() types.SpendMsg {
	// Creates Basic Spend Msg with no owners or recipients
	confirmSigs := [2]types.Signature{types.Signature{}, types.Signature{}}
	return types.SpendMsg{
		Blknum1:      1000,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       common.Address{},
		ConfirmSigs1: confirmSigs,
		Blknum2:      1000,
		Txindex2:     1,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       common.Address{},
		ConfirmSigs2: confirmSigs,
		Newowner1:    common.Address{},
		Denom1:       150,
		Newowner2:    common.Address{},
		Denom2:       50,
		Fee:          0,
	}
}

func GenSpendMsgWithAddresses() types.SpendMsg {
	// Creates Basic Spend Msg with owners and recipients
	confirmSigs := [2]types.Signature{types.Signature{}, types.Signature{}}
	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()

	return types.SpendMsg{
		Blknum1:      1000,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs1: confirmSigs,
		Blknum2:      1000,
		Txindex2:     1,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs2: confirmSigs,
		Newowner1:    utils.PrivKeyToAddress(privKeyB),
		Denom1:       150,
		Newowner2:    utils.PrivKeyToAddress(privKeyB),
		Denom2:       50,
		Fee:          0,
	}
}

// No signatures are provided
func TestNoSigs(t *testing.T) {
	ctx, mapper, txIndex, feeAmount := setup()

	var msg = GenSpendMsgWithAddresses()
	tx := types.NewBaseTx(msg, []types.Signature{})

	handler := NewAnteHandler(mapper, txIndex, feeAmount)
	_, res, abort := handler(ctx, tx)

	require.Equal(t, true, abort, "did not abort with no signatures")
	require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(4)), res.Code, "tx had processed with no signatures")
}

// The wrong amount of signatures are provided
func TestNotEnoughSigs(t *testing.T) {
	ctx, mapper, txIndex, feeAmount := setup()

	var msg = GenSpendMsgWithAddresses()
	priv, _ := ethcrypto.GenerateKey()
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ := ethcrypto.Sign(hash, priv)
	tx := types.NewBaseTx(msg, []types.Signature{types.Signature{sig}})

	handler := NewAnteHandler(mapper, txIndex, feeAmount)
	_, res, abort := handler(ctx, tx)

	require.Equal(t, true, abort, "did not abort with incorrect number of signatures")
	require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(4)), res.Code, "tx had processed with incorrect number of signatures")
}

// The transaction is not signed by the utxo owner
func TestWrongSigner(t *testing.T) {
	ctx, mapper, txIndex, feeAmount := setup()

	// Generate input utxos
	position1 := types.Position{1000, 0, 0, 0}
	position2 := types.Position{1000, 1, 0, 0}
	privA, _ := ethcrypto.GenerateKey()
	privB, _ := ethcrypto.GenerateKey()
	utxo1 := NewUTXO(privA, privB, position1)
	utxo2 := NewUTXO(privA, privB, position2)
	mapper.AddUTXO(ctx, utxo1)
	mapper.AddUTXO(ctx, utxo2)

	// Signature by non owner
	var msg = GenSpendMsgWithAddresses()
	msg.Owner1 = utils.PrivKeyToAddress(privB)
	msg.Owner2 = utils.PrivKeyToAddress(privB)
	priv, _ := ethcrypto.GenerateKey()
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, err := ethcrypto.Sign(hash, priv)
	require.NoError(t, err)
	tx := types.NewBaseTx(msg, []types.Signature{types.Signature{sig}, types.Signature{sig}})

	handler := NewAnteHandler(mapper, txIndex, feeAmount)
	_, res, abort := handler(ctx, tx)

	require.Equal(t, true, abort, "did not abort on wrong signer")
	require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(4)), res.Code, "signer address does not match owner address")
}

// Tests a valid single input transaction
func TestValidSingleInput(t *testing.T) {
	ctx, mapper, txIndex, feeAmount := setup()

	privKeyA, _ := ethcrypto.GenerateKey() //Input Owner
	privKeyB, _ := ethcrypto.GenerateKey() //ConfirmSig owner and recipient

	position1 := types.Position{1, 0, 0, 0}
	confirmSigHash := ethcrypto.Keccak256(position1.GetSignBytes())
	confirmSig, err := ethcrypto.Sign(confirmSigHash, privKeyB)
	require.NoError(t, err)
	confirmSigs := [2]types.Signature{types.Signature{confirmSig}, types.Signature{}}

	//Single input
	var msg = types.SpendMsg{
		Blknum1:      1,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs1: confirmSigs,
		Blknum2:      0,
		Txindex2:     0,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       common.Address{},
		ConfirmSigs2: confirmSigs,
		Newowner1:    utils.PrivKeyToAddress(privKeyA),
		Denom1:       50,
		Newowner2:    utils.PrivKeyToAddress(privKeyA),
		Denom2:       45,
		Fee:          5,
	}

	// Sign transaction
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, err := ethcrypto.Sign(hash, privKeyA)
	require.NoError(t, err)
	tx := types.NewBaseTx(msg, []types.Signature{types.Signature{sig}})

	handler := NewAnteHandler(mapper, txIndex, feeAmount)
	_, res, abort := handler(ctx, tx)

	require.Equal(t, true, abort, "did not abort on utxo that does not exist")
	require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(6)), res.Code, res.Log)

	utxo1 := newSingleInputUTXO(privKeyB, privKeyA, position1)
	mapper.AddUTXO(ctx, utxo1)

	_, res, abort = handler(ctx, tx)

	require.Equal(t, false, abort, "aborted with valid transaction")
	require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(0)), res.Code, res.Log)
}

// Tests a valid transaction
func TestValidTransaction(t *testing.T) {
	ctx, mapper, txIndex, feeAmount := setup()

	privKeyA, _ := ethcrypto.GenerateKey() //Input Owner
	privKeyB, _ := ethcrypto.GenerateKey() //ConfirmSig owner and recipient

	// Generate valid inputs
	position1 := types.Position{1, 0, 0, 0}
	position2 := types.Position{1, 1, 0, 0}
	confirmSigHash1 := ethcrypto.Keccak256(position1.GetSignBytes())
	confirmSigHash2 := ethcrypto.Keccak256(position2.GetSignBytes())
	confirmSig1, err := ethcrypto.Sign(confirmSigHash1, privKeyB)
	require.NoError(t, err)
	confirmSig2, err := ethcrypto.Sign(confirmSigHash2, privKeyB)
	require.NoError(t, err)
	confirmSigs1 := [2]types.Signature{types.Signature{confirmSig1}, types.Signature{confirmSig1}}
	confirmSigs2 := [2]types.Signature{types.Signature{confirmSig2}, types.Signature{confirmSig2}}

	//Single input
	var msg = types.SpendMsg{
		Blknum1:      1,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs1: confirmSigs1,
		Blknum2:      1,
		Txindex2:     1,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs2: confirmSigs2,
		Newowner1:    utils.PrivKeyToAddress(privKeyB),
		Denom1:       150,
		Newowner2:    utils.PrivKeyToAddress(privKeyB),
		Denom2:       45,
		Fee:          5,
	}
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, err := ethcrypto.Sign(hash, privKeyA)
	require.NoError(t, err)
	tx := types.NewBaseTx(msg, []types.Signature{types.Signature{sig}, types.Signature{sig}})

	handler := NewAnteHandler(mapper, txIndex, feeAmount)
	_, res, abort := handler(ctx, tx)

	require.Equal(t, true, abort, "did not abort on utxo that does not exist")
	require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(6)), res.Code, res.Log)

	utxo1 := newUTXO(privKeyB, privKeyA, position1)
	utxo2 := newUTXO(privKeyB, privKeyA, position2)
	mapper.AddUTXO(ctx, utxo1)
	mapper.AddUTXO(ctx, utxo2)

	_, res, abort = handler(ctx, tx)

	require.Equal(t, false, abort, "aborted with valid transaction")
	require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(0)), res.Code, res.Log)
}

// Check for double input that ante handler will
// prevent any malformed transactions with unequal
// input output fee balance from being spent
func TestDenomEquality(t *testing.T) {
	ctx, mapper, txIndex, feeAmount := setup()

	privKeyA, _ := ethcrypto.GenerateKey() //Input Owner
	privKeyB, _ := ethcrypto.GenerateKey() //ConfirmSig owner and recipient

	// Generate valid inputs
	position1 := types.Position{1, 0, 0, 0}
	position2 := types.Position{1, 1, 0, 0}
	confirmSigHash1 := ethcrypto.Keccak256(position1.GetSignBytes())
	confirmSigHash2 := ethcrypto.Keccak256(position2.GetSignBytes())
	confirmSig1, err := ethcrypto.Sign(confirmSigHash1, privKeyB)
	require.NoError(t, err)
	confirmSig2, err := ethcrypto.Sign(confirmSigHash2, privKeyB)
	require.NoError(t, err)
	confirmSigs1 := [2]types.Signature{types.Signature{confirmSig1}, types.Signature{confirmSig1}}
	confirmSigs2 := [2]types.Signature{types.Signature{confirmSig2}, types.Signature{confirmSig2}}

	//Single input
	var msg = types.SpendMsg{
		Blknum1:      1,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs1: confirmSigs1,
		Blknum2:      1,
		Txindex2:     1,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs2: confirmSigs2,
		Newowner1:    utils.PrivKeyToAddress(privKeyB),
		Denom1:       150,
		Newowner2:    utils.PrivKeyToAddress(privKeyB),
		Denom2:       50,
		Fee:          5,
	}
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, err := ethcrypto.Sign(hash, privKeyA)
	require.NoError(t, err)
	tx := types.NewBaseTx(msg, []types.Signature{types.Signature{sig}, types.Signature{sig}})

	handler := NewAnteHandler(mapper, txIndex, feeAmount)
	_, res, abort := handler(ctx, tx)

	require.Equal(t, true, abort, "did not abort on utxo that does not exist")
	require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(6)), res.Code, res.Log)

	utxo1 := newUTXO(privKeyB, privKeyA, position1)
	utxo2 := newUTXO(privKeyB, privKeyA, position2)
	mapper.AddUTXO(ctx, utxo1)
	mapper.AddUTXO(ctx, utxo2)

	_, res, abort = handler(ctx, tx)

	require.Equal(t, true, abort, "did not abort with invalid transaction")
	require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(4)), res.Code, res.Log)
}
