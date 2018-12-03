package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")

}

var balanceCmd = &cobra.Command{
	Use:   "balance <address>",
	Short: "Query Balances",
	Long:  "Query Balances",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewClientContextFromViper()

		ethAddr := common.HexToAddress(args[0])

		res, err2 := ctx.QuerySubspace(ethAddr.Bytes(), ctx.UTXOStore)
		if err2 != nil {
			return err2
		}

		for _, pair := range res {
			var resUTXO utxo.UTXO
			err := ctx.Codec.UnmarshalBinaryBare(pair.Value, &resUTXO)
			if err != nil {
				return err
			}
			if resUTXO.Valid {
				fmt.Printf("Position: %v \nAmount: %d \n", resUTXO.Position, resUTXO.Amount)
			}
		}

		return nil
	},
}
