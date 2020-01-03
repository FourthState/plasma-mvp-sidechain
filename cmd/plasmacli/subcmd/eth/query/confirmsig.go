package query

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
)

// SigCmd returns the query confirm sig command
func SigCmd() *cobra.Command {
	config.AddPersistentTMFlags(sigCmd)
	sigCmd.Flags().Bool(useNodeF, false, "trust connected full node")
	return sigCmd
}

var sigCmd = &cobra.Command{
	Use:   "sig <position>",
	Short: "Query confirm signature information for a given position",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// parse position
		pos, err := plasma.FromPositionString(args[0])
		if err != nil {
			return err
		}

		var sigs []byte
		ctx := context.NewCLIContext()
		sigs, err = getSigs(ctx, pos)
		if err != nil {
			return fmt.Errorf("failed to retrieve confirm signature information: %s", err)
		}

		switch len(sigs) {
		case 65:
			fmt.Printf("Confirmation Signatures: 0x%x\n", sigs[:])
		case 130:
			fmt.Printf("Confirmation Signatures: 0x%x, 0x%x\n", sigs[:65], sigs[65:])
		}

		return nil
	},
}

// Returns confirm sig results for given position
// Trusts connected full node
func getSigs(ctx context.CLIContext, position plasma.Position) ([]byte, error) {
	key := store.GetOutputKey(position)
	hash, err := ctx.QueryStore(key, store.DataStoreName)
	if err != nil {
		return nil, err
	}

	txKey := store.GetTxKey(hash)
	txBytes, err := ctx.QueryStore(txKey, store.DataStoreName)

	var tx store.Transaction
	if err := rlp.DecodeBytes(txBytes, &tx); err != nil {
		return nil, fmt.Errorf("transaction decoding failed: %s", err.Error())
	}

	// Look for confirmation signatures
	// Ignore error if no confirm sig currently exists in store
	var sigs []byte
	if len(tx.SpenderTxs[position.OutputIndex]) > 0 {
		queryPath := fmt.Sprintf("custom/data/tx/%s", tx.SpenderTxs[position.OutputIndex])
		data, err := ctx.Query(queryPath, nil)
		if err != nil {
			return nil, err
		}

		var spenderTx store.Transaction
		if err := json.Unmarshal(data, &spenderTx); err != nil {
			return nil, fmt.Errorf("unmarshaling json query response: %s", err)
		}
		for _, input := range spenderTx.Transaction.Inputs {
			if input.Position.String() == position.String() {
				for _, sig := range input.ConfirmSignatures {
					sigs = append(sigs, sig[:]...)
				}
			}
		}
	}

	return sigs, nil
}
