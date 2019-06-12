package query

import (
	"encoding/json"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/big"
)

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

func NewPlasmaQuerier(plasmaStore store.PlasmaStore) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		if len(path) == 0 {
			return nil, ErrInvalidPath("path not specified")
		}

		switch path[0] {
		case QueryBlock:
			if len(path) != 2 {
				return nil, ErrInvalidPath("expected block/<number>")
			}
			blockNum, ok := new(big.Int).SetString(path[1], 10)
			if !ok {
				return nil, ErrInvalidPath("block number must be provided in decimal format")
			} else if blockNum.Sign() < 0 {
				return nil, ErrInvalidPath("block number must be positive")
			}
			block, ok := plasmaStore.GetBlock(ctx, blockNum)
			if !ok {
				return nil, ErrInvalidPath("nonexistent plasma block")
			}
			data, err := json.Marshal(block)
			if err != nil {
				return nil, ErrSerialization("json: %s", err)
			}
			return data, nil
		case QueryBlocks:
			if len(path) > 2 {
				return nil, ErrInvalidPath("expected /blocks or /blocks/<number>")
			}

			var blockNum *big.Int
			if len(path) == 1 {
				// latest 10 blocks
				blockNum = plasmaStore.PlasmaBlockHeight(ctx)
				bigNine := big.NewInt(9)
				if blockNum.Cmp(bigNine) <= 0 {
					blockNum = big.NewInt(1)
				} else {
					blockNum = blockNum.Sub(blockNum, bigNine)
				}
			} else {
				// predefined starting point
				var ok bool
				blockNum, ok = new(big.Int).SetString(path[1], 10)
				if !ok || blockNum.Sign() < 0 {
					return nil, ErrInvalidPath("number must be in decimal format starting from 1. Got: %s", blockNum)
				}
			}

			blocks := queryBlocks(ctx, plasmaStore, blockNum)
			data, err := json.Marshal(blocks)
			if err != nil {
				return nil, ErrSerialization("json: %s", err)
			}
			return data, nil
		default:
			return nil, ErrInvalidPath("unregistered endpoint")
		}
	}
}

// queryBlocks will return an empty list of blocks if none are present
func queryBlocks(ctx sdk.Context, plasmaStore store.PlasmaStore, startPoint *big.Int) BlocksResp {
	resp := BlocksResp{startPoint, []plasma.Block{}}

	// want `startPoint` to remain the same
	blockHeight := new(big.Int).Add(startPoint, utils.Big0)
	for i := 0; i < 10; i++ {
		block, ok := plasmaStore.GetBlock(ctx, blockHeight)
		if !ok {
			return resp
		}

		resp.Blocks = append(resp.Blocks, block.Block)
		blockHeight = blockHeight.Add(blockHeight, utils.Big1)
	}

	return resp
}
