package query

import (
	"encoding/json"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/big"
)

const (
	// QueryBalance retrieves the aggregate value of
	// the set of owned by the specified address
	QueryBalance = "balance"

	// QueryInfo retrieves the entire utxo set owned
	// by the specified address
	QueryInfo = "info"
)

func NewUtxoQuerier(utxoStore store.UTXOStore) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		if len(path) == 0 {
			return nil, ErrInvalidPath("path not specified")
		}

		switch path[0] {
		case QueryBalance:
			if len(path) != 2 {
				return nil, ErrInvalidPath("expected balance/<address>")
			}
			addr := common.HexToAddress(path[1])
			utxos := utxoStore.GetUTXOSet(ctx, addr)

			total := big.NewInt(0)
			for _, utxo := range utxos {
				if !utxo.Spent {
					total = total.Add(total, utxo.Output.Amount)
				}
			}
			return []byte(total.String()), nil

		case QueryInfo:
			if len(path) != 2 {
				return nil, ErrInvalidPath("expected info/<address>")
			}
			addr := common.HexToAddress(path[1])
			utxos := utxoStore.GetUTXOSet(ctx, addr)
			data, err := json.Marshal(utxos)
			if err != nil {
				return nil, ErrSerialization("json: %s", err)
			}
			return data, nil

		default:
			return nil, ErrInvalidPath("unregistered endpoint")
		}
	}
}
