package query

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
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

		block, err := client.Block(ctx, num)
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

		var startingHeight *big.Int
		if len(args) > 0 {
			var ok bool
			startingHeight, ok = new(big.Int).SetString(args[0], 10)
			if !ok {
				return fmt.Errorf("block number must be in decimal format")
			}
		}

		blocks, err := client.Blocks(ctx, startingHeight)
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
