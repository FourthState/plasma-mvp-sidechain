package cmd

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	flagFrom  = "from"  // address to sign with
	flagOwner = "owner" // address that owns the utxo being queried for
)

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().String(flagFrom, "", "Address used to sign the confirmation signature")
	signCmd.Flags().String(flagOwner, "", "Owner of the newly created utxo")
	viper.BindPFlags(signCmd.Flags())
}

var signCmd = &cobra.Command{
	Use:   "sign <position>",
	Short: "Sign confirmation signatures for position provided (0.0.0.0), if it exists.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.NewClientContextFromViper()

		position, err := client.ParsePositions(args[0])
		if err != nil {
			return err
		}

		fromStr := viper.GetString(flagFrom)
		if fromStr == "" {
			return errors.New("must provide the address to sign with using the --from flag")
		}

		ownerStr := viper.GetString(flagOwner)
		if ownerStr == "" {
			return fmt.Errorf("must provide the address that owns position %v using the --owner flag", position)
		}

		ethAddr := common.HexToAddress(ownerStr)
		signerAddr := common.HexToAddress(fromStr)

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

		hash := tmhash.Sum(append(utxo.TxHash, blockhash...))
		signHash := utils.SignHash(hash)

		dir := viper.GetString(FlagHomeDir)
		ks := client.GetKeyStore(dir)

		acc := accounts.Account{
			Address: signerAddr,
		}
		// get account to sign with
		acct, err := ks.Find(acc)
		if err != nil {
			return err
		}

		// get passphrase
		passphrase, err := ctx.GetPassphraseFromStdin(signerAddr)
		if err != nil {
			return err
		}

		sig, err := ks.SignHashWithPassphrase(acct, passphrase, signHash)

		fmt.Printf("\nConfirmation Signature for utxo with\nposition: %v \namount: %d\n", utxo.Position, utxo.Amount)
		fmt.Printf("signature: %x\n", sig)

		inputLen := 1
		// check number of inputs
		if !utils.ZeroAddress(utxo.InputAddresses[1]) {
			inputLen = 2
		}
		fmt.Printf("UTXO had %d inputs\n", inputLen)
		return nil
	},
}
