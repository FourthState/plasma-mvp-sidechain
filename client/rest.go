package client

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"
	"math/big"
	"net/http"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/balance/{address}", balanceHandler(ctx)).Methods("GET")
	r.HandleFunc("/info/{address}", infoHandler(ctx)).Methods("GET")
	r.HandleFunc("/block/{height}", blockHandler(ctx)).Methods("GET")
	r.HandleFunc("/blocks/{height}", blocksHandler(ctx)).Methods("GET")
	r.HandleFunc("/tx/{hash}", txHandler(ctx)).Methods("GET")
	r.HandleFunc("/output/{position}", outputHandler(ctx)).Methods("GET")
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
			sdkerr, ok := err.(sdk.Error)
			if !ok || sdkerr.Code() != store.CodeDNE {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}

			w.Write([]byte(err.Error()))
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
			sdkerr, ok := err.(sdk.Error)
			if !ok || sdkerr.Code() != store.CodeDNE {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}

			w.Write([]byte(err.Error()))
			return
		}

		data, err := json.Marshal(txo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
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
			sdkerr, ok := err.(sdk.Error)
			if !ok || sdkerr.Code() != store.CodeDNE {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}

			w.Write([]byte(err.Error()))
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
		// Transcode base64 encoded into hex format
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

		data, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
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
			sdkerr, ok := err.(sdk.Error)
			if !ok || sdkerr.Code() != store.CodeDNE {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}

			w.Write([]byte(err.Error()))
			return
		}

		data, err := json.Marshal(blocks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func txHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func outputHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
