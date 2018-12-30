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
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"io"
	"math/big"
	"os"
)

const (
	appName = "plasmaMVP"
)

// Extended ABCI application
type PlasmaMVPChain struct {
	*baseapp.BaseApp

	txIndex   uint16
	feeAmount *big.Int
	numTxns   uint16

	// persistent stores
	utxoStore      store.UTXOStore
	plasmaStore    store.PlasmaStore
	fauxMerkleMode bool

	// smart contract connection
	ethConnection eth.Plasma

	/* Config */
	isOperator            bool // contract operator
	operatorPrivateKey    *ecdsa.PrivateKey
	plasmaContractAddress common.Address
	nodeURL               string // client that satisfies the web3 interface
	blockFinality         uint64 // presumed finality bound for the ethereum network
}

func NewPlasmaMVPChain(logger log.Logger, db dbm.DB, traceStore io.Writer, options ...func(*PlasmaMVPChain)) *PlasmaMVPChain {
	baseApp := baseapp.NewBaseApp(appName, logger, db, msgs.TxDecoder)
	baseApp.SetCommitMultiStoreTracer(traceStore)

	utxoStoreKey := sdk.NewKVStoreKey("utxo")
	plasmaStoreKey := sdk.NewKVStoreKey("plasma")
	app := &PlasmaMVPChain{
		BaseApp:   baseApp,
		txIndex:   0,
		numTxns:   0,
		feeAmount: big.NewInt(0), // we do not use `utils.BigZero` because the feeAmount is going to be updated

		utxoStore:   store.NewUTXOStore(utxoStoreKey),
		plasmaStore: store.NewPlasmaStore(plasmaStoreKey),
	}

	// set configs
	for _, option := range options {
		option(app)
	}

	// connect to remote client
	conn, err := eth.InitEthConn(app.nodeURL, logger)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	plasmaConn, err := eth.InitPlasma(app.plasmaContractAddress, app.operatorPrivateKey, conn, app.blockFinality, logger)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app.ethConnection = plasmaConn

	// mount and load stores
	// IAVL store used by default. `fauxMerkleMode` defaults to false
	app.MountStores(utxoStoreKey, plasmaStoreKey)
	if err := app.LoadLatestVersion(utxoStoreKey); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := app.LoadLatestVersion(utxoStoreKey); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Route spends to the handler
	nextTxIndex := func() uint16 {
		app.txIndex++
		return app.txIndex - 1
	}
	app.Router().AddRoute(msgs.SpendMsgRoute, handlers.NewSpendHandler(app.utxoStore, nextTxIndex))

	// Set the AnteHandler
	feeUpdater := func(amt *big.Int) sdk.Error {
		app.feeAmount = app.feeAmount.Add(app.feeAmount, amt)
		return nil
	}
	app.SetAnteHandler(handlers.NewAnteHandler(app.utxoStore, app.plasmaStore, feeUpdater, plasmaConn))

	// set the rest of the chain flow
	app.SetEndBlocker(app.endBlocker)
	app.SetInitChainer(app.initChainer)

	return app
}

func (app *PlasmaMVPChain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	// TODO is this now the whole genesis file?

	var genesisState GenesisState
	err := json.Unmarshal(stateJSON, &genesisState)
	if err != nil {
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
	blockHeight := big.NewInt(ctx.BlockHeight())
	if app.feeAmount.Sign() != 0 {
		// add a utxo in the store with position 2^16-1
		utxo := store.UTXO{
			Position: plasma.NewPosition(blockHeight, 1<<16-1, 0, nil),
			Output:   plasma.NewOutput(crypto.PubkeyToAddress(app.operatorPrivateKey.PublicKey), app.feeAmount),
			Spent:    false,
		}

		app.utxoStore.StoreUTXO(ctx, utxo)
	}

	var header [32]byte
	copy(header[:], ctx.BlockHeader().DataHash)
	block := plasma.NewBlock(header, app.numTxns, app.feeAmount)
	app.plasmaStore.StoreBlock(ctx, blockHeight, block)

	app.ethConnection.SubmitBlock(block)

	app.txIndex = 0
	app.feeAmount = big.NewInt(0)
	app.numTxns = app.txIndex

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
