package app

import (
	// TODO: Change to import from FourthState repo (currently not on there)
	types "plasma-mvp-sidechain/types" //points to a local package

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	rlp "github.com/ethereum/go-ethereum/rlp"
)

const (
	appName = "plasmaChildChain" // Can be changed
)

// Extended ABCI application
type ChildChain struct {
	*bam.BaseApp // Pointer to the Base App

	// keys to access the substores
	capKeyMainStore *sdk.KVStoreKey //capabilities key to access main store from multistore
	//Not sure if this is needed
	capKeyIBCStore *sdk.KVStoreKey //capabilities key to access IBC Store from multistore

	// Manage addition and deletion of unspent utxo's
	utxoMapper types.UTXOMapper
}

func NewChildChain(logger log.Logger, db dbm.DB) *ChildChain {
	var app = &ChildChain{
		BaseApp:         bam.NewBaseApp(appName, logger, db),
		capKeyMainStore: sdk.NewKVStoreKey("main"),
		capKeyIBCStore:  sdk.NewKVStoreKey("ibc"),
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
	// NOTE: type AnteHandler func(ctx Context, tx Tx) (newCtx Context, result Result, abort bool)

	// TODO: implement types.newantehandler
	app.setAnteHandler(types.NewAnteHandler(app.utxoMapper))

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
	err := rlp.DecodeBytes(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode("").TraceCause(err, "")
	}
	return tx, nil
}
