package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTo           = "to"
	flagPositions    = "position"
	flagConfirmSigs0 = "Input0ConfirmSigs"
	flagConfirmSigs1 = "Input1ConfirmSigs"
	flagAmounts      = "amounts"
)

func init() {
	rootCmd.AddCommand(sendTxCmd)
	sendTxCmd.Flags().String(flagTo, "", "Addresses sending to (separated by commas)")
	// Format for positions can be adjusted
	sendTxCmd.Flags().String(flagPositions, "", "UTXO Positions to be spent, format: blknum0.txindex0.oindex0.depositnonce0::blknum1.txindex1.oindex1.depositnonce1")

	sendTxCmd.Flags().String(flagConfirmSigs0, "", "Input Confirmation Signatures for first input to be spent (separated by commas)")
	sendTxCmd.Flags().String(flagConfirmSigs1, "", "Input Confirmation Signatures for second input to be spent (separated by commas)")

	sendTxCmd.Flags().String(flagAmounts, "", "Amounts to be spent, format: amount1, amount2, fee")

	sendTxCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
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

		csStr := viper.GetString(flagConfirmSigs0)
		cs0 := strings.Split(csStr, ",")
		csStr = viper.GetString(flagConfirmSigs1)
		cs1 := strings.Split(csStr, ",")

		confirmSigs0, err := getConfirmSigs(cs0)
		if err != nil {
			return err
		}
		confirmSigs1, err := getConfirmSigs(cs1)
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
		if utils.ZeroAddress(addr2) && amounts[1] != 0 {
			return fmt.Errorf("You are trying to send %d amount to the nil address. Please input the zero address if you would like to burn your amount", amounts[1])
		}
		msg := client.BuildMsg(from[0], from[1], addr1, addr2, position[0], position[1], confirmSigs0, confirmSigs1, amounts[0], amounts[1], amounts[2])
		res, err := ctx.SignBuildBroadcast(from, msg, dir)
		if err != nil {
			return err
		}
		fmt.Printf("Committed at block %d. Hash %s\n", res.Height, res.Hash.String())
		return nil
	},
}

func getConfirmSigs(sigs []string) (confirmSigs [][65]byte, err error) {
	var cs0, cs1 []byte
	var confirmSig0, confirmSig1 [65]byte

	if strings.Compare(sigs[0], "") == 0 {
		return confirmSigs, nil
	}
	switch len(sigs) {
	case 1:
		if cs0, err = hex.DecodeString(strings.TrimSpace(sigs[0])); err != nil {
			return confirmSigs, err
		}
		copy(confirmSig0[:], cs0)
		return append(confirmSigs, confirmSig0), nil
	case 2:
		if cs0, err = hex.DecodeString(strings.TrimSpace(sigs[0])); err != nil {
			return confirmSigs, err
		}
		if cs1, err = hex.DecodeString(strings.TrimSpace(sigs[1])); err != nil {
			return confirmSigs, err
		}
		copy(confirmSig0[:], cs0)
		copy(confirmSig1[:], cs1)

		return append(confirmSigs, confirmSig0, confirmSig1), nil
	}
	return confirmSigs, errors.New("the provided confirmSigs caused undefined behavior. Pass in 0, 1 or 2 confirm sigs per flag")
}
