package app

import (
	"encoding/json"
	auth "github.com/FourthState/plasma-mvp-sidechain/auth"
	plasmaDB "github.com/FourthState/plasma-mvp-sidechain/db"
	"github.com/FourthState/plasma-mvp-sidechain/types"
	abci "github.com/tendermint/abci/types"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/go-amino"
	crypto "github.com/tendermint/go-crypto"
	tmtypes "github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
)

const (
	appName = "plasmaChildChain"
)

// Extended ABCI application
type ChildChain struct {
	*bam.BaseApp

	cdc *amino.Codec

	txIndex *uint16

	feeAmount *uint64

	// keys to access the substores
	capKeyMainStore *sdk.KVStoreKey

	// Manage addition and deletion of utxo's
	utxoMapper types.UTXOMapper
}

func NewChildChain(logger log.Logger, db dbm.DB) *ChildChain {
	cdc := MakeCodec()
	var app = &ChildChain{
		BaseApp:         bam.NewBaseApp(appName, cdc, logger, db),
		cdc:             cdc,
		txIndex:         new(uint16),
		feeAmount:       new(uint64),
		capKeyMainStore: sdk.NewKVStoreKey("main"),
	}

	// define the utxoMapper
	app.utxoMapper = plasmaDB.NewUTXOMapper(
		app.capKeyMainStore, // target store
		cdc,
	)

	// UTXOKeeper to adjust spending and recieving of utxo's
	UTXOKeeper := plasmaDB.NewUTXOKeeper(app.utxoMapper)
	app.Router().
		AddRoute("spend", auth.NewHandler(UTXOKeeper, app.txIndex))

	// set the BaseApp txDecoder to use txDecoder with RLP
	app.SetTxDecoder(app.txDecoder)

	app.MountStoresIAVL(app.capKeyMainStore)

	app.SetInitChainer(app.initChainer)
	app.SetEndBlocker(app.endBlocker)

	// NOTE: type AnteHandler func(ctx Context, tx Tx) (newCtx Context, result Result, abort bool)
	app.SetAnteHandler(auth.NewAnteHandler(app.utxoMapper, app.txIndex, app.feeAmount))

	err := app.LoadLatestVersion(app.capKeyMainStore)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

func (app *ChildChain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	// TODO is this now the whole genesis file?

	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	// load the accounts
	for _, gutxo := range genesisState.UTXOs {
		utxo := ToUTXO(gutxo)
		app.utxoMapper.AddUTXO(ctx, utxo)
	}

	// load the initial stake information
	return abci.ResponseInitChain{}
}

func (app *ChildChain) endBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	// reset txIndex and fee
	*app.txIndex = 0
	*app.feeAmount = 0

	return abci.ResponseEndBlock{}
}

// RLP decodes the txBytes to a BaseTx
func (app *ChildChain) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx = types.BaseTx{}

	err := rlp.DecodeBytes(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode(err.Error())
	}
	return tx, nil
}

func MakeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(PlasmaGenTx{}, "app/PlasmaGenTx", nil)
	types.RegisterAmino(cdc)
	crypto.RegisterAmino(cdc)
	return cdc
}

func (app *ChildChain) ExportAppStateJSON() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	// TODO: Implement
	// Currently non-functional, just enough to compile
	tx := types.BaseTx{}
	appState, err = app.cdc.MarshalJSONIndent(tx, "", "\t")
	validators = []tmtypes.GenesisValidator{}
	return appState, validators, err
}
