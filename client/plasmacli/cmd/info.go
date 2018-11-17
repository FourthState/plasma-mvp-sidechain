package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")

}

var infoCmd = &cobra.Command{
	Use:   "info <address>",
	Short: "Information on owned utxos ",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewClientContextFromViper()

		ethAddr := common.HexToAddress(args[0])

		res, err2 := ctx.QuerySubspace(ethAddr.Bytes(), ctx.UTXOStore)
		if err2 != nil {
			return err2
		}

		for _, pair := range res {
			var utxo types.BaseUTXO
			err := ctx.Codec.UnmarshalBinaryBare(pair.Value, &utxo)
			if err != nil {
				return err
			}
			fmt.Printf("\nPosition: %v \nAmount: %d \n", utxo.Position, utxo.Amount)
			inputAddrHelper(utxo)
		}

		return nil
	},
}

func inputAddrHelper(utxo types.BaseUTXO) {
	if utxo.Position.DepositNum == 0 {
		fmt.Printf("First Input Address: %s \n", utxo.InputAddresses[0].Hex())
		if !utils.ZeroAddress(utxo.InputAddresses[1]) {
			fmt.Printf("Second Input Address: %s \n", utxo.InputAddresses[1].Hex())
		}
	}
}
