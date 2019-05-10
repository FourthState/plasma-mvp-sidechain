package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/query"
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

		queryPath := fmt.Sprintf("custom/plasma/block/%s", num)
		data, err := ctx.Query(queryPath, nil)
		if err != nil {
			return err
		}

		var resp query.BlockResp
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		}

		fmt.Printf("Block Header: 0x%x\n", resp.Header)
		fmt.Printf("Transaction Count: %d, FeeAmount: %d\n", resp.TxnCount, resp.FeeAmount)
		fmt.Printf("Tendermint BlockHeight: %d\n", resp.TMBlockHeight)
		return nil
	},
}
