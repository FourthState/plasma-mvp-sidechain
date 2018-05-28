package app

import (
	auth "plasma-mvp-sidechain/auth"
	plasmaDB "plasma-mvp-sidechain/db"
	types "plasma-mvp-sidechain/types"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-amino"
	crypto "github.com/tendermint/go-crypto"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
	//rlp "github.com/ethereum/go-ethereum/rlp"
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

	txHash []byte
}

func NewChildChain(logger log.Logger, db dbm.DB) *ChildChain {
	cdc := MakeCodec()
	var app = &ChildChain{
		BaseApp:         bam.NewBaseApp(appName, cdc, logger, db),
		cdc:             cdc,
		txIndex:         new(uint16),
		feeAmount:       new(uint64),
		capKeyMainStore: sdk.NewKVStoreKey("main"),
		txHash:          nil,
	}

	// define the utxoMapper
	app.utxoMapper = plasmaDB.NewUTXOMapper(
		app.capKeyMainStore, // target store
		cdc,
	)

	// UTXOKeeper to adjust spending and recieving of utxo's
	UTXOKeeper := plasmaDB.NewUTXOKeeper(app.utxoMapper)
	app.Router().
		AddRoute("txs", auth.NewHandler(UTXOKeeper, app.txIndex))

	// set the BaseApp txDecoder to use txDecoder with RLP
	app.SetTxDecoder(app.txDecoder)

	app.MountStoresIAVL(app.capKeyMainStore)

	// NOTE: type AnteHandler func(ctx Context, tx Tx) (newCtx Context, result Result, abort bool)
	app.SetAnteHandler(auth.NewAnteHandler(app.utxoMapper, app.txIndex, app.feeAmount))

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
		return nil, sdk.ErrTxDecode("")
	}
	return tx, nil
}

func (app *ChildChain) endBlocker(ctx sdk.Context, req abci.RequestEndBlock) {
	header := ctx.BlockHeader()
	app.txHash = header.GetDataHash()
}

func MakeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	types.RegisterAmino(cdc)
	crypto.RegisterAmino(cdc)
	return cdc
}
