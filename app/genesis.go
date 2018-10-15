package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	crypto "github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/ethereum/go-ethereum/common"
)

// State to Unmarshal
type GenesisState struct {
	Validator GenesisValidator `json: validator`
	UTXOs     []GenesisUTXO    `json:"UTXOs"`
}

type GenesisValidator struct {
	ConsPubKey crypto.PubKey `json: validator_pubkey`
	Address    string        `json: fee_address`
}

type GenesisUTXO struct {
	Address  string
	Denom    string
	Position [4]string
}

func NewGenesisUTXO(addr string, amount string, position [4]string) GenesisUTXO {
	utxo := GenesisUTXO{
		Address:  addr,
		Denom:    amount,
		Position: position,
	}
	return utxo
}

func ToUTXO(gutxo GenesisUTXO) utxo.UTXO {
	// Any failed str conversion defaults to 0
	addr := common.HexToAddress(gutxo.Address)
	amount, _ := strconv.ParseUint(gutxo.Denom, 10, 64)
	utxo := &types.BaseUTXO{
		InputAddresses: [2]common.Address{addr, common.Address{}},
		Address:        addr,
		Amount:         amount,
		Denom:          types.Denom,
	}
	blkNum, _ := strconv.ParseUint(gutxo.Position[0], 10, 64)
	txIndex, _ := strconv.ParseUint(gutxo.Position[1], 10, 16)
	oIndex, _ := strconv.ParseUint(gutxo.Position[2], 10, 8)
	depNum, _ := strconv.ParseUint(gutxo.Position[3], 10, 64)

	position := types.NewPlasmaPosition(blkNum, uint16(txIndex), uint8(oIndex), depNum)
	utxo.SetPosition(position)
	return utxo
}

var (
	flagAddress    = "address"
	flagClientHome = "home-client"
	flagOWK        = "owk"

	// UTXO amount awarded
	freeEtherVal = int64(100)

	// default home directories for expected binaries
	DefaultCLIHome  = os.ExpandEnv("$HOME/.plasmacli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.plasmad")
)

// get app init parameters for server init command
func PlasmaAppInit() server.AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)

	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(flagAddress, "", "address, required")
	fsAppGenTx.String(flagClientHome, DefaultCLIHome,
		"home directory for the client, used for key generation")
	fsAppGenTx.Bool(flagOWK, false, "overwrite the accounts created")

	return server.AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         PlasmaAppGenTx,
		AppGenState:      PlasmaAppGenStateJSON,
	}
}

// simple genesis tx
type PlasmaGenTx struct {
	// currently takes address as string because unmarshaling Ether address fails
	Address string `json:"address"`
}

// Generate a gaia genesis transaction with flags
func PlasmaAppGenTx(cdc *wire.Codec, pk crypto.PubKey, gentTxConfig config.GenTx) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	addrString := viper.GetString(flagAddress)
	overwrite := viper.GetBool(flagOWK)

	bz, err := cdc.MarshalJSON("success")
	cliPrint = json.RawMessage(bz)
	appGenTx, _, validator, err = PlasmaAppGenTxNF(cdc, pk, addrString, overwrite)
	return
}

// Generate a gaia genesis transaction without flags
func PlasmaAppGenTxNF(cdc *wire.Codec, pk crypto.PubKey, addr string, overwrite bool) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {

	var bz []byte
	plasmaGenTx := PlasmaGenTx{
		Address: addr,
	}
	bz, err = wire.MarshalJSONIndent(cdc, plasmaGenTx)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)

	validator = tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  1,
	}
	return
}

// Create the core parameters for genesis initialization for gaia
// note that the pubkey input is this machines pubkey
func PlasmaAppGenState(cdc *wire.Codec, appGenTxs []json.RawMessage) (genesisState GenesisState, err error) {

	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}

	// get genesis flag account information
	genUTXO := make([]GenesisUTXO, len(appGenTxs))
	for i, appGenTx := range appGenTxs {

		var genTx PlasmaGenTx
		err = cdc.UnmarshalJSON(appGenTx, &genTx)
		if err != nil {
			return
		}

		genUTXO[i] = NewGenesisUTXO(genTx.Address, "100", [4]string{"0", "0", "0", fmt.Sprintf("%d", i+1)})
	}

	// create the final app state
	genesisState = GenesisState{
		UTXOs: genUTXO,
	}
	return
}

// PlasmaAppGenState but with JSON
func PlasmaAppGenStateJSON(cdc *wire.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {

	// create the final app state
	genesisState, err := PlasmaAppGenState(cdc, appGenTxs)
	if err != nil {
		return nil, err
	}
	appState, err = wire.MarshalJSONIndent(cdc, genesisState)
	return
}
