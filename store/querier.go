package store

import (
	"encoding/json"
	"fmt"
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

	// QueryOutput retrieves a single output at
	// the given position
	QueryOutput = "output"
)

func NewOutputQuerier(outputStore OutputStore) sdk.Querier {
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
			total, err := queryBalance(ctx, outputStore, addr)
			if err != nil {
				return nil, sdk.ErrInternal("failed query balance")
			}
			return []byte(total.String()), nil

		case QueryInfo:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("expected info/<address>")
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
		case QueryOutput:
			if len(path) != 2 {
				return nil, sdk.ErrUnknownRequest("expected output/<position>")
			}
			pos, e := plasma.FromPositionString(path[1])
			if e != nil {
				return nil, sdk.ErrInternal("position decoding error")
			}
			output, err := queryOutput(ctx, outputStore, pos)
			if err != nil {
				return nil, err
			}
			data, e := json.Marshal(output)
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
		return nil, ErrAccountDNE(fmt.Sprintf("no account exists for the address provided: 0x%x", addr))
	}

	return acc.Balance, nil
}

func queryInfo(ctx sdk.Context, outputStore OutputStore, addr common.Address) ([]OutputInfo, sdk.Error) {
	acc, ok := outputStore.GetAccount(ctx, addr)
	if !ok {
		return nil, ErrAccountDNE(fmt.Sprintf("no account exists for the address provided: 0x%x", addr))
	}
	return outputStore.GetUnspentForAccount(ctx, acc), nil
}

func queryOutput(ctx sdk.Context, outputStore OutputStore, pos plasma.Position) (Output, sdk.Error) {
	output, ok := outputStore.GetOutput(ctx, pos)
	if !ok {
		return Output{}, ErrOutputDNE(fmt.Sprintf("no output exists for the position provided: %s", pos))
	}
	return output, nil
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

func NewBlockQuerier(blockStore BlockStore) sdk.Querier {
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
			block, ok := blockStore.GetBlock(ctx, blockNum)
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
				startingBlockNum = blockStore.PlasmaBlockHeight(ctx)
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

			blocks, sdkErr := queryBlocks(ctx, blockStore, startingBlockNum)
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

func queryBlocks(ctx sdk.Context, blockStore BlockStore, startPoint *big.Int) (BlocksResp, sdk.Error) {
	resp := BlocksResp{startPoint, []plasma.Block{}}

	// want `startPoint` to remain the same
	blockHeight := new(big.Int).Add(startPoint, utils.Big0)
	for i := 0; i < 10; i++ {
		block, ok := blockStore.GetBlock(ctx, blockHeight)
		if !ok {
			return resp, nil
		}

		resp.Blocks = append(resp.Blocks, block.Block)
		blockHeight = blockHeight.Add(blockHeight, utils.Big1)
	}

	return resp, nil
}
