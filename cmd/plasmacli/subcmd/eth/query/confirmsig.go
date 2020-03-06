package query

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
)

// SigCmd returns the query confirm sig command
func SigCmd() *cobra.Command {
	config.AddPersistentTMFlags(sigCmd)
	return sigCmd
}

var sigCmd = &cobra.Command{
	Use:   "sig <position>",
	Short: "Query confirm signature information for a given position",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// parse position
		pos, err := plasma.FromPositionString(args[0])
		if err != nil {
			return fmt.Errorf("error parsing position: %s", err)
		}

		var sigs []byte
		ctx := context.NewCLIContext()
		sigs, err = client.ConfirmSignatures(ctx, pos)
		if err != nil {
			return fmt.Errorf("failed to retrieve confirm signature information: %s", err)
		}

		cmd.SilenceUsage = true

		switch len(sigs) {
		case 0:
			fmt.Printf("No Confirm Signatures Found")
		case 65:
			fmt.Printf("Confirmation Signatures: 0x%x\n", sigs[:])
		case 130:
			fmt.Printf("Confirmation Signatures: 0x%x, 0x%x\n", sigs[:65], sigs[65:])
		}

		return nil
	},
}