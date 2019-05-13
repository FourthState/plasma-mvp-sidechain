package store

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/big"
)

const (
	QueryBalance = "balance"
	QueryInfo    = "info"
)

func NewUtxoQuerier(utxoStore UTXOStore) sdk.Querier {
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
			return []byte(total.String()), nil

		case QueryInfo:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("info query follows /info/<address>")
			}
			addr := common.HexToAddress(path[1])
			utxos, err := queryInfo(ctx, utxoStore, addr)
			if err != nil {
				return nil, err
			}
			data, e := json.Marshal(utxos)
			if e != nil {
				return nil, sdk.ErrInternal("serialization error")
			}
			return data, nil

		default:
			return nil, sdk.ErrUnknownRequest("unregistered endpoint")
		}
	}
}

func queryBalance(ctx sdk.Context, utxoStore UTXOStore, addr common.Address) (*big.Int, sdk.Error) {
	iter := sdk.KVStorePrefixIterator(utxoStore.KVStore(ctx), addr.Bytes())
	total := big.NewInt(0)
	for ; iter.Valid(); iter.Next() {
		utxo, ok := utxoStore.GetUTXOWithKey(ctx, iter.Key())
		if !ok {
			return nil, sdk.ErrInternal("failed utxo retrieval")
		}

		if !utxo.Spent {
			total = total.Add(total, utxo.Output.Amount)
		}
	}

	return total, nil
}

func queryInfo(ctx sdk.Context, utxoStore UTXOStore, addr common.Address) ([]UTXO, sdk.Error) {
	var utxos []UTXO
	iter := sdk.KVStorePrefixIterator(utxoStore.KVStore(ctx), addr.Bytes())
	for ; iter.Valid(); iter.Next() {
		utxo, ok := utxoStore.GetUTXOWithKey(ctx, iter.Key())
		if !ok {
			return nil, sdk.ErrInternal("failed utxo retrieval")
		}

		utxos = append(utxos, utxo)
	}

	return utxos, nil
}

const (
	QueryBlock = "block"
)

func NewPlasmaQuerier(plasmaStore PlasmaStore) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		if len(path) == 0 {
			return nil, sdk.ErrUnknownRequest("path not specified")
		}

		switch path[0] {
		case QueryBlock:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("block query follows /plasma/block/<number>")
			}
			blockNum, ok := new(big.Int).SetString(path[1], 10)
			if !ok {
				return nil, sdk.ErrUnknownRequest("block number must be provided in deicmal format")
			}
			block, ok := plasmaStore.GetBlock(ctx, blockNum)
			if !ok {
				return nil, sdk.ErrUnknownRequest("nonexistent plasma block")
			}
			data, err := json.Marshal(block)
			if err != nil {
				return nil, sdk.ErrInternal("serialization error")
			}
			return data, nil
		default:
			return nil, sdk.ErrUnknownRequest("unregistered endpoint")
		}
	}
}
