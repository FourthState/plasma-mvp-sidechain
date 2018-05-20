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
func newUTXO(privA *ecdsa.PrivateKey, privB *ecdsa.PrivateKey, position Position) UTXO {
	addrA := ethcrypto.PubkeyToAddress(privA.PublicKey).Bytes()
	addrB := ethcrypto.PubkeyToAddress(privB.PublicKey).Bytes()
	confirmAddr := [2]crypto.Address{addrA, addrA}
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
	priv := crypto.GenPrivKeySecp256k1()
	sig := priv.Sign(msg.GetSignBytes())
	tx := NewBaseTx(msg, []sdk.StdSignature{{
		PubKey:    priv.PubKey(),
		Signature: sig,
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
	privA := crypto.GenPrivKeySecp256k1()
	privB := crypto.GenPrivKeySecp256k1()
	utxo1 := NewUTXO(privA, privB, position1)
	utxo2 := NewUTXO(privA, privB, position2)
	mapper.AddUTXO(ctx, utxo1)
	mapper.AddUTXO(ctx, utxo2)
	var msg = GenSpendMsgWithAddresses()
	msg.Owner1 = privB.PubKey().Address()
	msg.Owner2 = privB.PubKey().Address()
	priv := crypto.GenPrivKeySecp256k1()
	sig := priv.Sign(msg.GetSignBytes())
	tx := NewBaseTx(msg, []sdk.StdSignature{{
		PubKey:    priv.PubKey(),
		Signature: sig,
	}, {
		PubKey:    priv.PubKey(),
		Signature: sig,
	}})

	handler := NewAnteHandler(mapper, txIndex, feeAmount)
	_, res, abort := handler(ctx, tx)

	assert.Equal(t, true, abort, "Did not abort on wrong signer")
	assert.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(4)), res.Code, "Signer address does not match owner address")
}

//Tests a valid single input transaction
/*func TestValidSingleInput(t *testing.T) {
	ctx, mapper, txIndex, feeAmount := setup()

	privKeyA, _ := ethcrypto.GenerateKey() //Input Owner
	privKeyB, _ := ethcrypto.GenerateKey() //ConfirmSig owner and recipient

	position1 := Position{1, 0, 0}
	confirmSig, _ := ethcrypto.Sign(position1.GetSignBytes(), privKeyB)
	confirmSig1 := crypto.SignatureSecp256k1(confirmSig)
	confrimSigs := [2]crypto.Signature{confirmSig1, crypto.SignatureSecp256k1{}}

	//Single input
	var msg = SpendMsg{
		Blknum1:      1,
		Txindex1:     0,
		Oindex1:      0,
		Indenom1:     200,
		Owner1:       ethcrypto.PubkeyToAddress(privKeyA.PublicKey).Bytes(),
		ConfirmSigs1: confrimSigs,
		Blknum2:      0,
		Txindex2:     0,
		Oindex2:      0,
		Indenom2:     0,
		Owner2:       crypto.Address([]byte("")),
		ConfirmSigs2: confrimSigs,
		Newowner1:    ethcrypto.PubkeyToAddress(privKeyA.PublicKey).Bytes(),
		Denom1:       150,
		Newowner2:    ethcrypto.PubkeyToAddress(privKeyA.PublicKey).Bytes(),
		Denom2:       45,
		Fee:          5,
	}
	sig, _ := ethcrypto.Sign(msg.GetSignBytes(), privKeyA)
	sig1 := crypto.SignatureSecp256k1(sig)
	priv, _ := crypto.PrivKeyFromBytes(ethcrypto.FromECDSA(privKeyA))
	pk := priv.PubKey()
	tx := NewBaseTx(msg, []sdk.StdSignature {{
			PubKey: 	pk,
			Signature:	sig1,
		},})

	handler := NewAnteHandler(mapper, txIndex, feeAmount)
	_, res, abort := handler(ctx, tx)

	assert.Equal(t, true, abort, "Did not abort on utxo that does not exist")
	assert.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(6)), res.Code, res.Log)

	utxo1 := newUTXO(privKeyB, privKeyA, position1)
	mapper.AddUTXO(ctx, utxo1)

	_, res, abort = handler(ctx, tx)

	assert.Equal(t, false, abort, "aborted with valid transaction")
	assert.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1),sdk.CodeType(0)), res.Code, res.Log)
}*/
