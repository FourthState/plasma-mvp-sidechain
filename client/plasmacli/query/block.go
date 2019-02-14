package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"math/big"
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
		blockNum, ok := new(big.Int).SetString(strings.TrimSpace(args[0]), 10)
		if !ok {
			return fmt.Errorf("block number must be provided in decimal format")
		}

		key := append([]byte("block::"), blockNum.Bytes()...)
		data, err := ctx.QueryStore(key, "plasma")
		if err != nil {
			fmt.Println("error querying store")
			return err
		}
		if data == nil {
			return fmt.Errorf("plasma block does not exist")
		}

		block := plasma.Block{}
		if err := rlp.DecodeBytes(data, &block); err != nil {
			fmt.Println("error decoding")
			return err
		}

		resp, err := json.Marshal(block)
		if err != nil {
			fmt.Println("error marshaling")
			return err
		}

		fmt.Printf("Block Data\n%s\n", resp)
		return nil
	},
}
