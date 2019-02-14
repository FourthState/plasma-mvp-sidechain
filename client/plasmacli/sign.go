package main

import (
	"errors"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

const (
	flagSinger = "signer"
	flagOwner  = "owner"
)

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().String(flagSinger, "", "Address used to sign the confirmation signature (required)")
	signCmd.Flags().String(flagOwner, "", "Owner of the output (required)")
	viper.BindPFlags(signCmd.Flags())
}

var signCmd = &cobra.Command{
	Use:   "sign <position>",
	Short: "Sign confirmation signatures for position provided (blockNum.txIndex.oIndex.depositNonce), if it exists.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext().WithCodec(codec.New()).WithTrustNode(true)

		position, err := plasma.FromPositionString(strings.TrimSpace(args[0]))
		if err != nil {
			fmt.Println("Error parsing positions")
			return err
		}

		// parse flags
		fromStr := utils.RemoveHexPrefix(strings.TrimSpace(viper.GetString(flagSinger)))
		if fromStr == "" || !common.IsHexAddress(fromStr) {
			return errors.New("must provide the address to sign with using the --from flag in hex format")
		}
		ownerStr := utils.RemoveHexPrefix(strings.TrimSpace(viper.GetString(flagOwner)))
		if ownerStr == "" || !common.IsHexAddress(ownerStr) {
			return fmt.Errorf("must provide the address that owns position %s using the --owner flag in hex format", position)
		}
		ownerAddr := common.HexToAddress(ownerStr)
		signerAddr := common.HexToAddress(fromStr)

		// retrieve the new output
		utxo := store.UTXO{}
		key := append(ownerAddr.Bytes(), position.Bytes()...)
		res, err := ctx.QueryStore(key, "utxo")
		if err != nil {
			return err
		}
		if err := rlp.DecodeBytes(res, &utxo); err != nil {
			return err
		}

		// create the signature
		hash := utils.ToEthSignedMessageHash(utxo.ConfirmationHash)
		acct, err := keystore.Find(signerAddr)
		if err != nil {
			return err
		}
		sig, err := keystore.SignHashWithPassphrase(acct.Address, hash)
		if err != nil {
			return err
		}

		// print the results
		fmt.Printf("Confirmation Signature for output with position: %s\n", utxo.Position)
		fmt.Printf("0x%x\n", sig)

		return nil
	},
}
