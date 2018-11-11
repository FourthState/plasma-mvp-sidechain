package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/FourthState/plasma-mvp-sidechain/types"
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

		var utxo types.BaseUTXO
		err = ctx.Codec.UnmarshalBinaryBare(res, &utxo)
		if err != nil {
			return err
		}

		result, err := ctx.Client.Tx(utxo.TxHash, true)
		if err != nil {
			return err
		}

		formattedProof, err := json.MarshalIndent(result.Proof.Proof, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(formattedProof))

		return nil
	},
}
