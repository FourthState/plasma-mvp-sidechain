package query

import (
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/store"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"math/big"
)

func init() {
	queryCmd.AddCommand(balanceCmd)
}

// TODO: Change to querying account, add flag for getting all info
var balanceCmd = &cobra.Command{
	Use:   "balance <name>",
	Short: "Query account balance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext().WithCodec(codec.New()).WithTrustNode(true)
		name := args[0]

		addr, err := ks.GetAccount(name)
		if err != nil {
			return err
		}

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
