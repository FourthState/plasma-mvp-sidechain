package query

import (
	"encoding/json"
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func InfoCmd() *cobra.Command {
	return infoCmd
}

var infoCmd = &cobra.Command{
	Use:          "info <account/address>",
	Short:        "Information on owned utxos valid and invalid",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()
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

		queryPath := fmt.Sprintf("custom/utxo/info/%s", addr)
		data, err := ctx.Query(queryPath, nil)
		if err != nil {
			return err
		}

		var utxos []store.TxOutput
		if err := json.Unmarshal(data, &utxos); err != nil {
			return fmt.Errorf("unmarshaling json query response: %s", err)
		}

		for i, utxo := range utxos {
			fmt.Printf("UTXO %d\n", i)
			fmt.Printf("Position: %s, Amount: %s, Spent: %t\nSpender Hash: %s\n", utxo.Position, utxo.Output.Amount.String(), utxo.Spent, utxo.SpenderTx)
			fmt.Printf("Transaction Hash: 0x%x\nConfirmationHash: 0x%x\n", utxo.TxHash, utxo.ConfirmationHash)
			// print inputs if applicable
			positions := utxo.InputPositions
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

func Info(ctx context.CLIContext, addr ethcmn.Address) ([]store.TxOutput, error) {
	// query for all utxos owned by this address
	queryRoute := fmt.Sprintf("custom/utxo/info/%s", addr.Hex())
	data, err := ctx.Query(queryRoute, nil)
	if err != nil {
		return nil, err
	}

	var utxos []store.TxOutput
	if err := json.Unmarshal(data, &utxos); err != nil {
		return nil, err
	}

	return utxos, nil
}
