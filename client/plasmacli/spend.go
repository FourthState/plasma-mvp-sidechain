package main

import (
	"encoding/hex"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/keys"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"strings"
)

func init() {
	rootCmd.AddCommand(spendCmd)
	spendCmd.Flags().String(flagTo, "", "Addresses sending to (separated by commas)")
	spendCmd.Flags().String(flagPositions, "", "UTXO Positions to be spent, format: (blknum0.txindex0.oindex0.depositnonce0)::(blknum1.txindex1.oindex1.depositnonce1)")

	spendCmd.Flags().String(flagConfirmSigs0, "", "Input Confirmation Signatures for first input to be spent (separated by commas)")
	spendCmd.Flags().String(flagConfirmSigs1, "", "Input Confirmation Signatures for second input to be spent (separated by commas)")

	spendCmd.Flags().String(flagAmounts, "", "Amounts to be spent, format: amount1, amount2, fee")
	spendCmd.Flags().StringP(flagAccount, "a", "", "Account to sign with")

	spendCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	spendCmd.Flags().BoolP(flagSync, "s", false, "wait for transaction commitment synchronously")
	viper.BindPFlags(spendCmd.Flags())
}

var spendCmd = &cobra.Command{
	Use:   "spend",
	Short: "Send a transaction spending utxos",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()
		name := viper.GetString(flagAccount)

		db, signer, err := keys.OpenAndGet(name)
		if err != nil {
			return err
		}
		defer db.Close()

		// validate addresses
		var toAddrs []ethcmn.Address
		toAddrTokens := strings.Split(strings.TrimSpace(viper.GetString(flagTo)), ",")
		if len(toAddrTokens) == 0 || len(toAddrTokens) > 2 {
			return fmt.Errorf("1 or 2 outputs must be specified")
		}

		for _, token := range toAddrTokens {
			token := strings.TrimSpace(token)
			if !ethcmn.IsHexAddress(token) {
				return fmt.Errorf("invalid address provided. please use hex format")
			}
			addr := ethcmn.HexToAddress(token)
			if utils.IsZeroAddress(addr) {
				return fmt.Errorf("cannot spend to the zero address")
			}
			toAddrs = append(toAddrs, addr)
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
			// empty confirmsig
			if len(confirmSigTokens) == 1 && confirmSigTokens[0] == "" {
				continue
			} else if len(confirmSigTokens) > 2 {
				return fmt.Errorf("only pass in 0, 1 or 2, confirm signatures")
			}

			var confirmSignature [][65]byte
			for _, token := range confirmSigTokens {
				token := strings.TrimSpace(token)
				sig, err := hex.DecodeString(token)
				if err != nil {
					return err
				}
				if len(sig) != 65 {
					return fmt.Errorf("signatures must be of length 65 bytes")
				}

				var signature [65]byte
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
		var amounts []*big.Int // [amount0, amount1, fee]
		amountTokens := strings.Split(strings.TrimSpace(viper.GetString(flagAmounts)), ",")
		if len(amountTokens) != 2 && len(amountTokens) != 3 {
			return fmt.Errorf("number of amounts must equal the number of outputs in addition to the fee")
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
		tx.Input0 = plasma.NewInput(inputs[0], [65]byte{}, confirmSignatures[0])
		if len(inputs) > 1 {
			tx.Input1 = plasma.NewInput(inputs[1], [65]byte{}, confirmSignatures[1])
		} else {
			tx.Input1 = plasma.NewInput(plasma.NewPosition(nil, 0, 0, nil), [65]byte{}, nil)
		}
		tx.Output0 = plasma.NewOutput(toAddrs[0], amounts[0])
		if len(toAddrs) > 1 {
			tx.Output1 = plasma.NewOutput(toAddrs[1], amounts[1])
			tx.Fee = amounts[2]
		} else {
			tx.Output1 = plasma.NewOutput(ethcmn.Address{}, nil)
			tx.Fee = amounts[1]
		}

		// create and fill in the signatures
		txHash := utils.ToEthSignedMessageHash(tx.TxHash())
		var signature [65]byte
		sig, err := keystore.SignHashWithPassphrase(signer, txHash)
		if err != nil {
			return err
		}
		copy(signature[:], sig)
		tx.Input0.Signature = signature
		if len(inputs) > 1 {
			sig, err := keystore.SignHashWithPassphrase(signer, txHash)
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
		if err := msg.ValidateBasic(); err != nil {
			return err
		}

		txBytes, err := rlp.EncodeToBytes(&msg)
		if err != nil {
			return err
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
