package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	"math/big"
	"strings"
)

func init() {
	queryCmd.AddCommand(blockCmd)
	queryCmd.AddCommand(blocksCmd)
}

var blockCmd = &cobra.Command{
	Use:   "block <block number>",
	Short: "Query information about a plasma block",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()
		num := strings.TrimSpace(args[0])

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
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()
		if len(args) > 1 {
			return fmt.Errorf("only maximum one argument")
		}

		startingBlockNum := ""
		if len(args) == 1 {
			startingBlockNum = strings.TrimSpace(args[0])
			if startingBlockNum == "0" {
				return fmt.Errorf("plasma blocks start at height 1")
			}

			_, ok := new(big.Int).SetString(startingBlockNum, 10)
			if !ok {
				return fmt.Errorf("provided block height must be in decimal format")
			}

			fmt.Printf("Querying blocks from height %s...\n", startingBlockNum)
		} else {
			fmt.Printf("No height specified. Querying the latest 10 blocks..\n")
		}

		// finished argument checking
		cmd.SilenceUsage = true

		blocksResp, err := BlocksMetadata(ctx, startingBlockNum)
		if err != nil {
			return err
		}

		if len(blocksResp.Blocks) == 0 {
			fmt.Printf("no block starting at height %s\n", startingBlockNum)
			return nil
		}

		blockHeight := blocksResp.StartingBlockHeight
		for _, block := range blocksResp.Blocks {
			fmt.Println("")
			fmt.Printf("Block Height: %s\n", blockHeight)
			fmt.Printf("Header: 0x%x\n", block.Header)
			fmt.Printf("Transaction Count: %d\n", block.TxnCount)
			fmt.Printf("Fee Amount: %s\n", block.FeeAmount)

			blockHeight = blockHeight.Add(blockHeight, utils.Big1)
		}

		return nil
	},
}

func Block(ctx context.CLIContext, num string) (store.Block, error) {
	queryPath := fmt.Sprintf("custom/plasma/block/%s", num)
	data, err := ctx.Query(queryPath, nil)
	if err != nil {
		return store.Block{}, err
	}

	var block store.Block
	if err := json.Unmarshal(data, &block); err != nil {
		return store.Block{}, err
	}

	return block, nil
}

func BlocksMetadata(ctx context.CLIContext, startingBlockNum string) (store.BlocksResp, error) {
	var queryPath string
	if startingBlockNum == "" {
		queryPath = "custom/plasma/blocks"
	} else {
		queryPath = fmt.Sprintf("custom/plasma/blocks/%s", startingBlockNum)
	}

	data, err := ctx.Query(queryPath, nil)
	if err != nil {
		return store.BlocksResp{}, err
	}

	var resp store.BlocksResp
	if err := json.Unmarshal(data, &resp); err != nil {
		return store.BlocksResp{}, err
	}

	return resp, nil
}
