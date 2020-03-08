package tx

import (
"fmt"
clistore "github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
"github.com/FourthState/plasma-mvp-sidechain/msgs"
"github.com/FourthState/plasma-mvp-sidechain/plasma"
"github.com/FourthState/plasma-mvp-sidechain/utils"
"github.com/cosmos/cosmos-sdk/client/context"
"github.com/ethereum/go-ethereum/rlp"
"github.com/spf13/cobra"
"github.com/spf13/viper"
"strings"
)

func BroadcastSigsCmd() *cobra.Command {
	includeCmd.Flags().Bool(asyncF, false, "wait for transaction commitment synchronously")
	return includeCmd
}

var broadcastSigsCmd = &cobra.Command{
	Use:   "broadcast-sigs <from, from> <input1, input2>",
	Short: "Broadcast confirm signatures",
	Long: `Broadcasts confirm signatures to network without spending funds. <from> must take form of account 
name, format: acc1::acc2. Inputs are UTXO Positions to be spent, format: (blknum0.txindex0.oindex0.depositnonce0)::(blknum1.txindex1.oindex1.depositnonce1) 

    Example usage:
	plasmacli broadcast-sigs <input1>
	plasmacli broadcast-sigs <input1, input2>
	plasmacli broadcast-sigs --confirmSigs0 <confirmSig> -confirmSigs1 <confirmSig>
	`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx := context.NewCLIContext()

		// parse accounts
		var accs []string
		names := args[0]

		accTokens := strings.Split(strings.TrimSpace(names), "::")
		if len(accTokens) == 0 || len(accTokens) > 2 {
			return fmt.Errorf("1 or 2 accounts must be specified")
		}
		for _, token := range accTokens {
			accs = append(accs, strings.TrimSpace(token))
		}

		// parse inputs
		var inputStrs []string
		values := args[1]

		inputTokens := strings.Split(strings.TrimSpace(values), "::")
		if len(inputTokens) > 2 {
			return fmt.Errorf("2 or fewer inputs must be specified")
		}
		for _, token := range inputTokens {
			inputStrs = append(inputStrs, strings.TrimSpace(token))
		}

		var positions []plasma.Position
		for _, token := range inputStrs {
			position, err := plasma.FromPositionString(token)
			if err != nil {
				return fmt.Errorf("error parsing position from string: %s", err)
			}
			positions = append(positions, position)
		}

		// get confirmation signatures from local storage
		confirmSignatures := getConfirmSignatures(positions)

		// override retrieved signatures if provided through flags - check if parseConfirmSignatures is right tool
		confirmSignatures, err := parseConfirmSignatures(confirmSignatures)
		if err != nil {
			return fmt.Errorf("error retrieving confirm Signatures %s", err)
		}

		// build transaction
		// create non-nil inputs with signatures
		input1 := plasma.Input{}
		if len(positions) > 0 {
			signer := accs[0]
			positionHash := utils.ToEthSignedMessageHash(positions[0].Bytes())
			var signature [65]byte
			sig, err := clistore.SignHashWithPassphrase(signer, positionHash)
			if err != nil {
				return err
			}
			copy(signature[:], sig)
			input1.Position = positions[0]
			input1.Signature = signature
			input1.ConfirmSignatures = confirmSignatures[0]
		}

		input2 := plasma.Input{}
		if len(positions) > 1 {
			signer := accs[1]
			positionHash := utils.ToEthSignedMessageHash(positions[1].Bytes())
			var signature [65]byte
			sig, err := clistore.SignHashWithPassphrase(signer, positionHash)
			if err != nil {
				return err
			}
			copy(signature[:], sig)

			input2.Position = positions[1]
			input2.Signature = signature
			input2.ConfirmSignatures = confirmSignatures[1]
		}

		cmd.SilenceUsage = true

		msg := msgs.ConfirmSigMsg{
			Input1: input1,
			Input2: input2,
		}
		if err := msg.ValidateBasic(); err != nil {
			return fmt.Errorf("failed on validating inputs. If you didn't provide the inputs please open an issue on github. Error: { %s }", err)
		}

		txBytes, err := rlp.EncodeToBytes(&msg)
		if err != nil {
			return err
		}

		// broadcast to the node
		if viper.GetBool(asyncF) {
			if _, err := ctx.BroadcastTxAsync(txBytes); err != nil {
				return err
			}
		} else {
			res, err := ctx.BroadcastTxAndAwaitCommit(txBytes)
			if err != nil {
				return err
			}
			fmt.Printf("Committed at block %d. Hash %s\n", res.Height, res.TxHash)
		}

		return nil
	},
}

