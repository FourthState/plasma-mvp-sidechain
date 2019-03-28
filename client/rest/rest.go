package rest

import (
	//"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"net/http"
	//clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	//"github.com/cosmos/sdk-application-tutorial/x/nameservice"
	//       "github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/query"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"
	//"io/ioutil"
	"math/big"
)

const (
	restName = "name"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// r.HandleFunc(fmt.Sprintf("/%s/names", storeName), buyNameHandler(cdc, cliCtx)).Methods("POST")
	// r.HandleFunc(fmt.Sprintf("/%s/names", storeName), setNameHandler(cdc, cliCtx)).Methods("PUT")
	// r.HandleFunc(fmt.Sprintf("/%s/names", storeName), setNameHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprint("/deposit/include"), postDepositHandler(cdc, cliCtx)).Methods("POST")
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

type PostDepositReq struct {
	//BaseReq      rest.BaseReq   `json:"baseReq"`
	DepositNonce string `json:"depositNonce"`
	OwnerAddress string `json:"ownerAddress"`
}

func RegisterCodec(cdc *codec.Codec) {
}

func postDepositHandler(cdc *codec.Codec, ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostDepositReq

		//body, err := ioutil.ReadAll(r.Body)
		//if err != nil {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		//}

		//err = json.Unmarshal(body, &req)
		//fmt.Println("body", ethcmn.ToHex(body))

		//if err != nil {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		//}
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		owner := ethcmn.HexToAddress(req.OwnerAddress)
		fmt.Println("owner", req.OwnerAddress)

		//baseReq := req.BaseReq.Sanitize()
		//if !baseReq.ValidateBasic(w) {
		//	return
		//}
		fmt.Println("deposit", req.DepositNonce)

		nonce, ok := new(big.Int).SetString(req.DepositNonce, 10)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request: deposit nonce must be string base 10")
			return
		}

		msg := msgs.IncludeDepositMsg{
			DepositNonce: nonce,
			Owner:        owner,
			ReplayNonce:  uint64(0),
		}

		fmt.Println("1")
		txBytes, err := rlp.EncodeToBytes(&msg)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println("2")
		// broadcast to the node
		res, err := ctx.BroadcastTxAndAwaitCommit(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println("3")
		fmt.Println("res")
		rest.PostProcessResponse(w, cdc, res.TxHash, ctx.Indent)
		return
	}
}
