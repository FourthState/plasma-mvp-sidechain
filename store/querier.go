package store

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/big"
)

const (
	QueryBalance = "balance"
	QueryInfo    = "info"
)

func NewOutputQuerier(outputStore OutputStore) sdk.Querier {
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
			total, err := queryBalance(ctx, outputStore, addr)
			if err != nil {
				return nil, sdk.ErrInternal("failed query balance")
			}
			return []byte(total.String()), nil

		case QueryInfo:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("info query follows /info/<address>")
			}
			addr := common.HexToAddress(path[1])
			utxos, err := queryInfo(ctx, outputStore, addr)
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

func queryBalance(ctx sdk.Context, outputStore OutputStore, addr common.Address) (*big.Int, sdk.Error) {
	acc, ok := outputStore.GetAccount(ctx, addr)
	if !ok {
		return nil, ErrAccountDNE(DefaultCodespace, fmt.Sprintf("no account exists for the address provided: 0x%x", addr))
	}

	return acc.Balance, nil
}

func queryInfo(ctx sdk.Context, outputStore OutputStore, addr common.Address) ([]QueryOutput, sdk.Error) {
	acc, ok := outputStore.GetAccount(ctx, addr)
	if !ok {
		return nil, ErrAccountDNE(DefaultCodespace, fmt.Sprintf("no account exists for the address provided: 0x%x", addr))
	}
	return outputStore.GetUnspentForAccount(ctx, acc), nil
}

const (
	QueryBlock = "block"
)

func NewBlockQuerier(blockStore BlockStore) sdk.Querier {
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
			block, ok := blockStore.GetBlock(ctx, blockNum)
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
