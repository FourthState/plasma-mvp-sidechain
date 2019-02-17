package main

import (
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(proveCmd)
}

var proveCmd = &cobra.Command{
	Use:   "prove",
	Short: "Prove transaction inclusion: prove <name> <position>",
	Args:  cobra.ExactArgs(2),
	Long:  "Returns proof for transaction inclusion. Use to exit transactions in the smart contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext().WithTrustNode(true)
		name := args[0]

		addr, err := ks.Get(name)
		if err != nil {
			return err
		}

		// parse position
		position, err := plasma.FromPositionString(args[1])
		if err != nil {
			return err
		}

		// query for the output
		key := append(addr.Bytes(), position.Bytes()...)
		res, err := ctx.QueryStore(key, "utxo")
		if err != nil {
			return err
		}
		utxo := store.UTXO{}
		if err := rlp.DecodeBytes(res, &utxo); err != nil {
			return err
		}

		// query tm node for information about this tx
		result, err := ctx.Client.Tx(utxo.MerkleHash, true)
		if err != nil {
			return err
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
