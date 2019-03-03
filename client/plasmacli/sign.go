package main

import (
	"fmt"
	clistore "github.com/FourthState/plasma-mvp-sidechain/client/store"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().StringP(flagAccount, "a", "", "Account to sign the confirmation signature with (required)")
	signCmd.Flags().String(flagOwner, "", "Owner of the output (required)")
}

var signCmd = &cobra.Command{
	Use:   "sign <position>",
	Short: "Sign confirmation signatures for position provided (blockNum.txIndex.oIndex.depositNonce), if it exists.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx := context.NewCLIContext().WithCodec(codec.New()).WithTrustNode(true)

		position, err := plasma.FromPositionString(strings.TrimSpace(args[0]))
		if err != nil {
			fmt.Println("Error parsing positions")
			return err
		}

		name := viper.GetString(flagAccount)

		ownerStr := utils.RemoveHexPrefix(strings.TrimSpace(viper.GetString(flagOwner)))
		if ownerStr == "" || !ethcmn.IsHexAddress(ownerStr) {
			return fmt.Errorf("must provide the address that owns position %s using the --owner flag in hex format", position)
		}
		ownerAddr := ethcmn.HexToAddress(ownerStr)

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
		sig, err := clistore.SignHashWithPassphrase(name, hash)
		if err != nil {
			return err
		}

		if err := clistore.SaveSig(position, sig); err != nil {
			return err
		}

		// print the results
		fmt.Printf("Confirmation Signature for output with position: %s\n", utxo.Position)
		fmt.Printf("0x%x\n", sig)

		return nil
	},
}
