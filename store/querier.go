package store

import (
	"encoding/json"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
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

func NewUtxoQuerier(utxoStore UTXOStore) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		if len(path) == 0 {
			return nil, sdk.ErrUnknownRequest("path not specified")
		}

		switch path[0] {
		case QueryBalance:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("exprected balance/<address>")
			}
			addr := common.HexToAddress(path[1])
			total, err := queryBalance(ctx, utxoStore, addr)
			if err != nil {
				return nil, sdk.ErrInternal("failed query balance")
			}
			return []byte(total.String()), nil

		case QueryInfo:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("expected info/<address>")
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
	// QueryBlocks retrieves full information about a
	// speficied block
	QueryBlock = "block"

	// QueryBlocs retrieves metadata about 10 blocks from
	// a specified start point or the last 10 from the latest
	// block
	QueryBlocks = "blocks"
)

type BlocksResp struct {
	StartingBlockHeight *big.Int
	Blocks              []plasma.Block
}

func NewPlasmaQuerier(plasmaStore PlasmaStore) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		if len(path) == 0 {
			return nil, sdk.ErrUnknownRequest("path not specified")
		}

		switch path[0] {
		case QueryBlock:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("expected block/<number>")
			}
			blockNum, ok := new(big.Int).SetString(path[1], 10)
			if !ok {
				return nil, sdk.ErrUnknownRequest("block number must be provided in decimal format")
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
		case QueryBlocks:
			if len(path) > 2 {
				return nil, sdk.ErrUnknownRequest("expected /blocks or /blocks/<starting block num>")
			}

			var startingBlockNum *big.Int
			if len(path) == 1 {
				// latest 10 blocks
				startingBlockNum = plasmaStore.PlasmaBlockHeight(ctx)
				bigNine := big.NewInt(9)
				if startingBlockNum.Cmp(bigNine) <= 0 {
					startingBlockNum = big.NewInt(1)
				} else {
					startingBlockNum = startingBlockNum.Sub(startingBlockNum, bigNine)
				}
			} else {
				// predefined starting point
				var ok bool
				startingBlockNum, ok = new(big.Int).SetString(path[1], 10)
				if !ok {
					return nil, sdk.ErrUnknownRequest("block number must be in decimal format")
				}
			}

			blocks, sdkErr := queryBlocks(ctx, plasmaStore, startingBlockNum)
			if sdkErr != nil {
				return nil, sdkErr
			}
			data, err := json.Marshal(blocks)
			if err != nil {
				return nil, sdk.ErrInternal("serialization error")
			}
			return data, nil
		default:
			return nil, sdk.ErrUnknownRequest("unregistered endpoint")
		}
	}
}

func queryBlocks(ctx sdk.Context, plasmaStore PlasmaStore, startPoint *big.Int) (BlocksResp, sdk.Error) {
	resp := BlocksResp{startPoint, []plasma.Block{}}

	// want `startPoint` to remain the same
	blockHeight := new(big.Int).Add(startPoint, utils.Big0)
	for i := 0; i < 10; i++ {
		block, ok := plasmaStore.GetBlock(ctx, blockHeight)
		if !ok {
			return resp, nil
		}

		resp.Blocks = append(resp.Blocks, block.Block)
		blockHeight = blockHeight.Add(blockHeight, utils.Big1)
	}

	return resp, nil
}
