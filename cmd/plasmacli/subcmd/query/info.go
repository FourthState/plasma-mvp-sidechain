package query

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/flags"
	ks "github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// InfoCmd returns the query information command
func InfoCmd() *cobra.Command {
	return infoCmd
}

var infoCmd = &cobra.Command{
	Use:   "info <account/address>",
	Short: "Information on owned utxos valid and invalid",
	Long: `Information on owned utxos valid and invalid.
If --verbose is set, the input positions that created each output will also be displayed`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()
		cmd.SilenceUsage = true
		var (
			addr ethcmn.Address
			err  error
		)

		if !ethcmn.IsHexAddress(args[0]) {
			if addr, err = ks.GetAccount(args[0]); err != nil {
				return fmt.Errorf("failed local account retrieval: %s", err)
			}
		} else {
			addr = ethcmn.HexToAddress(args[0])
		}

		utxos, err := client.Info(ctx, addr)
		if err != nil {
			return err
		}

		for i, utxo := range utxos {
			fmt.Printf("UTXO %d\n", i)
			fmt.Printf("Position: %s, Amount: %s, Spent: %t\nSpender Hash: %s\n", utxo.Position, utxo.Output.Amount.String(), utxo.Spent, utxo.SpenderTx)
			fmt.Printf("Transaction Hash: 0x%x\nConfirmationHash: 0x%x\n", utxo.TxHash, utxo.ConfirmationHash)

			// query for the inputs that created this out
			if viper.GetBool(flags.Verbose) {
				txInput, err := client.TxInput(ctx, utxo.Position)
				if err != nil {
					fmt.Printf("Error retrieving further information about the inputs of this UTXO: %s", err)
				} else {
					for i, input := range txInput.InputPositions {
						fmt.Printf("Input %d: Position: %s\n", i, input)
					}
				}
			}

			fmt.Printf("End UTXO %d info\n\n", i)
		}

		if len(utxos) == 0 {
			fmt.Println("No information available for this address")
		}

		return nil
	},
}
