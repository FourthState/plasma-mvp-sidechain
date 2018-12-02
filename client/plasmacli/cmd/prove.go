package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(proveCmd)
	proveCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")

}

var proveCmd = &cobra.Command{
	Use:   "prove",
	Short: "Prove tx inclusion: prove <address> <position>",
	Long:  "Returns proof for tx inclusion. Use to exit tx on rootchain",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewClientContextFromViper()

		ethAddr := common.HexToAddress(args[0])
		position, err := client.ParsePositions(args[1])
		if err != nil {
			return err
		}

		posBytes, err := ctx.Codec.MarshalBinaryBare(position[0])
		if err != nil {
			return err
		}

		key := append(ethAddr.Bytes(), posBytes...)

		res, err2 := ctx.QueryStore(key, ctx.UTXOStore)
		if err2 != nil {
			return err2
		}

		var resUTXO utxo.UTXO
		err = ctx.Codec.UnmarshalBinaryBare(res, &resUTXO)
		if err != nil {
			return err
		}

		result, err := ctx.Client.Tx(resUTXO.TxHash, true)
		if err != nil {
			return err
		}

		fmt.Printf("Roothash: 0x%s\n", hex.EncodeToString(result.Proof.RootHash))
		fmt.Printf("Total: %d\n", result.Proof.Proof.Total)
		fmt.Printf("LeafHash: 0x%s\n", hex.EncodeToString(result.Proof.Proof.LeafHash))
		fmt.Printf("TxBytes: 0x%s\n", hex.EncodeToString(result.Tx))

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
