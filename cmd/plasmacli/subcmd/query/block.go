package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/store"
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
			if !ok {
				return fmt.Errorf("block number must be in decimal format")
			}
		}

		blocks, err := BlocksMetadata(ctx, blockNum)
		if err != nil {
			return err
		}

		if len(blocks) == 0 {
			fmt.Println("no blocks")
			return nil
		}

		for _, block := range blocks {
			fmt.Printf("Block Height: %d\n", block.Height)
			fmt.Printf("Header: 0x%x\n", block.Header)
			fmt.Printf("Transaction Count: %d\n", block.TxnCount)
			fmt.Printf("Fee Amount: %s\n\n", block.FeeAmount)
		}

		return nil
	},
}

func Block(ctx context.CLIContext, num *big.Int) (store.Block, error) {
	if num == nil || num.Sign() <= 0 {
		return store.Block{}, fmt.Errorf("block number starts at 1")
	}

	queryPath := fmt.Sprintf("custom/data/block/%s", num)
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

// BlocksMetadata will query metadata about 10 plasma blocks starting from `startingBlockNum`.
// The latest 10 blocks will be retrieved if startingBlockNum is nil
func BlocksMetadata(ctx context.CLIContext, startingBlockNum *big.Int) ([]store.Block, error) {
	var queryPath string
	if startingBlockNum == nil {
		queryPath = "custom/data/blocks/latest"
	} else if startingBlockNum.Sign() <= 0 {
		return nil, fmt.Errorf("block number starts at 1")
	} else {
		queryPath = fmt.Sprintf("custom/data/blocks/%s", startingBlockNum)
	}

	data, err := ctx.Query(queryPath, nil)
	if err != nil {
		return nil, err
	}

	var resp []store.Block
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("json unmarshal: %s", err)
	}

	return resp, nil
}
