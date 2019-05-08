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
	QueryBalance = "balance"
	QueryInfo    = "info"
)

type BalanceResp struct {
	Address common.Address
	Total   *big.Int
}

type InfoResp struct {
	Address common.Address
	Utxos   []store.UTXO
}

func NewUtxoQuerier(utxoStore store.UTXOStore) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		if len(path) == 0 {
			return nil, sdk.ErrUnknownRequest("path not specified")
		}

		switch path[0] {
		case QueryBalance:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("balance query follows balance/<address>")
			}
			addr := common.HexToAddress(path[1])
			total, err := queryBalance(ctx, utxoStore, addr)
			if err != nil {
				return nil, sdk.ErrInternal("failed query balance")
			}
			data, err := json.Marshal(BalanceResp{addr, total})
			if err != nil {
				return nil, sdk.ErrInternal("serialization error")
			}
			return data, nil

		case QueryInfo:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("info query follows /info/<address>")
			}
			addr := common.HexToAddress(path[1])
			utxos, err := queryInfo(ctx, utxoStore, addr)
			if err != nil {
				return nil, sdk.ErrInternal("failed utxo retrieval")
			}
			data, err := json.Marshal(InfoResp{addr, utxos})
			if err != nil {
				return nil, sdk.ErrInternal("serialization error")
			}
			return data, nil

		default:
			return nil, sdk.ErrUnknownRequest("unregistered endpoint")
		}
	}
}

const (
	QueryBlock = "block"
)

func NewPlasmaQuerier(plasmaStore store.PlasmaStore) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		if len(path) == 0 {
			return nil, sdk.ErrUnknownRequest("path not specified")
		}

		switch path[0] {
		default:
			return nil, sdk.ErrUnknownRequest("unregistered endpoint")
		}
	}
}
