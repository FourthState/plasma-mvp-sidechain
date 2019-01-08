package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"strings"
)

const (
	flagTo           = "to"
	flagPositions    = "position"
	flagConfirmSigs0 = "Input0ConfirmSigs"
	flagConfirmSigs1 = "Input1ConfirmSigs"
	flagAmounts      = "amounts"
	flagAddress      = "address"
	flagSync         = "sync"
)

func init() {
	rootCmd.AddCommand(sendTxCmd)
	sendTxCmd.Flags().String(flagTo, "", "Addresses sending to (separated by commas)")
	sendTxCmd.Flags().String(flagPositions, "", "UTXO Positions to be spent, format: (blknum0.txindex0.oindex0.depositnonce0)::(blknum1.txindex1.oindex1.depositnonce1)")

	sendTxCmd.Flags().String(flagConfirmSigs0, "", "Input Confirmation Signatures for first input to be spent (separated by commas)")
	sendTxCmd.Flags().String(flagConfirmSigs1, "", "Input Confirmation Signatures for second input to be spent (separated by commas)")

	sendTxCmd.Flags().String(flagAmounts, "", "Amounts to be spent, format: amount1, amount2, fee")
	sendTxCmd.Flags().String(flagAddress, "", "Addresses to sign with. One address will be used to sign both inputs if two addresses are not provided (seperated by commas)")

	sendTxCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	sendTxCmd.Flags().BoolP(flagSync, "s", false, "wait for transaction commitment synchronously")
	viper.BindPFlags(sendTxCmd.Flags())
}

