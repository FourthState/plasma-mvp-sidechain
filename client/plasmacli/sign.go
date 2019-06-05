package main

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/query"
	clistore "github.com/FourthState/plasma-mvp-sidechain/client/store"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	cosmoscli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	ethcmn "github.com/ethereum/go-ethereum/common"
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

		utxos, err := query.Info(ctx, signerAddr)
		if err != nil {
			return err
		}

		for _, output := range utxos {

			if output.Output.Spent {
				spenderPositions := output.Tx.SpenderPositions()
				spenderAddresses := output.Tx.SpenderAddresses()
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
	// query for output for the specified position
	output, err := query.Output(ctx, position)
	if err != nil {
		return err
	}

	if err := verifyAndSign(output, signerAddr, name); err != nil {
		return err
	}
	return nil
}

// verify that the inputs provided are correct
// signing address should match one of the input addresses
// generate confirmation signature for given utxo
func verifyAndSign(output store.Output, signerAddr ethcmn.Address, name string) error {
	sig, _ := clistore.GetSig(output.Position)
	inputAddrs := output.InputAddresses()

	if len(sig) == 130 || (len(sig) == 65 && len(inputAddrs) == 1) {
		return nil
	}

	for i, input := range inputAddrs {
		if input != signerAddr {
			continue
		}

		// get confirmation to generate signature
		fmt.Printf("\nUTXO\nPosition: %s\nOwner: 0x%x\nValue: %d\n", output.OutputPosition, output.Output.Owner, output.Output.Amount)
		buf := cosmoscli.BufferStdin()
		auth, err := cosmoscli.GetString(signPrompt, buf)
		if err != nil {
			return err
		}
		if auth != "Y" {
			return nil
		}

		hash := utils.ToEthSignedMessageHash(output.Tx.ConfirmationHash)
		sig, err := clistore.SignHashWithPassphrase(name, hash)
		if err != nil {
			return fmt.Errorf("failed to generate confirmation signature: { %s }", err)
		}

		if err := clistore.SaveSig(output.Output.Position, sig, i == 0); err != nil {
			return err
		}

		// print the results
		fmt.Printf("Confirmation Signature for output with position: %s\n", output.Output.Position)
		fmt.Printf("0x%x\n", sig)
	}
	return nil
}
