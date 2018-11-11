package cmd

import (
	"encoding/binary"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/FourthState/plasma-mvp-sidechain/app"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func init() {
	rootCmd.AddCommand(signCmd)
	viper.BindPFlags(signCmd.Flags())
}

var signCmd = &cobra.Command{
	Use:   "sign <position> <address>",
	Short: "Sign confirmation signatures for position provided (0.0.0.0), if it exists",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.NewClientContextFromViper()
		cdc := app.MakeCodec()
		ctx = ctx.WithCodec(cdc)

		ethAddr := common.HexToAddress(args[1])
		position, err := client.ParsePositions(args[0])
		if err != nil {
			return err
		}

		posBytes, err := ctx.Codec.MarshalBinaryBare(position[0])
		if err != nil {
			return err
		}

		key := append(ethAddr.Bytes(), posBytes...)
		res, err := ctx.QueryStore(key, ctx.UTXOStore)
		if err != nil {
			return err
		}

		var utxo types.BaseUTXO
		err = ctx.Codec.UnmarshalBinaryBare(res, &utxo)

		blknumKey := make([]byte, binary.MaxVarintLen64)
		binary.PutUvarint(blknumKey, utxo.GetPosition().Get()[0].Uint64())

		blockhash, err := ctx.QueryStore(blknumKey, ctx.MetadataStore)
		if err != nil {
			return err
		}

		hash := ethcrypto.Keccak256(append(utxo.MsgHash, blockhash...))

		dir := viper.GetString(FlagHomeDir)
		ks := client.GetKeyStore(dir)

		acc := accounts.Account{
			Address: ethAddr,
		}
		// get account to sign with
		acct, err := ks.Find(acc)
		if err != nil {
			return err
		}

		// get passphrase
		passphrase, err := ctx.GetPassphraseFromStdin(ethAddr)
		if err != nil {
			return err
		}

		sig, err := ks.SignHashWithPassphrase(acct, passphrase, hash)

		fmt.Printf("\nConfirmation Signature for utxo with\nposition: %v \namount: %d\n", utxo.Position, utxo.Amount)
		fmt.Printf("signature:%x\n", sig)

		inputLen := 1
		// check number of inputs
		if !utils.ZeroAddress(utxo.InputAddresses[1]) {
			inputLen = 2
		}
		fmt.Printf("UTXO had %d inputs\n", inputLen)
		return nil
	},
}
