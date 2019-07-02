package tx

import (
	"fmt"
	clistore "github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	cosmoscli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

const (
	signPrompt = "Would you like to finalize this transaction? [Y/n]"
)

func SignCmd() *cobra.Command {
	signCmd.Flags().String(ownerF, "", "Owner of the output (required with position flag)")
	signCmd.Flags().String(positionF, "", "Position of transaction to finalize (required with owner flag)")
	return signCmd
}

var signCmd = &cobra.Command{
	Use:   "sign <account>",
	Short: "Sign confirmation signatures for pending transactions",
	Long: `Iterate over all unfinalized transaction corresponding to the provided account. 
Prompt the user for confirmation to finailze the pending transactions. Owner and Position flags can be used to finalize a specific transaction.

Usage:
	plasmacli sign <account>
	plasmacli sign <account> --position "(blknum.txindex.oindex.depositnonce)"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx := context.NewCLIContext()

		name := args[0]

		signerAddr, err := clistore.GetAccount(name)
		if err != nil {
			return err
		}

		positionS := viper.GetString(positionF)
		if positionS != "" {
			position, err := plasma.FromPositionString(strings.TrimSpace(viper.GetString(positionF)))
			if err != nil {
				return err
			}

			err = signSingleConfirmSig(ctx, position, signerAddr, name)
			return err
		}

		utxos, err := query.Info(ctx, signerAddr)
		if err != nil {
			return err
		}

		for _, output := range utxos {

			if output.Output.Spent {
				tx, err := query.Tx(ctx, output.Output.SpenderTx)
				if err != nil {
					return err
				}

				for i, pos := range tx.InputPositions() {
					err = signSingleConfirmSig(ctx, pos, signerAddr, name)
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
// verify that the inputs provided are correct
// signing address should match one of the input addresses
// generate confirmation signature for given utxo
func signSingleConfirmSig(ctx context.CLIContext, position plasma.Position, signerAddr ethcmn.Address, name string) error {
	// query for output for the specified position
	output, err := query.OutputInfo(ctx, position)
	if err != nil {
		return err
	}

	sig, _ := clistore.GetSig(output.Tx.Position)
	inputAddrs := output.Tx.InputAddresses()

	if len(sig) == 130 || (len(sig) == 65 && len(inputAddrs) == 1) {
		return nil
	}

	for i, input := range inputAddrs {
		if input != signerAddr {
			continue
		}
		// TODO: fix to use correct position
		// get confirmation to generate signature
		fmt.Printf("\nUTXO\nPosition: %s\nOwner: 0x%x\nValue: %d\n", output.Tx.Position, output.Output.Output.Owner, output.Output.Amount)
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

		if err := clistore.SaveSig(output.Tx.Position, sig, i == 0); err != nil {
			return err
		}

		// print the results
		fmt.Printf("Confirmation Signature for output with position: %s\n", output.Tx.Position)
		fmt.Printf("0x%x\n", sig)
	}
	return nil
}