package app //change to app

//modeled after basecoinapp in cosmos/cosmos-sdk/examples
import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbm "github.com/tendermint/tmlibs/db"
	//crypto "github.com/tendermint/go-crypto"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/log"
	types "plasma-mvp-sidechain/types" 
	//"fmt" //for testing
	//"github.com/ethereum/go-ethereum/rlp"

)

const (
	appName = "plasmaChildChain"
)

// Extended ABCI application

type ChildChain struct {
	*bam.BaseApp // Pointer to the Base App
	// TODO: Add RLP here?

	// keys to access the substores
	capKeyMainStore *sdk.KVStoreKey //capabilities key to access main store from multistore
	//Not sure if this is needed
	capKeyIBCStore *sdk.KVStoreKey //capabilities key to access IBC Store from multistore

	// Manage addition and deletion of unspent utxo's 
	utxoMapper types.UtxoMapper
}

func NewChildChain(logger log.Logger, db dbm.DB) *ChildChain {
	var app = &ChildChain{
		BaseApp:			bam.NewBaseApp(appName, logger, db),
		capKeyMainStore:	sdk.NewKVStoreKey("main"),
		capKeyIBCStore:  	sdk.NewKVStoreKey("ibc"),
	}

	// define the utxoMapper
	app.utxoMapper = types.NewUTXOMapper(
		app.capKeyMainStore, // target store
		&types.BaseUTXO{},
	)

	// TODO: add handlers/router
	// UTXOKeeper to adjust spending and recieving of utxo's
	UTXOKeeper := types.NewUTXOKeeper(app.utxoMapper)
	app.Router().
		AddRoute("txs", types.NewHandler(UTXOKeeper))

	// initialize BaseApp
	// set the BaseApp txDecoder to use txDecoder with RLP
	app.SetTxDecoder(app.txDecoder)
	
	// TODO: set initChainer?
	// TO-UNDERSTAND: Not sure what mounting does yet
	app.MountStoresIAVL(app.capKeyMainStore)
	
	// TO-UNDERSTAND: What does ante handler do, do i need to make a new one
	//app.setAnteHandler()
	
	// TODO: set ante handler
	err := app.LoadLatestVersion(app.capKeyMainStore)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// TODO: change sdk.Tx to different transaction struct
func (app *ChildChain) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	// TODO: implement method
	return nil, nil
}

// TODO: Add initChainer?



// Current big idea TODO List:
// - Implement RLP Encoding/Decoding in app.go and tx.go
// - Implement AnteHandler
// - Write Basic Test Cases