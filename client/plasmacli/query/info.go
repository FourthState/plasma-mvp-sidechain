package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	queryCmd.AddCommand(infoCmd)
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
		addr := common.HexToAddress(addrStr)
		fmt.Printf("Querying information for 0x%x\n", addr)

		// query for all utxos owned by this address
		queryRoute := fmt.Sprintf("custom/utxo/info/%s", addr.Hex())
		data, err := ctx.Query(queryRoute, nil)
		if err != nil {
			return err
		}

		var utxos []store.UTXO
		if err := json.Unmarshal(data, &utxos); err != nil {
			return err
		}

		for i, utxo := range utxos {
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

		if len(utxos) == 0 {
			fmt.Println("no information available for this address")
		}

		return nil
	},
}
