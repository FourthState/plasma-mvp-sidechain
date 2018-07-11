package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/ethereum/go-ethereum/common"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTo        = "to"
	flagPositions = "position"
	// ConfirmSigs possibly to be taken out
	flagConfirmSigs = "confirmSigs"
	flagAmounts     = "amounts"
)

func init() {
	rootCmd.AddCommand(sendTxCmd)
	sendTxCmd.Flags().String(flagTo, "", "Addresses sending to (separated by commas)")
	// Format for positions can be adjusted
	sendTxCmd.Flags().String(flagPositions, "", "UTXO Positions to be spent, format: blknum1.txindex1.oindex1.depositnonce1::blknum2.txindex2.oindex2.depositnonce2")
	// ConfirmSigs possibly to be taken out
	sendTxCmd.Flags().String(flagConfirmSigs, "", "Confirmation Signatures for inputs to be spent")
	sendTxCmd.Flags().String(flagAmounts, "", "Amounts to be spent, format: amount1, amount2, fee")
	sendTxCmd.Flags().String(client.FlagNode, "tcp://localhost:46657", "<host>:<port> to tendermint rpc interface for this chain")
	sendTxCmd.Flags().String(client.FlagAddress, "", "Address to sign with")
	viper.BindPFlags(sendTxCmd.Flags())
}

var sendTxCmd = &cobra.Command{
	Use:   "send",
	Short: "Build, Sign, and Send transactions",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewClientContextFromViper()

		// get the directory for our keystore
		dir := viper.GetString(FlagHomeDir)

		// get the from/to address
		from, err := ctx.GetInputAddresses(dir)
		if err != nil {
			return err
		}
		toStr := viper.GetString(flagTo)
		if toStr == "" {
			return errors.New("must provide an address to send to")
		}

		toAddrs := strings.Split(toStr, ",")
		if len(toAddrs) > 2 {
			return errors.New("incorrect amount of addresses provided")
		}

		var addr1, addr2 common.Address
		addr1, err = client.StrToAddress(toAddrs[0])
		if err != nil {
			return err
		}
		if len(toAddrs) > 1 && toAddrs[1] != "" {
			addr2, err = client.StrToAddress(toAddrs[1])
			if err != nil {
				return err
			}
		}

		csStr := viper.GetString(flagConfirmSigs)
		cs := strings.Split(csStr, ",")
		switch len(cs) {
		case 1:
			cs = append(cs, "", "", "")
		case 2:
			cs = append(cs, "", "")
		case 3:
			cs = append(cs, "")
		}
		confirmSigs1, err := getConfirmSigs(cs[0], cs[1])
		if err != nil {
			return err
		}
		confirmSigs2, err := getConfirmSigs(cs[2], cs[3])
		if err != nil {
			return err
		}

		// Get positions for transaction inputs
		posStr := viper.GetString(flagPositions)
		position, err := client.ParsePositions(posStr)
		if err != nil {
			return err
		}

		// Get amounts and fee
		amtStr := viper.GetString(flagAmounts)
		amounts, err := client.ParseAmounts(amtStr)
		msg := client.BuildMsg(from[0], from[1], addr1, addr2, position[0], position[1], confirmSigs1, confirmSigs2, amounts[0], amounts[1], amounts[2])
		res, err := ctx.SignBuildBroadcast(from, msg, dir)
		if err != nil {
			return err
		}
		fmt.Printf("Committed at block %d. Hash %s\n", res.Height, res.Hash.String())
		return nil
	},
}

func getConfirmSigs(sig1, sig2 string) (sigs [2]types.Signature, err error) {
	var cs1, cs2 []byte
	if sig1 != "" {
		cs1, err = hex.DecodeString(strings.TrimSpace(sig1))
		if err != nil {
			return sigs, err
		}
	}
	if sig2 != "" {
		cs2, err = hex.DecodeString(strings.TrimSpace(sig2))
		if err != nil {
			return sigs, err
		}
	}
	sigs = [2]types.Signature{types.Signature{cs1}, types.Signature{cs2}}
	return sigs, nil
}
