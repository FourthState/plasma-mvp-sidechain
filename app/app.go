package app

import (
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	auth "github.com/FourthState/plasma-mvp-sidechain/auth"
	"github.com/FourthState/plasma-mvp-sidechain/eth"
	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/FourthState/plasma-mvp-sidechain/x/kvstore"
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	abci "github.com/tendermint/tendermint/abci/types"
	"io"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	ethcmn "github.com/ethereum/go-ethereum/common"
)

const (
	appName = "plasmaChildChain"
)

// Extended ABCI application
type ChildChain struct {
	*bam.BaseApp

	cdc *amino.Codec

	txIndex uint16

	feeAmount uint64

	// keys to access the substores
	capKeyMainStore *sdk.KVStoreKey

	capKeyPlasmaStore *sdk.KVStoreKey

	// Manage addition and deletion of utxo's
	utxoMapper utxo.Mapper

	plasmaStore kvstore.KVStore

	/* Validator Information */
	isValidator bool

	// Address that validator uses to collect fees
	validatorAddress ethcmn.Address

	// Private key for submitting blocks to rootchain
	validatorPrivKey *ecdsa.PrivateKey

	// Rootchain contract address
	rootchain ethcmn.Address

	// NodeURL for connecting to ethereum client
	nodeURL string

	// Minimum Fee a validator is willing to accept
	minFees uint64

	// Number of blocks required for a submitted block to be considered final
	blockFinality uint64

	ethConnection *eth.Plasma
}

func NewChildChain(logger log.Logger, db dbm.DB, traceStore io.Writer, options ...func(*ChildChain)) *ChildChain {
	cdc := MakeCodec()

	bapp := bam.NewBaseApp(appName, logger, db, txDecoder)
	bapp.SetCommitMultiStoreTracer(traceStore)

	var app = &ChildChain{
		BaseApp:           bapp,
		cdc:               cdc,
		txIndex:           0,
		feeAmount:         0,
		capKeyMainStore:   sdk.NewKVStoreKey("main"),
		capKeyPlasmaStore: sdk.NewKVStoreKey("plasma"),
	}

	for _, option := range options {
		option(app)
	}

	// define the utxoMapper
	app.utxoMapper = utxo.NewBaseMapper(
		app.capKeyMainStore, // target store
		cdc,
	)

	app.plasmaStore = kvstore.NewKVStore(
		app.capKeyPlasmaStore,
	)

	app.Router().
		AddRoute("spend", utxo.NewSpendHandler(app.utxoMapper, app.nextPosition))

	app.MountStoresIAVL(app.capKeyMainStore)
	app.MountStoresIAVL(app.capKeyPlasmaStore)

	app.SetInitChainer(app.initChainer)
	app.SetEndBlocker(app.endBlocker)

	// Set Ethereum connection
	client, err := eth.InitEthConn(app.nodeURL, bapp.Logger)
	if err != nil {
		panic(err)
	}

	plasmaClient, err := eth.InitPlasma(app.rootchain.Hex(), app.validatorPrivKey, client, app.BaseApp.Logger, app.blockFinality, app.isValidator)
	if err != nil {
		panic(err)
	}

	app.ethConnection = plasmaClient

	// NOTE: type AnteHandler func(ctx Context, tx Tx) (newCtx Context, result Result, abort bool)
	app.SetAnteHandler(auth.NewAnteHandler(app.utxoMapper, app.plasmaStore, app.feeUpdater, app.ethConnection))

	err = app.LoadLatestVersion(app.capKeyMainStore)
	if err != nil {
		cmn.Exit(err.Error())
	}
	err = app.LoadLatestVersion(app.capKeyPlasmaStore)
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
		app.utxoMapper.ReceiveUTXO(ctx, utxo)
	}

	app.validatorAddress = ethcmn.HexToAddress(genesisState.Validator.Address)

	// load the initial stake information
	return abci.ResponseInitChain{Validators: []abci.ValidatorUpdate{abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(genesisState.Validator.ConsPubKey),
		Power:  1,
	}}}
}

func (app *ChildChain) endBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	if app.feeAmount != 0 {
		position := types.PlasmaPosition{
			Blknum:     uint64(ctx.BlockHeight()),
			TxIndex:    uint16(1<<16 - 1),
			Oindex:     0,
			DepositNum: 0,
		}
		utxo := utxo.NewUTXO(app.validatorAddress.Bytes(), app.feeAmount, types.Denom, position)
		app.utxoMapper.ReceiveUTXO(ctx, utxo)
	}

	// reset txIndex and fee
	app.txIndex = 0
	app.feeAmount = 0

	blknumKey := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(blknumKey, uint64(ctx.BlockHeight()))
	key := append(utils.RootHashPrefix, blknumKey...)

	if ctx.BlockHeader().DataHash != nil {
		app.plasmaStore.Set(ctx, key, ctx.BlockHeader().DataHash)
	}
	return abci.ResponseEndBlock{}
}

// RLP decodes the txBytes to a BaseTx
func txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx = types.BaseTx{}

	err := rlp.DecodeBytes(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode(err.Error())
	}
	return tx, nil
}

// Return the next output position given ctx
// and secondary flag which indicates if it is for secondary outputs from single tx.
func (app *ChildChain) nextPosition(ctx sdk.Context, secondary bool) utxo.Position {
	if !secondary {
		app.txIndex++
		return types.NewPlasmaPosition(uint64(ctx.BlockHeight()), app.txIndex-1, 0, 0)
	}
	return types.NewPlasmaPosition(uint64(ctx.BlockHeight()), app.txIndex-1, 1, 0)
}

// Fee Updater passed into antehandler
func (app *ChildChain) feeUpdater(output []utxo.Output) sdk.Error {
	if len(output) != 1 || output[0].Denom != types.Denom {
		return utxo.ErrInvalidFee(2, "Fee must be paid in Eth")
	}
	app.feeAmount += output[0].Amount
	return nil
}

func MakeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(PlasmaGenTx{}, "app/PlasmaGenTx", nil)
	types.RegisterAmino(cdc)
	utxo.RegisterAmino(cdc)
	cryptoAmino.RegisterAmino(cdc)
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
