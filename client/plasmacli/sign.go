package main

import (
	"fmt"
	clistore "github.com/FourthState/plasma-mvp-sidechain/client/store"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	cosmoscli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

const (
	signPrompt = "Would you like to finalize this transaction? [Y/n]"
)

func init() {
	signCmd.Flags().String(ownerF, "", "Owner of the output (required with position flag)")
	signCmd.Flags().String(positionF, "", "Position of transaction to finalize (required with owner flag)")
}

var signCmd = &cobra.Command{
	Use:   "sign <account>",
	Short: "Sign confirmation signatures for pending transactions",
	Long: `Iterate over all unfinalized transaction corresponding to the provided account. 
Prompt the user for confirmation to finailze the pending transactions. Owner and Position flags can be used to finalize a specific transaction.

Usage:
	plasmacli sign <account>
	plasmacli sign <account> --owner <address> --position "(blknum.txindex.oindex.depositnonce)"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx := context.NewCLIContext().WithCodec(codec.New()).WithTrustNode(true)

		name := args[0]

		signerAddr, err := clistore.GetAccount(name)
		if err != nil {
			return err
		}

		ownerS := viper.GetString(ownerF)
		positionS := viper.GetString(positionF)
		if ownerS != "" || positionS != "" {
			position, err := plasma.FromPositionString(strings.TrimSpace(viper.GetString(positionF)))
			if err != nil {
				return err
			}

			ownerStr := utils.RemoveHexPrefix(strings.TrimSpace(viper.GetString(ownerF)))
			if ownerStr == "" || !ethcmn.IsHexAddress(ownerStr) {
				return fmt.Errorf("must provide the address that owns position %s using the --owner flag in hex format", position)
			}
			ownerAddr := ethcmn.HexToAddress(ownerStr)

			err = signSingleConfirmSig(ctx, position, signerAddr, ownerAddr, name)
			return err
		}

		res, err := ctx.QuerySubspace(signerAddr.Bytes(), "utxo")
		if err != nil {
			return err
		}

		utxo := store.UTXO{}
		for _, pair := range res {
			if err := rlp.DecodeBytes(pair.Value, &utxo); err != nil {
				return err
			}

			if utxo.Spent {
				spenderPositions := utxo.SpenderPositions()
				spenderAddresses := utxo.SpenderAddresses()
				for i, pos := range spenderPositions {
					err = signSingleConfirmSig(ctx, pos, signerAddr, spenderAddresses[i], name)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}

		return nil
	},
}

// generate confirmation signature for specified owner and position
func signSingleConfirmSig(ctx context.CLIContext, position plasma.Position, signerAddr, owner ethcmn.Address, name string) error {
	// retrieve the new output
	utxo := store.UTXO{}
	key := append(owner.Bytes(), position.Bytes()...)
	res, err := ctx.QueryStore(key, "utxo")
	if err != nil {
		return err
	}

	if err := rlp.DecodeBytes(res, &utxo); err != nil {
		return err
	}

	if err := verifyAndSign(utxo, signerAddr, name); err != nil {
		return err
	}
	return nil
}

// verify that the inputs provided are correct
// signing address should match one of the input addresses
// generate confirmation signature for given utxo
func verifyAndSign(utxo store.UTXO, signerAddr ethcmn.Address, name string) error {
	sig, _ := clistore.GetSig(utxo.Position)
	inputAddrs := utxo.InputAddresses()

	if len(sig) == 130 || (len(sig) == 65 && len(inputAddrs) == 1) {
		return nil
	}

	for i, input := range inputAddrs {
		if input != signerAddr {
			continue
		}

		// get confirmation to generate signature
		fmt.Printf("\nUTXO\nPosition: %s\nOwner: 0x%x\nValue: %d\n", utxo.Position, utxo.Output.Owner, utxo.Output.Amount)
		buf := cosmoscli.BufferStdin()
		auth, err := cosmoscli.GetString(signPrompt, buf)
		if err != nil {
			return err
		}
		if auth != "Y" {
			return nil
		}

		hash := utils.ToEthSignedMessageHash(utxo.ConfirmationHash)
		sig, err := clistore.SignHashWithPassphrase(name, hash)
		if err != nil {
			return fmt.Errorf("failed to generate confirmation signature: { %s }", err)
		}

		if err := clistore.SaveSig(utxo.Position, sig, i == 0); err != nil {
			return err
		}

		// print the results
		fmt.Printf("Confirmation Signature for output with position: %s\n", utxo.Position)
		fmt.Printf("0x%x\n", sig)
	}
	return nil
}
