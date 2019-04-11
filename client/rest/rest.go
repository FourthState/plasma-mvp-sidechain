package rest

import (
	"fmt"
	"net/http"
	url "net/url"

	//clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	//sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	//"github.com/cosmos/sdk-application-tutorial/x/nameservice"
	ethcli "github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/eth"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	//"github.com/FourthState/plasma-mvp-sidechain/utils"
	ethcmn "github.com/ethereum/go-ethereum/common"
	hex "github.com/ethereum/go-ethereum/common/hexutil"
	//"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"
	tm "github.com/tendermint/tendermint/rpc/core/types"

	plasma "github.com/FourthState/plasma-mvp-sidechain/plasma"

	//"io/ioutil"
	"math/big"
)

const (
	ownerAddress = "ownerAddress"
	position     = "position"
	signature    = "signature"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// r.HandleFunc(fmt.Sprintf("/%s/names", storeName), buyNameHandler(cdc, cliCtx)).Methods("POST")
	// r.HandleFunc(fmt.Sprintf("/%s/names", storeName), setNameHandler(cdc, cliCtx)).Methods("PUT")
	// r.HandleFunc(fmt.Sprintf("/%s/names", storeName), setNameHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprint("/deposit/include"), postDepositHandler(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/balance/{%s}", ownerAddress), queryBalanceHandler(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprint("/utxo"), queryUTXOHandler(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprint("/proof"), queryProofHandler(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprint("/health"), healthHandler(cdc, cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprint("/spend"), postSpendHandler(cdc, cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprint("/tx/hash"), postTxHashHandler(cdc, cliCtx)).Methods("POST")
}

func healthHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rest.PostProcessResponse(w, cdc, "healthy", cliCtx.Indent)
	}
}

type ProofResp struct {
	ResultTx   tm.ResultTx `json:"transaction"`
	ProofAunts string      `json:"proofAunts"`
}

func queryProofHandler(cdc *codec.Codec, ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vals, err := url.ParseQuery(r.URL.RawQuery)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ownerStr := vals[ownerAddress][0]
		posStr := vals[position][0]

		owner := ethcmn.HexToAddress(ownerStr)
		pos, err := plasma.FromPositionString(posStr)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var proofAunts []byte

		tx, _, err := ethcli.GetProof(owner, pos)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// flatten proof
		for _, aunt := range tx.Proof.Proof.Aunts {
			proofAunts = append(proofAunts, aunt...)
		}

		proofAuntsHex := hex.Encode(proofAunts)
		proofResp := ProofResp{*tx, proofAuntsHex}

		rest.PostProcessResponse(w, cdc, proofResp, ctx.Indent)
	}

}

func queryUTXOHandler(cdc *codec.Codec, ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vals, err := url.ParseQuery(r.URL.RawQuery)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ownerStr := vals[ownerAddress][0]
		posStr := vals[position][0]

		owner := ethcmn.HexToAddress(ownerStr)
		pos, err := plasma.FromPositionString(posStr)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		key := store.GetUTXOStoreKey(owner, pos)
		res, err := ctx.QuerySubspace(key, "utxo")

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if len(res) != 1 {
			rest.WriteErrorResponse(w, http.StatusNotFound, "No UTXO found")
			return
		}

		utxoBytes := res[0]

		utxo := store.UTXO{}
		if err := rlp.DecodeBytes(utxoBytes.Value, &utxo); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		rest.PostProcessResponse(w, cdc, utxo, ctx.Indent)

	}

}

func queryBalanceHandler(cdc *codec.Codec, ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[ownerAddress]

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
		//fmt.Println("owner", req.OwnerAddress)

		//baseReq := req.BaseReq.Sanitize()
		//if !baseReq.ValidateBasic(w) {
		//	return
		//}
		//fmt.Println("deposit", req.DepositNonce)

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

		//fmt.Println("1")
		txBytes, err := rlp.EncodeToBytes(&msg)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		//fmt.Println("2")
		// broadcast to the node
		res, err := ctx.BroadcastTxAndAwaitCommit(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		//fmt.Println("3")
		//fmt.Println("res")
		rest.PostProcessResponse(w, cdc, res.TxHash, ctx.Indent)
		return
	}
}

func postSpendHandler(cdc *codec.Codec, ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req plasma.Transaction

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Failed to parse request")
			return
		}

		//fmt.Println("tx to spend: ", req)

		//txHash := req.TxHash()

		//txHashMessage := utils.ToEthSignedMessageHash(txHash)
		//fmt.Println("TX HASH MESSAGE ", ethcmn.ToHex(txHashMessage))

		//pubKey, err := crypto.SigToPub(txHashMessage, req.Input0.Signature[:])
		//if err != nil {
		//	rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		//	return
		//}
		//signAddr := crypto.PubkeyToAddress(*pubKey)
		//fmt.Printf("address from signature: %x \n", signAddr)

		// create SpendMsg and txBytes
		msg := msgs.SpendMsg{
			Transaction: req,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid `SpendMsg`: "+err.Error())
			return
		}

		txBytes, err := rlp.EncodeToBytes(&msg)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Failed to encode `SpendMsg`: "+err.Error())
			return
		}

		// broadcast to the node
		res, err := ctx.BroadcastTxAndAwaitCommit(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Failed to broadcast `SpendMsg` "+err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res.TxHash, ctx.Indent)
		return
	}
}

func postTxHashHandler(cdc *codec.Codec, ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req plasma.Transaction

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Failed to parse request")
			return
		}

		txHash := ethcmn.ToHex(req.TxHash())
		//fmt.Println("txHash (hex): ", txHash)

		rest.PostProcessResponse(w, cdc, txHash, ctx.Indent)
	}
}
