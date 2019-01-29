package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"math/big"
	"strings"
)

func init() {
	rootCmd.AddCommand(balanceCmd)
}

var balanceCmd = &cobra.Command{
	Use:   "balance <address>",
	Short: "Query Balances",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext().WithCodec(codec.New()).WithTrustNode(true)
		addrStr := utils.RemoveHexPrefix(strings.TrimSpace(args[0]))
		if !common.IsHexAddress(addrStr) {
			return fmt.Errorf("address must be provided in hex format")
		}

		addr := common.HexToAddress(addrStr)

		res, err := ctx.QuerySubspace(addr.Bytes(), "utxo")
		if err != nil {
			return err
		}

		total := big.NewInt(0)
		utxo := store.UTXO{}
		for _, pair := range res {
			if err := rlp.DecodeBytes(pair.Value, &utxo); err != nil {
				return err
			}

			if !utxo.Spent {
				fmt.Printf("Position: %s , Amount: %d\n", utxo.Position, utxo.Output.Amount.Uint64())
				total = total.Add(total, utxo.Output.Amount)
			}
		}

		fmt.Printf("Total: %d\n", total.Uint64())
		return nil
	},
}
