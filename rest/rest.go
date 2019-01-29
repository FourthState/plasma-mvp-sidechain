package rest

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/types"
	"math/big"
	"net/http"
)

// RegisterRoutes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/plasma/block/{blockNum:[0-9]+}", blockHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/plasma/block/{blockNum:[0-9]+}/txs", txHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/plasma/block/", latestBlockHandler(cliCtx)).Methods("GET")
}

// retrieve the latest block
func latestBlockHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// retrieve the latest block
		key := []byte("plasmaBlockNum")
		data, err := cliCtx.QueryStore(key, "plasma")
		if err != nil || data == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		key = append([]byte("block::"), data...)
		blockData, err := cliCtx.QueryStore(key, "plasma")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		plasmaBlock := plasma.Block{}
		if err := rlp.DecodeBytes(blockData, &plasmaBlock); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// serialize the plasma block
		type blk struct {
			Header    string `json:"header"`
			TxnCount  uint16 `json:"txnCount"`
			FeeAmount string `json:"feeAmount"`
			BlockNum  string `json:"blockNum"`
		}

		resp, err := json.Marshal(blk{
			Header:    fmt.Sprintf("%x", plasmaBlock.Header),
			TxnCount:  plasmaBlock.TxnCount,
			FeeAmount: plasmaBlock.FeeAmount.String(),
			BlockNum:  new(big.Int).SetBytes(data).String(),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("json marshal error"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func blockHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		blockNum, ok := new(big.Int).SetString(vars["blockNum"], 10)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid block number"))
			return
		}

		key := append([]byte("block::"), blockNum.Bytes()...)
		plasmaBlockData, err := cliCtx.QueryStore(key, "plasma")
		if err != nil || plasmaBlockData == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("non existent block"))
			return
		}

		plasmaBlock := plasma.Block{}
		if err := rlp.DecodeBytes(plasmaBlockData, &plasmaBlock); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("corrupt store"))
			return
		}

		// serialize the plasma block
		type blk struct {
			Header    string `json:"header"`
			TxnCount  uint16 `json:"TxnCount"`
			FeeAmount string `json:"FeeAmount"`
		}

		resp, err := json.Marshal(blk{
			Header:    fmt.Sprintf("%x", plasmaBlock.Header),
			TxnCount:  plasmaBlock.TxnCount,
			FeeAmount: plasmaBlock.FeeAmount.String(),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("json marshal error"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func txHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		blockNum, ok := new(big.Int).SetString(vars["blockNum"], 10)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid block number"))
			return
		}

		// find the tm block number
		key := append([]byte("plasmatotm::"), blockNum.Bytes()...)
		tmBlockNumData, err := cliCtx.QueryStore(key, "plasma")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("non existent block"))
			return
		}
		tmBlockNum := new(big.Int).SetBytes(tmBlockNumData).Int64()

		// query the tmp block
		tmBlock, err := cliCtx.Client.Block(&tmBlockNum)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("non existent tm block"))
			return
		}

		type txStruct struct {
			Txs types.Txs `json:"txs"`
		}

		resp, err := json.Marshal(txStruct{
			Txs: tmBlock.Block.Data.Txs,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("json marshal error"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}
