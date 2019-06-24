package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/store/query"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	"math/big"
)

func BlockCmd() *cobra.Command {
	return blockCmd
}

func BlocksCmd() *cobra.Command {
	return blocksCmd
}

var blockCmd = &cobra.Command{
	Use:   "block <block number>",
	Short: "Query information about a plasma block",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()
		num, ok := new(big.Int).SetString(args[0], 10)
		if !ok {
			return fmt.Errorf("number must be in decimal format")
		}
		cmd.SilenceUsage = true

		block, err := Block(ctx, num)
		if err != nil {
			return err
		}

		fmt.Printf("Block Header: 0x%x\n", block.Header)
		fmt.Printf("Transaction Count: %d, FeeAmount: %d\n", block.TxnCount, block.FeeAmount)
		fmt.Printf("Tendermint BlockHeight: %d\n", block.TMBlockHeight)

		return nil
	},
}

var blocksCmd = &cobra.Command{
	Use:   "blocks <number>",
	Short: "Query Metadata about blocks",
	Long:  "Query Metadata about blocks. If no height is provided, the latest 10 blocks will be queried. Otherwise 10 starting from the specified height",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()

		var blockNum *big.Int
		if len(args) > 0 {
			var ok bool
			blockNum, ok = new(big.Int).SetString(args[0], 10)
			if !ok || blockNum.Sign() <= 0 {
				return fmt.Errorf("number must be in decimal format starting from 1. Got: %s", args[0])
			}
		}

		// done with argument checking
		cmd.SilenceUsage = true

		blocksResp, err := BlocksMetadata(ctx, blockNum)
		if err != nil {
			return err
		}

		if len(blocksResp.Blocks) == 0 {
			fmt.Printf("no block starting at height %s\n", blocksResp.StartingBlockHeight)
			return nil
		}

		blockHeight := blocksResp.StartingBlockHeight
		for _, block := range blocksResp.Blocks {
			fmt.Printf("Block Height: %s\n", blockHeight)
			fmt.Printf("Header: 0x%x\n", block.Header)
			fmt.Printf("Transaction Count: %d\n", block.TxnCount)
			fmt.Printf("Fee Amount: %s\n\n", block.FeeAmount)

			blockHeight = blockHeight.Add(blockHeight, utils.Big1)
		}

		return nil
	},
}

func Block(ctx context.CLIContext, num *big.Int) (store.Block, error) {
	if num == nil || num.Sign() <= 0 {
		return store.Block{}, fmt.Errorf("block number starts at 1")
	}

	queryPath := fmt.Sprintf("custom/plasma/block/%s", num)
	data, err := ctx.Query(queryPath, nil)
	if err != nil {
		return store.Block{}, err
	}

	var block store.Block
	if err := json.Unmarshal(data, &block); err != nil {
		return store.Block{}, fmt.Errorf("json unmarshal: %s", err)
	}

	return block, nil
}

// BlocksMetadata will query metadate about 10 plasma blocks starting from `startingBlockNum`.
// The latest 10 blocks will be retrieved if startingBlockNum is nil
func BlocksMetadata(ctx context.CLIContext, startingBlockNum *big.Int) (query.BlocksResp, error) {
	var queryPath string
	if startingBlockNum == nil {
		queryPath = "custom/plasma/blocks"
	} else if startingBlockNum.Sign() <= 0 {
		return query.BlocksResp{}, fmt.Errorf("block number starts at 1")
	} else {
		queryPath = fmt.Sprintf("custom/plasma/blocks/%s", startingBlockNum)
	}

	data, err := ctx.Query(queryPath, nil)
	if err != nil {
		return query.BlocksResp{}, err
	}

	var resp query.BlocksResp
	if err := json.Unmarshal(data, &resp); err != nil {
		return query.BlocksResp{}, fmt.Errorf("json unmarshal: %s", err)
	}

	return resp, nil
}
