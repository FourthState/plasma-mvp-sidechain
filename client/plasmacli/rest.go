package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/query"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/server/app"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/types"
	"io/ioutil"
	"net/http"
)

func init() {
	cdc := app.MakeCodec()
	rootCmd.AddCommand(lcd.ServeCommand(cdc, registerRoutes))
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerRoutes(rs *lcd.RestServer) {
	ctx := rs.CliCtx.WithTrustNode(true)
	r := rs.Mux

	r.HandleFunc("/balance/{address}", balanceHandler(ctx)).Methods("GET")
	r.HandleFunc("/block/{num:[0-9]+}", blockHandler(ctx)).Methods("GET")
	r.HandleFunc("/blocks", blocksHandler(ctx)).Methods("GET")
	r.HandleFunc("/blocks/{num:[0-9]+}", blocksHandler(ctx)).Methods("GET")
	r.HandleFunc("/info/{address}", infoHandler(ctx)).Methods("GET")

	r.HandleFunc("/submit", submitHandler(ctx)).Methods("POST")
}

func balanceHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		addr := vars["address"]
		if !common.IsHexAddress(addr) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("address must be a 20-byte hex string"))
			return
		}

		total, err := query.Balance(ctx, common.HexToAddress(addr))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			// Log the error
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
		if !common.IsHexAddress(addr) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("address must be a 20-byte hex string"))
			return
		}

		utxos, err := query.Info(ctx, common.HexToAddress(addr))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// Log the error
			return
		}

		data, err := json.Marshal(utxos)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			// log the error
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func submitHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txBytes, err := ioutil.ReadAll(hex.NewDecoder(r.Body))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("unable to read transaction bytes. Body must be in hex format"))
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
		_, err = ctx.BroadcastTxAndAwaitCommit(txBytes)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// retrieve full information about a block
func blockHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		num := mux.Vars(r)["num"]

		if num == "0" || num == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("plasma blocks start at 1"))
			return
		}

		block, err := query.Block(ctx, num)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			// TODO: check the against the codespace type
			// maybe the block does not exist?
			return
		}

		type resp struct {
			plasma.Block
			Txs []types.Tx
		}

		// Query the tendermint block
		height := int64(block.TMBlockHeight)
		tmBlock, err := ctx.Client.Block(&height)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		data, err := json.Marshal(resp{block.Block, tmBlock.Block.Data.Txs})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// retrieve metadata about the last 10 blocks
func blocksHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		num, ok := mux.Vars(r)["num"]
		if !ok {
			num = ""
		}

		if num == "0" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("plasma blocks start at 1"))
			return
		}

		blocksResp, err := query.BlocksMetadata(ctx, num)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			// TODO: check the against the codespace type
			// maybe the block does not exist?
			return
		}

		data, err := json.Marshal(blocksResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
