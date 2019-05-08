package app

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/eth"
	"github.com/FourthState/plasma-mvp-sidechain/handlers"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/query"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"io"
	"math/big"
	"os"
	"time"
)

const (
	appName = "plasmaMVP"
)

// Extended ABCI application
type PlasmaMVPChain struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	txIndex   uint16
	feeAmount *big.Int

	// persistent stores
	utxoStore   store.UTXOStore
	plasmaStore store.PlasmaStore

	// smart contract connection
	ethConnection *eth.Plasma

	/* Config */
	isOperator            bool // contract operator
	operatorPrivateKey    *ecdsa.PrivateKey
	operatorAddress       common.Address
	plasmaContractAddress common.Address
	blockCommitmentRate   time.Duration
	nodeURL               string // client that satisfies the web3 interface
	blockFinality         uint64 // presumed finality bound for the ethereum network
}

func NewPlasmaMVPChain(logger log.Logger, db dbm.DB, traceStore io.Writer, options ...func(*PlasmaMVPChain)) *PlasmaMVPChain {
	baseApp := baseapp.NewBaseApp(appName, logger, db, msgs.TxDecoder)
	cdc := MakeCodec()
	baseApp.SetCommitMultiStoreTracer(traceStore)

	utxoStoreKey := sdk.NewKVStoreKey("utxo")
	utxoStore := store.NewUTXOStore(utxoStoreKey)
	plasmaStoreKey := sdk.NewKVStoreKey("plasma")
	plasmaStore := store.NewPlasmaStore(plasmaStoreKey)
	app := &PlasmaMVPChain{
		BaseApp:   baseApp,
		cdc:       cdc,
		txIndex:   0,
		feeAmount: big.NewInt(0), // we do not use `utils.BigZero` because the feeAmount is going to be updated

		utxoStore:   utxoStore,
		plasmaStore: plasmaStore,
	}

	// set configs
	for _, option := range options {
		option(app)
	}

	// connect to remote client
	ethClient, err := eth.InitEthConn(app.nodeURL, logger)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	plasmaClient, err := eth.InitPlasma(app.plasmaContractAddress, ethClient, app.blockFinality, app.blockCommitmentRate, logger,
		app.isOperator, app.operatorPrivateKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app.ethConnection = plasmaClient

	// query for the operator address
	addr, err := plasmaClient.OperatorAddress()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app.operatorAddress = addr

	// Route spends to the handler
	nextTxIndex := func() uint16 {
		app.txIndex++
		return app.txIndex - 1
	}
	feeUpdater := func(amt *big.Int) sdk.Error {
		app.feeAmount = app.feeAmount.Add(app.feeAmount, amt)
		return nil
	}
	app.Router().AddRoute(msgs.SpendMsgRoute, handlers.NewSpendHandler(app.utxoStore, app.plasmaStore, nextTxIndex, feeUpdater))
	app.Router().AddRoute(msgs.IncludeDepositMsgRoute, handlers.NewDepositHandler(app.utxoStore, app.plasmaStore, nextTxIndex, plasmaClient))

	// custom queriers
	app.QueryRouter().
		AddRoute("utxo", query.NewUtxoQuerier(utxoStore)).
		AddRoute("plasma", query.NewPlasmaQuerier(plasmaStore))

	// Set the AnteHandler
	app.SetAnteHandler(handlers.NewAnteHandler(app.utxoStore, app.plasmaStore, plasmaClient))

	// set the rest of the chain flow
	app.SetEndBlocker(app.endBlocker)
	app.SetInitChainer(app.initChainer)

	// mount and load stores
	// IAVL store used by default. `fauxMerkleMode` defaults to false
	app.MountStores(utxoStoreKey, plasmaStoreKey)
	if err := app.LoadLatestVersion(utxoStoreKey); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return app
}

func (app *PlasmaMVPChain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := GenesisState{}
	if err := app.cdc.UnmarshalJSON(stateJSON, &genesisState); err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	// load the initial stake information
	return abci.ResponseInitChain{Validators: []abci.ValidatorUpdate{abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(genesisState.Validator.ConsPubKey),
		Power:  1,
	}}}
}

// Reset state at the end of each block
func (app *PlasmaMVPChain) endBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	// skip if the block is empty
	if app.txIndex == 0 {
		// try to commit any headers in the store
		app.ethConnection.CommitPlasmaHeaders(ctx, app.plasmaStore)
		return abci.ResponseEndBlock{}
	}

	tmBlockHeight := big.NewInt(ctx.BlockHeight())

	var header [32]byte
	copy(header[:], ctx.BlockHeader().DataHash)
	block := plasma.NewBlock(header, app.txIndex, app.feeAmount)
	plasmaBlockNum := app.plasmaStore.StoreBlock(ctx, tmBlockHeight, block)
	app.ethConnection.CommitPlasmaHeaders(ctx, app.plasmaStore)

	if app.feeAmount.Sign() != 0 {
		// add a utxo in the store with position 2^16-1
		utxo := store.UTXO{
			Position: plasma.NewPosition(plasmaBlockNum, 1<<16-1, 0, nil),
			Output:   plasma.NewOutput(app.operatorAddress, app.feeAmount),
			Spent:    false,
		}

		app.utxoStore.StoreUTXO(ctx, utxo)
	}

	app.txIndex = 0
	app.feeAmount = big.NewInt(0)

	return abci.ResponseEndBlock{}
}

func (app *PlasmaMVPChain) ExportAppStateJSON() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	// TODO: Implement
	// Currently non-functional, just enough to compile
	tx := msgs.SpendMsg{}
	appState, err = json.MarshalIndent(tx, "", "\t")
	validators = []tmtypes.GenesisValidator{}
	return appState, validators, err
}

func MakeCodec() *codec.Codec {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
