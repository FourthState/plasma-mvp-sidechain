package tx

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
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
	Use:   "broadcast-sigs <input1, input2>",
	Short: "Broadcast confirm signatures tied to input",
	Long: `Broadcasts confirm signatures tied to an input without spending funds. 
           Inputs take the format: 
           (blknum0.txindex0.oindex0.depositnonce0)::(blknum1.txindex1.oindex1.depositnonce1) 

    Example usage:
	plasmacli broadcast-sigs <input1>
	plasmacli broadcast-sigs <input1, input2>
	`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx := context.NewCLIContext()

		// parse inputs
		var inputStrs []string
		values := args[0]

		inputTokens := strings.Split(strings.TrimSpace(values), "::")
		if len(inputTokens) > 2 {
			return fmt.Errorf("2 or fewer inputs must be specified")
		}
		for _, token := range inputTokens {
			inputStrs = append(inputStrs, strings.TrimSpace(token))
		}

		// retrieve positions
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

		// build transaction
		// create non-nil inputs with null signatures
		input1 := plasma.Input{}
		if len(positions) >= 1 {
			input1.Position = positions[0]
			input1.Signature = [65]byte{1}
			input1.ConfirmSignatures = confirmSignatures[0]
		}

		input2 := plasma.Input{}
		if len(positions) >= 2 {
			input2.Position = positions[1]
			input2.Signature = [65]byte{1}
			input2.ConfirmSignatures = confirmSignatures[1]
		}

		cmd.SilenceUsage = true

		msg := msgs.ConfirmSigMsg{
			Input1: input1,
			Input2: input2,
		}
		if err := msg.ValidateBasic(); err != nil {
			return fmt.Errorf("failed on validating inputs. If you didn't provide the inputs please open an issue on github. error: %s", err)
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
