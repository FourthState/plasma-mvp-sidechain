package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	//clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	//"github.com/cosmos/sdk-application-tutorial/x/nameservice"
	//       "github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/query"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"
)

const (
	restName = "name"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// r.HandleFunc(fmt.Sprintf("/%s/names", storeName), buyNameHandler(cdc, cliCtx)).Methods("POST")
	// r.HandleFunc(fmt.Sprintf("/%s/names", storeName), setNameHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/balance/{%s}", restName), queryBalanceHandler(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprint("/health"), healthHandler(cdc, cliCtx)).Methods("GET")
}

func healthHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rest.PostProcessResponse(w, cdc, "healthy", cliCtx.Indent)
	}
}

func queryBalanceHandler(cdc *codec.Codec, ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		addr := ethcmn.HexToAddress(paramType)
		res, err := ctx.QuerySubspace(addr.Bytes(), "utxo")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		var utxos []store.UTXO

		utxo := store.UTXO{}

		for _, pair := range res {
			if err := rlp.DecodeBytes(pair.Value, &utxo); err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			if !utxo.Spent {
				utxos = append(utxos, utxo)
			}
		}

		rest.PostProcessResponse(w, cdc, utxos, ctx.Indent)
	}

}
