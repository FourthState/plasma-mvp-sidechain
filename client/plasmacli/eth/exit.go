package eth

import (
	"fmt"
	"math/big"
	"strconv"

	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	ethCmd.AddCommand(exitCmd)
	exitCmd.Flags().String(feeF, "0", "fee committed in an unfinalized spend of the input")
	exitCmd.Flags().BoolP(trustNodeF, "t", false, "trust connected full node")
	exitCmd.Flags().StringP(txBytesF, "b", "", "bytes of the transaction that created the utxo ")
	exitCmd.Flags().String(proofF, "", "merkle proof of inclusion")
	exitCmd.Flags().StringP(sigsF, "S", "", "confirmation signatures for the utxo")
	viper.BindPFlags(exitCmd.Flags())
}

var exitCmd = &cobra.Command{
	Use:   "exit <account> <position>",
	Short: "Start an exit for the given position",
	Long: `Starts an exit for the given position. If the trust-node flag is set, 
the necessary information will be retrieved from the connected full node. 
Otherwise, the transaction bytes, merkle proof, and confirmation signatures must be given

Usage:
	plasmacli exit <account> <position> -t
	plasmacli exit <account> <position> -t --fee <amount>
	plasmacli exit <account> <position> -b <tx-bytes> --proof <merkle-proof> -S <confirmation-signatures> --fee <amount>`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		key, err := ks.GetKey(args[0])
		if err != nil {
			return err
		}
		addr := crypto.PubkeyToAddress(key.PublicKey)

		// parse position
		position, err := plasma.FromPositionString(args[1])
		if err != nil {
			return err
		}

		fee, err := strconv.ParseInt(viper.GetString(feeF), 10, 64)
		if err != nil {
			//return fmt.Errorf("failed to parse fee - %v", err)
			fee = 0
		}

		var txBytes, proof, confirmSignatures []byte
		if viper.GetBool(trustNodeF) { // query full node
			txBytes, proof, confirmSignatures, err = proveExit(addr, position)
			if err != nil {
				return err
			}
		} else { // use inputted flags
			txBytesStr := viper.GetString(txBytesF)
			proofStr := viper.GetString(proofF)
			sigsStr := viper.GetString(sigsF)

			if txBytesStr != "" {
				txBytes = []byte(txBytesStr)
			} else {
				return fmt.Errorf("please provide transaction bytes for the given position")
			}

			if proofStr != "" {
				proof = []byte(proofStr)
			} else {
				return fmt.Errorf("please provide a merkle proof for the given position")
			}

			if sigsStr != "" {
				confirmSignatures = []byte(sigsStr)
			} else {
				return fmt.Errorf("please provide confirmation signatures for the given position")
			}
		}

		// Send exit transaction
		auth := bind.NewKeyedTransactor(key)
		defer func() {
			rc.session.TransactOpts = bind.TransactOpts{}
		}()
		rc.session.TransactOpts = bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: 3141592,
			Value:    big.NewInt(minExitBond),
		}

		// TODO: Add start fee exit via pos.IsFee()
		if position.IsDeposit() {
			if _, err := rc.session.StartDepositExit(position.DepositNonce, big.NewInt(fee)); err != nil {
				return err
			}
			fmt.Println("Started deposit exit")
		} else {
			txPos := [3]*big.Int{position.BlockNum, big.NewInt(int64(position.TxIndex)), big.NewInt(int64(position.OutputIndex))}
			if _, err := rc.session.StartTransactionExit(txPos, txBytes, proof, confirmSignatures, big.NewInt(fee)); err != nil {
				return err
			}
			fmt.Println("Started transaction exit")
		}

		return nil
	},
}
