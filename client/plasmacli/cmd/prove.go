package cmd

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/ethereum/go-ethereum/common"
	rlp "github.com/ethereum/go-ethereum/rlp"
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

		fmt.Printf("Roothash: 0x%s\n", hex.EncodeToString(result.Proof.RootHash))
		fmt.Printf("Total: %d\n", result.Proof.Proof.Total)
		fmt.Printf("LeafHash: 0x%s\n", hex.EncodeToString(result.Proof.Proof.LeafHash))
		fmt.Printf("TxBytes: 0x%s\n", hex.EncodeToString(result.Tx))

		var tx types.BaseTx
		rlp.DecodeBytes(result.Tx, &tx)

		fmt.Printf("Transaction: %+v\n", tx)

		var proof []byte
		for _, aunt := range result.Proof.Proof.Aunts {
			proof = append(proof, aunt...)
		}

		fmt.Printf("Proof: 0x%s\n", hex.EncodeToString(proof))

		hasher := sha256.New()
		var buf [10]byte
		n := binary.PutUvarint(buf[:], 32)
		_, err = hasher.Write(buf[0:n])
		fmt.Printf("buffer: %s\n", hex.EncodeToString(buf[0:n]))
		hasher.Write(proof)
		hasher.Write(buf[0:n])
		hasher.Write(result.Proof.Proof.LeafHash)
		value := hasher.Sum(nil)

		fmt.Println()
		fmt.Printf("Verify: %s\n", hex.EncodeToString(value))

		return nil
	},
}
