package app

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/eth"
	"github.com/FourthState/plasma-mvp-sidechain/handlers"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/store/query"
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
	outputStore store.OutputStore
	blockStore  store.BlockStore

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

	outputStoreKey := sdk.NewKVStoreKey("outputs")
	outputStore := store.NewOutputStore(outputStoreKey)
	blockStoreKey := sdk.NewKVStoreKey("block")
	blockStore := store.NewBlockStore(blockStoreKey)

	app := &PlasmaMVPChain{
		BaseApp:   baseApp,
		cdc:       cdc,
		txIndex:   0,
		feeAmount: big.NewInt(0), // we do not use `utils.BigZero` because the feeAmount is going to be updated

		blockStore:  blockStore,
		outputStore: outputStore,
	}

	// set configs
	for _, option := range options {
		option(app)
	}

	// connect to remote client
	eth.SetLogger(logger)
	ethClient, err := eth.InitEthConn(app.nodeURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	plasmaClient, err := eth.InitPlasma(app.plasmaContractAddress, ethClient, app.blockFinality)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if app.isOperator {
		plasmaClient, err = plasmaClient.WithOperatorSession(app.operatorPrivateKey, app.blockCommitmentRate)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app.ethConnection = plasmaClient

	// query for the operator address
	addr, err := plasmaClient.OperatorAddress()
	if err != nil {
		logger.Error("unable to query the contract for the operator address")
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
	app.Router().AddRoute(msgs.SpendMsgRoute, handlers.NewSpendHandler(app.outputStore, app.blockStore, nextTxIndex, feeUpdater))
	app.Router().AddRoute(msgs.IncludeDepositMsgRoute, handlers.NewDepositHandler(app.outputStore, app.blockStore, nextTxIndex, plasmaClient))

	// custom queriers
	app.QueryRouter().
		AddRoute(store.QueryOutputStore, query.NewOutputQuerier(outputStore)).
		AddRoute(store.QueryBlockStore, query.NewBlockQuerier(blockStore))

	// Set the AnteHandler
	app.SetAnteHandler(handlers.NewAnteHandler(app.outputStore, app.blockStore, plasmaClient))

	// set the rest of the chain flow
	app.SetEndBlocker(app.endBlocker)
	app.SetInitChainer(app.initChainer)

	// mount and load stores
	// IAVL store used by default. `fauxMerkleMode` defaults to false
	app.MountStores(outputStoreKey, blockStoreKey)
	if err := app.LoadLatestVersion(outputStoreKey); err != nil {
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
		app.ethConnection.CommitPlasmaHeaders(ctx, app.blockStore)
		return abci.ResponseEndBlock{}
	}

	tmBlockHeight := uint64(ctx.BlockHeight())

	var header [32]byte
	copy(header[:], ctx.BlockHeader().DataHash)
	block := plasma.NewBlock(header, app.txIndex, app.feeAmount)
	app.blockStore.StoreBlock(ctx, tmBlockHeight, block)
	app.ethConnection.CommitPlasmaHeaders(ctx, app.blockStore)

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
