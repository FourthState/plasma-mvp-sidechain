package cmd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"strconv"
	"strings"
)

const (
	flagFrom  = "from"  // address to sign with
	flagOwner = "owner" // address that owns the utxo being queried for
)

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().String(flagFrom, "", "Address used to sign the confirmation signature (required)")
	signCmd.Flags().String(flagOwner, "", "Owner of the newly created utxo (required)")
	viper.BindPFlags(signCmd.Flags())
}

var signCmd = &cobra.Command{
	Use:   "sign <position>",
	Short: "Sign confirmation signatures for position provided (0.0.0.0), if it exists.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCLIContext()

		position, err := plasma.FromPositionString(args[0])
		if err != nil {
			return err
		}

		// parse flags
		fromStr := viper.GetString(flagFrom)
		if fromStr == "" {
			return errors.New("must provide the address to sign with using the --from flag")
		}
		ownerStr := viper.GetString(flagOwner)
		if ownerStr == "" {
			return fmt.Errorf("must provide the address that owns position %s using the --owner flag", position)
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
		if err := rlp.DecodeBytes(res, &input); err != nil {
			return err
		}

		// create the signature
		hash := utils.ToEthSignedMessageHash(utxo.ConfirmationHash)
		dir := viper.GetString(client.FlagHomeDir)
		keystore.InitKeyStore(dir)
		acct, err := keystore.Find(signerAddr)
		if err != nil {
			return err
		}
		sig, err := keystore.SignHashWithPassphrase(acct, hash)
		if err != nil {
			return err
		}

		// print the results
		fmt.Printf("Confirmation Signature for utxo with position: %s\n", input.Position, input.Amount.String())
		fmt.Printf("Signature: %x\n", sig)

		return nil
	},
}
