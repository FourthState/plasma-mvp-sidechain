package eth

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	tm "github.com/tendermint/tendermint/rpc/core/types"
)

func ProveCmd() *cobra.Command {
	config.AddPersistentTMFlags(proveCmd)
	return proveCmd
}

var proveCmd = &cobra.Command{
	Use:   "prove <position>",
	Short: "Prove transaction inclusion: prove <account> <position>",
	Args:  cobra.ExactArgs(2),
	Long:  "Returns proof for transaction inclusion. Use to exit transactions in the smart contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()

		// parse position
		position, err := plasma.FromPositionString(args[1])
		if err != nil {
			return err
		}

		result, sigs, err := getProof(ctx, position)
		if err != nil {
			return err
		}

		// print meta data
		fmt.Printf("Roothash: 0x%x\n", result.Proof.RootHash)
		fmt.Printf("Total: %d\n", result.Proof.Proof.Total)
		fmt.Printf("LeafHash: 0x%x\n", result.Proof.Proof.LeafHash)
		fmt.Printf("TxBytes: 0x%x\n", []byte(result.Tx))

		switch len(sigs) {
		case 65:
			fmt.Printf("Confirmation Signatures: 0x%x\n", sigs[:])
		case 130:
			fmt.Printf("Confirmation Signatures: 0x%x, 0x%x\n", sigs[:65], sigs[65:])
		}

		// flatten aunts
		var proof []byte
		for _, aunt := range result.Proof.Proof.Aunts {
			proof = append(proof, aunt...)
		}

		if len(proof) == 0 {
			if result.Proof.Proof.Total == 1 {
				fmt.Println("No proof required since this was the only transaction in the block")
			} else {
				fmt.Printf("Proof: nil\n")
			}
		} else {
			fmt.Printf("Proof: 0x%x\n", proof)
		}

		return nil
	},
}

// Returns transaction results for given position
// Trusts connected full node
func getProof(ctx context.CLIContext, position plasma.Position) (*tm.ResultTx, []byte, error) {
	key := store.GetOutputKey(position)
	hash, err := ctx.QueryStore(key, store.OutputStoreName)
	if err != nil {
		return &tm.ResultTx{}, nil, err
	}

	txKey := store.GetTxKey(hash)
	txBytes, err := ctx.QueryStore(txKey, store.OutputStoreName)

	var tx store.Transaction
	if err := rlp.DecodeBytes(txBytes, &tx); err != nil {
		return &tm.ResultTx{}, nil, fmt.Errorf("Transaction decoding failed: %s", err.Error())
	}

	// query tm node for information about this tx
	result, err := ctx.Client.Tx(tx.Transaction.MerkleHash(), true)
	if err != nil {
		return &tm.ResultTx{}, nil, err
	}

	// Look for confirmation signatures
	// Ignore error if no confirm sig currently exists in store
	var sigs []byte
	if len(tx.SpenderTxs[position.OutputIndex]) > 0 {
		queryPath := fmt.Sprintf("custom/data/tx/%s", tx.SpenderTxs[position.OutputIndex])
		data, err := ctx.Query(queryPath, nil)
		if err != nil {
			return &tm.ResultTx{}, nil, err
		}

		var spenderTx store.Transaction
		if err := json.Unmarshal(data, &spenderTx); err != nil {
			return &tm.ResultTx{}, nil, fmt.Errorf("unmarshaling json query response: %s", err)
		}
		for _, input := range spenderTx.Transaction.Inputs {
			if input.Position.String() == position.String() {
				for _, sig := range input.ConfirmSignatures {
					sigs = append(sigs, sig[:]...)
				}
			}
		}
	}

	return result, sigs, nil
}
