package app 

import (
	// TODO: Change to import from FourthState repo (currently not on there)
	types "plasma-mvp-sidechain/types" //points to a local package

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbm "github.com/tendermint/tmlibs/db"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/log"
	crypto "github.com/tendermint/go-crypto"

	"github.com/tendermint/go-amino" // Not necessary once switched to RLP
	//rlp "github.com/ethereum/go-ethereum/rlp" // TODO: Change from amino to RLP


)

const (
	appName = "plasmaChildChain" // Can be changed
)

// Extended ABCI application
type ChildChain struct {
	*bam.BaseApp // Pointer to the Base App

	cdc *amino.Codec

	// keys to access the substores
	capKeyMainStore *sdk.KVStoreKey //capabilities key to access main store from multistore
	//Not sure if this is needed
	capKeyIBCStore *sdk.KVStoreKey //capabilities key to access IBC Store from multistore

	// Manage addition and deletion of unspent utxo's 
	utxoMapper types.UTXOMapper
}

func NewChildChain(logger log.Logger, db dbm.DB) *ChildChain {
	var app = &ChildChain{
		BaseApp:			bam.NewBaseApp(appName, logger, db),
		cdc: MakeCodec(),
		capKeyMainStore:	sdk.NewKVStoreKey("main"),
		capKeyIBCStore:  	sdk.NewKVStoreKey("ibc"),
	}

	// define the utxoMapper
	app.utxoMapper = types.NewUTXOMapper(
		app.capKeyMainStore, // target store
		// MYNOTE: may need to change proto
		&types.BaseUTXOHolder{}, // UTXOHolder is a struct that holds BaseUTXO's
							 // BaseUTXO implemented UTXO interface
	)

	// TODO: add handlers/router
	// UTXOKeeper to adjust spending and recieving of utxo's
	UTXOKeeper := types.NewUTXOKeeper(app.utxoMapper)
	app.Router().
		AddRoute("txs", types.NewHandler(UTXOKeeper))

	// initialize BaseApp
	// set the BaseApp txDecoder to use txDecoder with RLP
	app.SetTxDecoder(app.txDecoder)
	
	// TO-UNDERSTAND: Not sure what mounting does yet
	app.MountStoresIAVL(app.capKeyMainStore)
	
	// TODO: Make ante handler
	//
	err := app.LoadLatestVersion(app.capKeyMainStore)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// TODO: change sdk.Tx to different transaction struct
func (app *ChildChain) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	// TODO: implement method with RLP
	var tx = types.BaseTx{}
	// BaseTx is struct for Msg wrapped with authentication data
	err := app.cdc.UnmarshalBinary(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode("").TraceCause(err, "")
	}
	return tx, nil
}

// TODO: Add initChainer?

func MakeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	types.RegisterAmino(cdc)   // Register bank.[SendMsg,IssueMsg] types.
	// crypto.RegisterWire(cdc)
	cdc.RegisterConcrete(crypto.PubKey{}, "go-crypto/PubKey", nil)
	cdc.RegisterConcrete(crypto.PrivKey{}, "go-crypto/PrivKey", nil)
	cdc.RegisterConcrete(crypto.Signature{}, "go-crypto/Signature", nil)
	cdc.RegisterConcrete(sdk.StdSignature{}, "sdk/StdSignature", nil)
	cdc.RegisterInterface((*crypto.PubKeyInner)(nil), nil)
	cdc.RegisterConcrete(crypto.PubKeySecp256k1{}, "go-crypto/PubKeySecpk1", nil)
	cdc.RegisterConcrete(crypto.SignatureSecp256k1{}, "go-crypto/SignatureSecpk1", nil)
	cdc.RegisterInterface((*crypto.SignatureInner)(nil), nil)
	return cdc
}



// Current TODO List:
// - Implement RLP Encoding/Decoding in app.go and tx.go
// - Implement AnteHandler
// - Write Basic Test Cases
