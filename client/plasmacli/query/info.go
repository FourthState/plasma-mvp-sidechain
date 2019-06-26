package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	queryCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info <address>",
	Args:  cobra.ExactArgs(1),
	Short: "Information on owned utxos valid and invalid",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()
		addr := args[0]
		if !common.IsHexAddress(addr) {
			return fmt.Errorf("Invalid address provided. Please use hex format")
		}
		addr := common.HexToAddress(addrStr)
		fmt.Printf("Querying information for 0x%x\n\n", addr)

		// valid arguments
		cmd.SilenceUsage = true

		queryPath := fmt.Sprintf("custom/tx/info/%s", addr)
		data, err := ctx.Query(queryPath, nil)
		if err != nil {
			return err
		}

		var utxos []store.OuptputInfo
		if err := json.Unmarshal(data, &utxos); err != nil {
			return fmt.Errorf("unmarshaling json query response: %s", err)
		}

		for i, utxo := range utxos {
			fmt.Printf("UTXO %d\n", i)
			fmt.Printf("Position: %s, Amount: %s, Spent: %t\nSpender Hash: %s\n", utxo.Tx.Position, utxo.Output.Output.Amount.String(), utxo.Output.Spent, utxo.Output.SpenderTx)
			fmt.Printf("Transaction Hash: 0x%x\nConfirmationHash: 0x%x\n", utxo.Tx.Transaction.TxHash(), utxo.Tx.ConfirmationHash)
			// print inputs if applicable
			positions := utxo.Tx.Transaction.InputPositions()
			for i, p := range positions {
				fmt.Printf("Input %d Position: %s\n", i, p)
			}

			fmt.Printf("End UTXO %d info\n\n", i)
		}

		if len(utxos) == 0 {
			fmt.Println("no information available for this address")
		}

		return nil
	},
}

func Info(ctx context.CLIContext, addr common.Address) ([]store.OutputInfo, error) {
	// query for all utxos owned by this address
	queryRoute := fmt.Sprintf("custom/utxo/info/%s", addr.Hex())
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return nil, err
	}

	var utxos []store.OutputInfo
	if err := json.Unmarshal(data, &utxos); err != nil {
		return nil, err
	}

	return utxos, nil
}
