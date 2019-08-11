package client

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"
	"math/big"
	"net/http"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	// Getters
	r.HandleFunc("/balance/{address}", balanceHandler(ctx)).Methods("GET")
	r.HandleFunc("/info/{address}", infoHandler(ctx)).Methods("GET")
	r.HandleFunc("/block/{height}", blockHandler(ctx)).Methods("GET")
	r.HandleFunc("/blocks/{height}", blocksHandler(ctx)).Methods("GET")
	r.HandleFunc("/tx/{hash}", txHandler(ctx)).Methods("GET")
	r.HandleFunc("/output/{position}", outputHandler(ctx)).Methods("GET")

	// Post
	r.HandleFunc("/submit", submitHandler(ctx)).Methods("POST")
}

func balanceHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["address"]
		if !ethcmn.IsHexAddress(addr) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("address must be an ethereum 20-byte hex string"))
			return
		}

		total, err := Balance(ctx, ethcmn.HexToAddress(addr))
		if err != nil {
			writeClientRetrievalErr(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(total))
	}
}

func infoHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["address"]
		if !ethcmn.IsHexAddress(addr) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("address must be an ethereum 20-byte hex string"))
			return
		}

		txo, err := Info(ctx, ethcmn.HexToAddress(addr))
		if err != nil {
			writeClientRetrievalErr(w, err)
			return
		}

		writeJSONResponse(txo, w)
	}
}

func blockHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		num, ok := new(big.Int).SetString(mux.Vars(r)["height"], 10)
		if !ok || num.Sign() <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("block height must be in decimal format starting from 1"))
			return
		}

		block, err := Block(ctx, num)
		if err != nil {
			writeClientRetrievalErr(w, err)
			return
		}

		resp := struct {
			store.Block
			Txs [][]byte
		}{
			Block: block,
			Txs:   [][]byte{},
		}

		// Query the tendermint block.
		// Tendermint stores transactions in base64 format
		//
		// Transcode base64 to hex
		height := int64(block.TMBlockHeight)
		tmBlock, err := ctx.Client.Block(&height)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		for _, tx := range tmBlock.Block.Data.Txs {
			hexFormat, err := base64.StdEncoding.DecodeString(string(tx))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error with transaction transcoding"))
				return
			}

			resp.Txs = append(resp.Txs, hexFormat)
		}

		writeJSONResponse(resp, w)
	}
}

func blocksHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var blockHeight *big.Int
		arg := mux.Vars(r)["height"]
		if arg != "latest" {
			var ok bool
			if blockHeight, ok = new(big.Int).SetString(arg, 10); !ok || blockHeight.Sign() <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("block height must be in decimal format starting from 1. /blocks/latest for the lastest 10 blocks"))
			}
		}

		blocks, err := Blocks(ctx, blockHeight)
		if err != nil {
			writeClientRetrievalErr(w, err)
			return
		}

		writeJSONResponse(blocks, w)
	}
}

func txHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txHash := utils.RemoveHexPrefix(mux.Vars(r)["hash"])

		// validation
		_, err := hex.DecodeString(txHash)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("tx hash expected in hexadecimal format")))
			return
		} else if len(txHash) != 64 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("tx hash expected to be 32 bytes in length"))
			return
		}

		tx, err := Tx(ctx, txHash)
		if err != nil {
			writeClientRetrievalErr(w, err)
			return
		}

		writeJSONResponse(tx, w)
	}
}

func outputHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//pos := mux.Vars(r)["position"]
	}
}

func submitHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type reqBody struct {
			Async   bool   `json:"async"`
			TxBytes string `json:"txBytes"`
		}

		var body reqBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("unable to read request body: " + err.Error()))
			return
		}

		// clean up txBytes string
		if len(body.TxBytes) > 2 && body.TxBytes[:2] == "0x" || body.TxBytes[:2] == "0X" {
			body.TxBytes = body.TxBytes[2:]
		}
		if len(body.TxBytes)%2 != 0 {
			body.TxBytes = "0" + body.TxBytes
		}
		txBytes, err := hex.DecodeString(body.TxBytes)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("tx bytes must be in hexadecimal format"))
			return
		}

		var tx plasma.Transaction
		if err := rlp.DecodeBytes(txBytes, &tx); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("malformed tx bytes"))
			return
		}

		if err := tx.ValidateBasic(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// deliver the tx
		if body.Async {
			_, err = ctx.BroadcastTxAsync(txBytes)
		} else {
			_, err = ctx.BroadcastTxAndAwaitCommit(txBytes)
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

/****  Helpers ****/

func writeJSONResponse(obj interface{}, w http.ResponseWriter) {
	data, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("json:" + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func writeClientRetrievalErr(w http.ResponseWriter, err error) {
	// If the client request (GET) could not be fulfilled for some
	// other reason than the requested information not existing (DNE), the request
	// must have been malformed or an internal server error must have occured
	sdkerr, ok := err.(sdk.Error)
	if !ok || sdkerr.Code() != store.CodeDNE {
		w.WriteHeader(http.StatusBadRequest)
		// TODO: log the error
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

	w.Write([]byte(err.Error()))
}
