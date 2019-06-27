package subcmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/app"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/subcmd/query"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/types"
	"math/big"
	"net/http"
)

func RestServerCmd() *cobra.Command {
	serverCmd.Flags().String(client.FlagListenAddr, "tcp://localhost:1317", "The address for the server to listen on")
	serverCmd.Flags().Bool(client.FlagTLS, false, "Enable SSL/TLS layer")
	serverCmd.Flags().String(client.FlagSSLHosts, "", "Comma-separated hostnames and IPs to generate a certificate for")
	serverCmd.Flags().String(client.FlagSSLCertFile, "", "Path to a SSL certificate file. If not supplied, a self-signed certificate will be generated.")
	serverCmd.Flags().String(client.FlagSSLKeyFile, "", "Path to a key file; ignored if a certificate file is not supplied.")
	serverCmd.Flags().String(client.FlagCORS, "", "Set the domains that can make CORS requests (* for all)")
	serverCmd.Flags().Int(client.FlagMaxOpenConnections, 1000, "The number of maximum open connections")
	return serverCmd
}

var serverCmd = &cobra.Command{
	Use:   "rest-server",
	Short: "Start LCD (light-client daemon), a local REST server",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		rs := lcd.NewRestServer(app.MakeCodec())

		registerRoutes(rs)

		// Start the rest server and return error if one exists
		err := rs.Start(
			viper.GetString(client.FlagListenAddr),
			viper.GetString(client.FlagSSLHosts),
			viper.GetString(client.FlagSSLCertFile),
			viper.GetString(client.FlagSSLKeyFile),
			viper.GetInt(client.FlagMaxOpenConnections),
			viper.GetBool(client.FlagTLS))

		return err
	},
}

func registerRoutes(rs *lcd.RestServer) {
	ctx := rs.CliCtx
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

		queryPath := fmt.Sprintf("custom/utxo/balance/%s", addr)
		total, err := ctx.Query(queryPath, nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(total)
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

		queryPath := fmt.Sprintf("custom/utxo/info/%s", addr)
		data, err := ctx.Query(queryPath, nil)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
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

// retrieve full information about a block
func blockHandler(ctx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		num, _ := new(big.Int).SetString(mux.Vars(r)["num"], 10)

		block, err := query.Block(ctx, num)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
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

		var blockNum *big.Int
		if num != "" {
			var ok bool
			blockNum, ok = new(big.Int).SetString(num, 10)
			if !ok || blockNum.Sign() <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("number must be in decimal format starting from 1"))
				return
			}
		}

		blocksResp, err := query.BlocksMetadata(ctx, blockNum)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
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
