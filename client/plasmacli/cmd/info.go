package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")

}

var infoCmd = &cobra.Command{
	Use:   "info <address>",
	Short: "Information on owned utxos valid and invalid",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewClientContextFromViper()

		ethAddr := common.HexToAddress(args[0])

		res, err2 := ctx.QuerySubspace(ethAddr.Bytes(), ctx.UTXOStore)
		if err2 != nil {
			return err2
		}

		for i, pair := range res {
			fmt.Printf("Start UTXO %d info:\n", i)
			var resUTXO utxo.UTXO
			err := ctx.Codec.UnmarshalBinaryBare(pair.Value, &resUTXO)
			if err != nil {
				return err
			}
			fmt.Printf("\nPosition: %v \nAmount: %d \nDenomination: %s \nValid: %t\n", resUTXO.Position, resUTXO.Amount, resUTXO.Denom, resUTXO.Valid)

			if resUTXO.InputKeys != nil {
				inputOwners := resUTXO.InputAddresses()
				inputs := resUTXO.InputPositions(ctx.Codec, types.ProtoPosition)
				for i, key := range resUTXO.InputKeys {
					plasmaInput, _ := inputs[i].(*types.PlasmaPosition)
					fmt.Printf("\nInput Owner %d: %s\nInput Position %d: %v\nInputKey %d in UTXO store: %s\n", i, hex.EncodeToString(inputOwners[i]),
						i, *plasmaInput, i, hex.EncodeToString(key))
				}
			}

			if resUTXO.SpenderKeys != nil {
				spenders := resUTXO.SpenderAddresses()
				spenderPositions := resUTXO.SpenderPositions(ctx.Codec, types.ProtoPosition)
				for i, key := range resUTXO.SpenderKeys {
					plasmaPosition, _ := spenderPositions[i].(*types.PlasmaPosition)
					fmt.Printf("\nSpender %d: %s\nSpending Position %d: %v\nSpendKey %d in UTXO store: %s\n", i, hex.EncodeToString(spenders[i]),
						i, *plasmaPosition, i, hex.EncodeToString(key))
				}
			}

			fmt.Printf("End UTXO %d info:\n", i)
		}

		fmt.Println()

		return nil
	},
}
