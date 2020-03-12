package query

import (
	"encoding/hex"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	"strings"
)

// TxCommand -
func TxCommand() *cobra.Command {
	return txCmd
}

var txCmd *cobra.Command = &cobra.Command{
	Use:   "tx <txHash/position>",
	Short: "Query for information about a single transaction or transaction output",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()
		cmd.SilenceUsage = true
		arg := strings.TrimSpace(args[0])

		// argument validation
		if pos, err := plasma.FromPositionString(arg); err != nil {
			// utxo position
			txOutput, err := client.TxOutput(ctx, pos)
			if err != nil {
				fmt.Printf("Error quering transaction output %s: %s\n", pos, err)
				return err
			}

			// display output information
			fmt.Printf("TxOutput: %s\n", pos)
			fmt.Printf("Confirmation Hash: 0x%x\n", txOutput.ConfirmationHash)
			fmt.Printf("Transaction Hash: 0x%x\n", txOutput.TxHash)
			fmt.Printf("Spent: %v\n", txOutput.Spent)
			if txOutput.Spent {
				fmt.Printf("Spending Tx Hash: 0x%x\n", txOutput.SpenderTx)
			}

		} else {
			// transaction hash
			txHash, err := hex.DecodeString(utils.RemoveHexPrefix(arg))
			if err != nil {
				fmt.Printf("Error decoding tx hash hex-encoded string: %s\n", err)
				return err
			}

			tx, err := client.Tx(ctx, txHash)
			if err != nil {
				fmt.Printf("Error queyring trasnsaction hash: 0x%x\n", txHash)
				return err
			}

			// display transaction information
			fmt.Printf("Block %d, Tx Index: %d\n", tx.Position.BlockNum, tx.Position.TxIndex)
			fmt.Printf("TxHash: 0x%x\n", txHash)
			fmt.Printf("Confirmation Hash: 0x%x\n", tx.ConfirmationHash)
			// we expect the length of `SpenderTx` and `Spent` to be the same
			for i, spender := range tx.SpenderTxs {
				spent := tx.Spent[i]
				fmt.Printf("Ouput %d:\n\tSpent: %v\n", i, spent)
				if spent {
					// keep the indentendation
					fmt.Printf("\tSpender Tx: 0x%x\n", spender)
				}
			}
		}

		return nil
	},
}
