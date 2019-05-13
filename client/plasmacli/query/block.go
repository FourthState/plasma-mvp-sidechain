package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	queryCmd.AddCommand(blockCmd)
}

var blockCmd = &cobra.Command{
	Use:   "block <block number>",
	Short: "Query information about a plasma block",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext().WithTrustNode(true)
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
