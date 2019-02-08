package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")

}

var infoCmd = &cobra.Command{
	Use:   "info <address>",
	Args:  cobra.ExactArgs(1),
	Short: "Information on owned utxos valid and invalid",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext().WithCodec(codec.New())
		addrStr := strings.TrimSpace(args[0])
		if !common.IsHexAddress(addrStr) {
			return fmt.Errorf("Invalid address provided. Please use hex format")
		}

		// query for all utxos owned by this address
		res, err := ctx.QuerySubspace(common.HexToAddress(addrStr).Bytes(), "utxo")
		if err != nil {
			return err
		}

		utxo := store.UTXO{}
		for i, pair := range res {
			if err := rlp.DecodeBytes(pair.Value, &utxo); err != nil {
				return err
			}

			fmt.Printf("UTXO %d\n", i)
			fmt.Printf("Position: %s, Amount: %s, Spent: %t\n", utxo.Position, utxo.Output.Amount.String(), utxo.Spent)

			// print inputs if applicable
			inputAddresses := utxo.InputAddresses()
			positions := utxo.InputPositions()
			for i, _ := range inputAddresses {
				fmt.Printf("Input Owner %d, Position: %s\n", i, positions[i])
			}

			// print spenders if applicable
			spenderAddresses := utxo.SpenderAddresses()
			positions = utxo.SpenderPositions()
			for i, _ := range spenderAddresses {
				fmt.Printf("Spender Owner %d, Position: %s\n", i, positions[i])
			}

			fmt.Printf("End UTXO %d info\n\n", i)
		}

		fmt.Println()

		return nil
	},
}