var sendTxCmd = &cobra.Command{
	Use:   "send",
	Short: "Build, Sign, and Send transactions",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()

		// validate addresses
		var fromAddrs, toAddrs, signers []common.Address
		fromAddrTokens := strings.Split(strings.TrimSpace(viper.GetString(flagFrom)), ",")
		toAddrTokens := strings.Split(strings.TrimSpace(viper.GetString(flagTo)), ",")
		signerAddrTokens := strings.Split(strings.TrimSpace(viper.GetString(flagAddress)), ",")
		if len(fromAddrTokens) == 0 || len(toAddrTokens) == 0 {
			return fmt.Errorf("at least one input and one output must be specified")
		}
		if len(signerAddrTokens) == 0 || len(signerAddrTokens) > 2 {
			return fmt.Errorf("at least 1 or 2 signers must be provided. Same signer will be used for both inputs if 1 is provided")
		}
		if len(fromAddrTokens) > 2 || len(toAddrTokens) > 2 {
			return fmt.Errorf("can only spend at most 2 inputs to at most 2 outputs")
		}
		for _, token := range fromAddrTokens {
			token := strings.TrimSpace(token)
			if !common.IsHexAddress(token) {
				return fmt.Errorf("invalid address provided. please use hex format")
			}
			fromAddrs = append(fromAddrs, common.HexToAddress(token))
		}
		for _, token := range toAddrTokens {
			token := strings.TrimSpace(token)
			if !common.IsHexAddress(token) {
				return fmt.Errorf("invalid address provided. please use hex format")
			}
			addr := common.HexToAddress(token)
			if utils.IsZeroAddress(addr) {
				return fmt.Errorf("cannot spend to the zero address")
			}
			toAddrs = append(toAddrs, addr)
		}
		for _, token := range signerAddrTokens {
			token := strings.TrimSpace(token)
			if !common.IsHexAddress(token) {
				return fmt.Errorf("invalid address provided. please use hex format")
			}
			signers = append(signers, common.HexToAddress(token))
		}

		// validate confirm signatures
		var confirmSignatures [2][][65]byte
		for i := 0; i < 2; i++ {
			var flag string
			if i == 0 {
				flag = flagConfirmSigs0
			} else {
				flag = flagConfirmSigs1

			}
			confirmSigTokens := strings.Split(strings.TrimSpace(viper.GetString(flag)), ",")
			if len(confirmSigTokens) > 2 {
				return fmt.Errorf("only pass in 0, 1 or 2, confirm signatures")
			}

			var confirmSignature [][65]byte
			for _, token := range confirmSigTokens {
				token := strings.TrimSpace(token)
				if len(token) != 65 {
					return fmt.Errorf("signatures must be of length 65 bytes")
				}
				if !common.IsHexAddress(token) {
					return fmt.Errorf("invalid first confirm signature provided. please use hex format")
				}

				var signature [65]byte
				sig, err := hex.DecodeString(token)
				if err != nil {
					return err
				}

				copy(signature[:], sig)
				confirmSignature = append(confirmSignature, signature)
			}

			confirmSignatures[i] = confirmSignature
		}

		// validate inputs
		var inputs []plasma.Position
		positions := strings.Split(strings.TrimSpace(viper.GetString(flagPositions)), "::")
		if len(positions) > 2 || len(positions) == 0 {
			return fmt.Errorf("only pass in 1 or 2 positions")
		}
		for _, token := range positions {
			token = strings.TrimSpace(token)
			position, err := plasma.FromPositionString(token)
			if err != nil {
				return err
			}

			inputs = append(inputs, position)
		}

		// validate amounts and fee
		var amounts []*big.Int // [ammount0, amount1, fee]
		amountTokens := strings.Split(strings.TrimSpace(viper.GetString(flagAmounts)), ",")
		if len(amountTokens) != 3 {
			return fmt.Errorf("3 amounts must be passed in. amount0, amount1, fee")
		}
		if len(amountTokens)-1 != len(toAddrs) {
			return fmt.Errorf("provided amounts to not match the number of outputs")
		}
		for _, token := range amountTokens {
			token = strings.TrimSpace(token)
			num, ok := new(big.Int).SetString(token, 10)
			if !ok {
				return fmt.Errorf("error parsing number: %s", token)
			}
			amounts = append(amounts, num)
		}

		// create the transaction without signatures
		tx := plasma.Transaction{}
		tx.Input0 = plasma.NewInput(inputs[0], fromAddrs[0], [65]byte{}, confirmSignatures[0])
		if len(inputs) > 1 {
			tx.Input1 = plasma.NewInput(inputs[1], fromAddrs[1], [65]byte{}, confirmSignatures[1])
		}
		tx.Output0 = plasma.NewOutput(toAddrs[0], amounts[0])
		if len(toAddrs) > 1 {
			tx.Output1 = plasma.NewOutput(toAddrs[1], amounts[1])
		}
		tx.Fee = amounts[2]

		// create and fill in the signatures
		txHash := tx.TxHash()
		var signature [65]byte
		sig, err := keystore.SignHashWithPassphrase(signers[0], txHash[:])
		if err != nil {
			return err
		}
		copy(signature[:], sig)
		tx.Input0.Signature = signature
		if len(inputs) > 1 {
			var signer common.Address
			if len(signers) > 1 {
				signer = signers[1]
			} else {
				signer = signers[0]
			}
			sig, err := keystore.SignHashWithPassphrase(signer, txHash[:])
			if err != nil {
				return err
			}
			copy(signature[:], sig)
			tx.Input1.Signature = signature
		}

		// create SpendMsg and txBytes
		msg := msgs.SpendMsg{
			Transaction: tx,
		}
		txBytes, err := rlp.EncodeToBytes(&msg)
		if err != nil {
			return nil
		}

		// broadcast to the node
		if viper.GetBool(flagSync) {
			res, err := ctx.BroadcastTxAndAwaitCommit(txBytes)
			if err != nil {
				return err
			}
			fmt.Printf("Committed at block %d. Hash %s\n", res.Height, res.Hash.String())
		} else {
			if _, err := ctx.BroadcastTxAsync(txBytes); err != nil {
				return err
			}
		}

		return nil
	},
}
