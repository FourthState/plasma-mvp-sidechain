package eth

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"strconv"
	"time"
)

func init() {
	queryCmd.AddCommand(blockCmd)
	blockCmd.Flags().String(limitF, "1", "number of plasma blocks to be displayed")
	viper.BindPFlags(blockCmd.Flags())
}

var blockCmd = &cobra.Command{
	Use:   "block <number>",
	Short: "Query a plasma block submitted to the rootchain",
	Long: `Returns the reported block header, number of transactions, fee amount,
and creation time for the requested plasma block.

Usage:
	plasmacli eth query block <number>
	plasmacli eth query block <number> --limit <number>`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		curr, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse block number: { %s }", err)
		}

		lim, err := strconv.ParseInt(viper.GetString(limitF), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse limit: { %s }", err)
		}

		end := curr + lim
		for curr < end {
			block, err := rc.contract.PlasmaChain(nil, big.NewInt(curr))
			if err != nil {
				return fmt.Errorf("failed to retrieve block: { %s }", err)
			}
			curr++
			fmt.Printf("Header: 0x%x\nTxs: %d\nFee: %d\nCreated: %v\n",
				block.Header, block.NumTxns, block.FeeAmount, time.Unix(block.CreatedAt.Int64(), 0))
		}
		return nil
	},
}
