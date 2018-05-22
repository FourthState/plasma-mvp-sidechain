package types

import (
	"crypto/ecdsa"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/tmlibs/log"
	"testing"
)

func setup() (sdk.Context, UTXOMapper, *uint16, *uint64) {
	ms, capKey := setupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil, log.NewNopLogger())
	mapper := NewUTXOMapper(capKey, MakeCodec())

	return ctx, mapper, new(uint16), new(uint64)

}

/// @param privA confirmSig Address
/// @param privB owner address
/// single input
func newSingleInputUTXO(privA *ecdsa.PrivateKey, privB *ecdsa.PrivateKey, position Position) UTXO {
	addrA := EthPrivKeyToSDKAddress(privA)
	addrB := EthPrivKeyToSDKAddress(privB)
	confirmAddr := [2]crypto.Address{addrA, crypto.Address([]byte(""))}
	return NewBaseUTXO(addrB, confirmAddr, 100, position)
}

func TestNoSigs(t *testing.T) {
	ctx, mapper, txIndex, feeAmount := setup()

	var msg = GenSpendMsgWithAddresses()
	tx := NewBaseTx(msg, []sdk.StdSignature{})

	handler := NewAnteHandler(mapper, txIndex, feeAmount)
	_, res, abort := handler(ctx, tx)

	assert.Equal(t, true, abort, "Did not abort with no signatures")
	assert.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(4)), res.Code, "Tx had processed with no signatures")
}

func TestNotEnoughSigs(t *testing.T) {
	ctx, mapper, txIndex, feeAmount := setup()

	var msg = GenSpendMsgWithAddresses()
	priv, _ := ethcrypto.GenerateKey()
	hash := ethcrypto.Keccak256(msg.GetSignBytes()) 
	sig, _ := ethcrypto.Sign(hash, priv)
	tx := NewBaseTx(msg, []sdk.StdSignature{{
		PubKey:     nil,
		Signature:  crypto.SignatureSecp256k1(sig),
	}})

	handler := NewAnteHandler(mapper, txIndex, feeAmount)
	_, res, abort := handler(ctx, tx)

	assert.Equal(t, true, abort, "Did not abort with incorrect number of signatures")
	assert.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(4)), res.Code, "Tx had processed with incorrect number of signatures")
}

func TestWrongSigner(t *testing.T) {
	ctx, mapper, txIndex, feeAmount := setup()

	position1 := Position{1000, 0, 0, 0}
	position2 := Position{1000, 1, 0, 0}
	privA, _ := ethcrypto.GenerateKey()
	privB, _ := ethcrypto.GenerateKey()
	utxo1 := NewUTXO(privA, privB, position1)
	utxo2 := NewUTXO(privA, privB, position2)
	mapper.AddUTXO(ctx, utxo1)
	mapper.AddUTXO(ctx, utxo2)
	var msg = GenSpendMsgWithAddresses()
	msg.Owner1 = EthPrivKeyToSDKAddress(privB)
	msg.Owner2 = EthPrivKeyToSDKAddress(privB)
	priv, _ := ethcrypto.GenerateKey()
	hash := ethcrypto.Keccak256(msg.GetSignBytes()) 
	sig, _ := ethcrypto.Sign(hash, priv)
	tx := NewBaseTx(msg, []sdk.StdSignature{{
		PubKey:    nil,
		Signature:  crypto.SignatureSecp256k1(sig),
	}, {
		PubKey:    nil,
		Signature:  crypto.SignatureSecp256k1(sig),
	}})

	handler := NewAnteHandler(mapper, txIndex, feeAmount)
	_, res, abort := handler(ctx, tx)

	assert.Equal(t, true, abort, "Did not abort on wrong signer")
	assert.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(4)), res.Code, "Signer address does not match owner address")
}

//Tests a valid single input transaction
func TestValidSingleInput(t *testing.T) {
	ctx, mapper, txIndex, feeAmount := setup()

	privKeyA, _ := ethcrypto.GenerateKey() //Input Owner
	privKeyB, _ := ethcrypto.GenerateKey() //ConfirmSig owner and recipient

	position1 := Position{1, 0, 0, 0}
	confirmSigHash := ethcrypto.Keccak256(position1.GetSignBytes())
	confirmSig, _ := ethcrypto.Sign(confirmSigHash, privKeyB)
	confirmSig1 := crypto.SignatureSecp256k1(confirmSig)
	confirmSigs := [2]crypto.Signature{confirmSig1, crypto.SignatureSecp256k1{}}

	//Single input
	var msg = SpendMsg{
		Blknum1:      1,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       EthPrivKeyToSDKAddress(privKeyA),
		ConfirmSigs1: confirmSigs,
		Blknum2:      0,
		Txindex2:     0,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       crypto.Address([]byte("")),
		ConfirmSigs2: confirmSigs,
		Newowner1:    EthPrivKeyToSDKAddress(privKeyA),
		Denom1:       150,
		Newowner2:    EthPrivKeyToSDKAddress(privKeyA),
		Denom2:       45,
		Fee:          5,
	}
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ := ethcrypto.Sign(hash, privKeyA)
	sig1 := crypto.SignatureSecp256k1(sig)
	tx := NewBaseTx(msg, []sdk.StdSignature{{
			PubKey: 	nil,
			Signature:	sig1,
		}})

	handler := NewAnteHandler(mapper, txIndex, feeAmount)
	_, res, abort := handler(ctx, tx)

	assert.Equal(t, true, abort, "Did not abort on utxo that does not exist")
	assert.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(6)), res.Code, res.Log)

	utxo1 := newSingleInputUTXO(privKeyB, privKeyA, position1) 
	mapper.AddUTXO(ctx, utxo1)

	_, res, abort = handler(ctx, tx)

	assert.Equal(t, false, abort, "aborted with valid transaction")
	assert.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1),sdk.CodeType(0)), res.Code, res.Log)
}
