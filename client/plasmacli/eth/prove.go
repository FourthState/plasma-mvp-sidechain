package eth

import (
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	tm "github.com/tendermint/tendermint/rpc/core/types"
)

func ProveCmd() *cobra.Command {
	return proveCmd
}

var proveCmd = &cobra.Command{
	Use:   "prove",
	Short: "Prove transaction inclusion: prove <name> <position>",
	Args:  cobra.ExactArgs(2),
	Long:  "Returns proof for transaction inclusion. Use to exit transactions in the smart contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, err := ks.Get(args[0])
		if err != nil {
			return err
		}

		// parse position
		position, err := plasma.FromPositionString(args[1])
		if err != nil {
			return err
		}

		result, confirmSignatures, err := getProof(addr, position)
		if err != nil {
			return err
		}
		// print meta data
		fmt.Printf("Roothash: 0x%x\n", result.Proof.RootHash)
		fmt.Printf("Total: %d\n", result.Proof.Proof.Total)
		fmt.Printf("LeafHash: 0x%x\n", result.Proof.Proof.LeafHash)
		fmt.Printf("TxBytes: 0x%x\n", result.Tx)

		switch len(confirmSignatures) {
		case 1:
			fmt.Printf("Confirmation Signatures: 0x%x\n", confirmSignatures[0])
		case 2:
			fmt.Printf("Confirmation Signatures: 0x%x, 0x%x\n", confirmSignatures[0], confirmSignatures[1])
		}

		// flatten aunts
		var proof []byte
		for _, aunt := range result.Proof.Proof.Aunts {
			proof = append(proof, aunt...)
		}

		if len(proof) == 0 {
			fmt.Printf("Proof: nil\n")
		} else {
			fmt.Printf("Proof: 0x%x\n", proof)
		}

		return nil
	},
}

// Returns transaction results for given position
// Trusts connected full node
func getProof(addr ethcmn.Address, position plasma.Position) (*tm.ResultTx, [][65]byte, error) {
	ctx := context.NewCLIContext().WithTrustNode(true)

	// query for the output
	key := append(addr.Bytes(), position.Bytes()...)
	res, err := ctx.QueryStore(key, "utxo")
	if err != nil {
		return &tm.ResultTx{}, nil, err
	}

	utxo := store.UTXO{}
	if err := rlp.DecodeBytes(res, &utxo); err != nil {
		return &tm.ResultTx{}, nil, err
	}

	// query tm node for information about this tx
	result, err := ctx.Client.Tx(utxo.MerkleHash, true)
	if err != nil {
		return &tm.ResultTx{}, nil, err
	}

	// Look for confirmation signatures
	var confirmSignatures [][65]byte
	key = append([]byte("confirmSignature"), utxo.Position.Bytes()...)
	res, err = ctx.QueryStore(key, "plasma")

	if err == nil { // confirm signatures exist
		var signature [65]byte
		copy(signature[:], res)
		confirmSignatures = append(confirmSignatures, signature)

		if len(res) > 65 {
			copy(signature[:], res[65:])
			confirmSignatures = append(confirmSignatures, signature)
		}
	}

	return result, confirmSignatures, nil
}

// Returns all information necessary to start an exit
// return (position, txbytes, proof, confirmation signatures, error)
// Trust connected full node
// TODO: Add nodeURL flag
func proveExit(addr ethcmn.Address, position plasma.Position) ([]byte, []byte, []byte, error) {
	result, confirmSignatures, err := getProof(addr, position)
	if err != nil {
		return nil, nil, nil, err
	}
	var sigs []byte
	for i, sig := range confirmSignatures {
		sigs = append(sigs, sig[i])
	}

	// flatten aunts
	var proof []byte
	for _, aunt := range result.Proof.Proof.Aunts {
		proof = append(proof, aunt...)
	}

	return result.Tx, proof, sigs, nil
}
