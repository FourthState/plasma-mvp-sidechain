package cmd

import (
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

const (
	flagTo        = "to"
	flagPositions = "position"
)

func init() {
	rootCmd.AddCommand(sendTxCmd)
	sendTxCmd.Flags().String(flagTo, "", "Addresses sending to (separated by commas)")
	sendTxCmd.Flags().String(flagPositions, "", "UTXO Positions to be spent")
}

var sendTxCmd = &cobra.Command{
	Use:   "send",
	Short: "Build, Sign, and Send transactions",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewClientContextFromViper()

		// get the from/to address
		from, err := ctx.GetFromAddress()
		if err != nil {
			return err
		}

		toStr := viper.GetString(flagTo)

		toAddrs := strings.Split(toStr, ",")
		if len(toAddrs) > 2 || len(toAddrs) == 0 {
			return errors.New("incorrect amount of addresses provided")
		}

		// Assert that addresses are convertible to hex address
		if !common.IsHexAddress(toAddrs[0]) || (toAddrs[1] && !common.IsHexAddress(toAddrs[1])) {
			return errors.New("address cannot be converted to hex address")
		}

		var addr1, addr2 common.Address
		addr1 = common.HexToAddress(toAddrs[0])
		if toAddrs[1] {
			addr2 = common.HexToAddress(toAddrs[2])
		}

		// Get positions, amounts, fee

		msg := client.BuildMsg(from, addr1, addr2, position1, position2, inputPosition1, inputPosition2, amount1, amount2, fee)
		res, err := ctx.SignBuildBroadcast(ctx.FromAddressName, msg)
		if err != nil {
			return err
		}
		fmt.Printf("Committed at block %d. Hash %s\n", res.Height, res.Hash.String())
		return nil
	},
}
