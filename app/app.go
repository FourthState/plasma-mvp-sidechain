package app

import (
	// TODO: Change to import from FourthState repo (currently not on there)
	types "plasma-mvp-sidechain/types" //points to a local package

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
	crypto "github.com/tendermint/go-crypto"

	"github.com/tendermint/go-amino" 
	//rlp "github.com/ethereum/go-ethereum/rlp" 
)

const (
	appName = "plasmaChildChain"
)

// Extended ABCI application
type ChildChain struct {
	*bam.BaseApp // Pointer to the Base App

	cdc *amino.Codec

	// keys to access the substores
	capKeyMainStore *sdk.KVStoreKey //capabilities key to access main store from multistore
	
	capKeySigStore *sdk.KVStoreKey //capabilities key to access confirm signature store from multistore
	//Not sure if this is needed
	capKeyIBCStore *sdk.KVStoreKey //capabilities key to access IBC Store from multistore

	// Manage addition and deletion of unspent utxo's
	utxoMapper types.UTXOMapper
}

func NewChildChain(logger log.Logger, db dbm.DB) *ChildChain {
	var app = &ChildChain{
		BaseApp:			bam.NewBaseApp(appName, logger, db),
		cdc: 				MakeCodec(),
		capKeyMainStore:	sdk.NewKVStoreKey("main"),
		capKeySigStore: 	sdk.NewKVStoreKey("sig"),
		capKeyIBCStore:  	sdk.NewKVStoreKey("ibc"),

	}

	// define the utxoMapper
	app.utxoMapper = types.NewUTXOMapper(
		app.capKeyMainStore, // target store
		app.capKeySigStore,
	)

	// UTXOKeeper to adjust spending and recieving of utxo's
	UTXOKeeper := types.NewUTXOKeeper(app.utxoMapper)
	app.Router().
		AddRoute("txs", types.NewHandler(UTXOKeeper))

	// initialize BaseApp
	// set the BaseApp txDecoder to use txDecoder with RLP
	app.SetTxDecoder(app.txDecoder)

	app.MountStoresIAVL(app.capKeyMainStore)

	// NOTE: type AnteHandler func(ctx Context, tx Tx) (newCtx Context, result Result, abort bool)
	//app.setAnteHandler(types.NewAnteHandler(app.utxoMapper))

	err := app.LoadLatestVersion(app.capKeyMainStore)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

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


func MakeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	types.RegisterAmino(cdc)   // Register SpendMsg, BaseTx
	crypto.RegisterAmino(cdc)
	return cdc
}