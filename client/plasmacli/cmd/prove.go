package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	"strings"
)

func init() {
	rootCmd.AddCommand(proveCmd)
	proveCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
}

var proveCmd = &cobra.Command{
	Use:   "prove",
	Short: "Prove transaction inclusion: prove <address> <position>",
	Args:  cobra.ExactArgs(2),
	Long:  "Returns proof for transaction inclusion. Use to exit transactions in the smart contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()

		// validate arguments
		addrStr := strings.TrimSpace(args[0])
		if !common.IsHexAddress(addrStr) {
			return fmt.Errorf("Invalid address provided. Please use hex format")
		}
		position, err := plasma.FromPositionString(args[1])
		if err != nil {
			return err
		}

		// query for the output
		key := append(ethAddr.Bytes(), position.Bytes()...)
		res, err := ctx.QueryStore(key, ctx.UTXOStore)
		if err != nil {
			return err2
		}
		utxo := store.UTXO{}
		if err := rlp.DecodeBytes(res, &utxo); err != nil {
			return err
		}

		// query tm node for information about this tx
		result, err := ctx.Client.Tx(utxo.MerkleHash[:], true)
		if err != nil {
			return err
		}

		// Look for confirmation signatures
		cdc := amino.NewCodec()
		pos := [2]uint64{position[0].Blknum, uint64(position[0].TxIndex)}
		bz, err := cdc.MarshalBinaryBare(pos)
		if err != nil {
			return err
		}

		key = append(utils.ConfirmSigPrefix, bz...)
		res, err = ctx.QueryStore(key, ctx.PlasmaStore)

		var sigs [][65]byte
		if err == nil {
			err = ctx.Codec.UnmarshalBinaryBare(res, &sigs)
			if err != nil {
				return err
			}
		}

		// print meta data
		fmt.Printf("Roothash: 0x%x\n", result.Proof.RootHash)
		fmt.Printf("Total: %d\n", result.Proof.Proof.Total)
		fmt.Printf("LeafHash: 0x%x\n", result.Proof.Proof.LeafHash)
		fmt.Printf("TxBytes: 0x%x\n", result.Tx)

		switch len(sigs) {
		case 1:
			fmt.Printf("Confirmation Signatures: %v\n", sigs[0])
		case 2:
			fmt.Printf("Confirmation Signatures: %v,%v\n", sigs[0], sigs[1])
		}

		// flatten aunts
		var proof []byte
		for _, aunt := range result.Proof.Proof.Aunts {
			proof = append(proof, aunt...)
		}

		if len(proof) == 0 {
			fmt.Println("Proof: nil")
		} else {
			fmt.Printf("Proof: 0x%s\n", hex.EncodeToString(proof))
		}

		return nil
	},
}
