package eth

import (
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tm "github.com/tendermint/tendermint/rpc/core/types"
	"math/big"
	"strconv"
)

func init() {
	ethCmd.AddCommand(challengeCmd)
	challengeCmd.Flags().BoolP(trustNodeF, "t", false, "trust connected full node")
	challengeCmd.Flags().StringP(gasLimitF, "g", "21000", "gas limit for ethereum transaction")
	challengeCmd.Flags().String(ownerF, "", "owner of the challenging transaction, required if different from the specified account")
	challengeCmd.Flags().String(proofF, "", "merkle proof of inclusion")
	challengeCmd.Flags().StringP(sigsF, "S", "", "confirmation signatures for the challenging transaction")
	challengeCmd.Flags().StringP(txBytesF, "b", "", "bytes of the challenging transaction")
	viper.BindPFlags(challengeCmd.Flags())
}

var challengeCmd = &cobra.Command{
	Use:   "challenge <exiting position> <challenging position> <account>",
	Short: "Challenge an existing exit",
	Long:  ``,
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		// parse positions
		exitingPos, err := plasma.FromPositionString(args[0])
		if err != nil {
			return err
		}

		challengingPos, err := plasma.FromPositionString(args[1])
		if err != nil {
			return err
		}

		gasLimit, err := strconv.ParseUint(viper.GetString(gasLimitF), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse gas limit: { %s }", err)
		}

		key, err := ks.GetKey(args[2])
		if err != nil {
			return fmt.Errorf("failed to retrieve account key: { %s }", err)
		}

		var owner ethcmn.Address
		if viper.GetString(ownerF) != "" {
			owner = ethcmn.HexToAddress(viper.GetString(ownerF))
		} else {
			owner = crypto.PubkeyToAddress(key.PublicKey)
		}

		// bind key
		auth := bind.NewKeyedTransactor(key)
		transactOpts := &bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: gasLimit,
		}

		var txBytes, proof, confirmSignatures []byte
		if viper.GetBool(trustNodeF) {
			var result *tm.ResultTx
			result, confirmSignatures, err = getProof(owner, challengingPos)
			if err != nil {
				fmt.Errorf("failed to retrieve exit information: { %s }", err)
			}

			txBytes = result.Tx

			// flatten proof
			for _, aunt := range result.Proof.Proof.Aunts {
				proof = append(proof, aunt...)
			}
		}

		txBytes, proof, confirmSignatures, err = parseProof(txBytes, proof, confirmSignatures)

		exitPos := [4]*big.Int{exitingPos.BlockNum, big.NewInt(int64(exitingPos.TxIndex)), big.NewInt(int64(exitingPos.OutputIndex)), exitingPos.DepositNonce}
		challengePos := [2]*big.Int{challengingPos.BlockNum, big.NewInt(int64(challengingPos.TxIndex))}
		tx, err := rc.contract.ChallengeExit(transactOpts, exitPos, challengePos, txBytes, proof, confirmSignatures)
		if err != nil {
			return fmt.Errorf("failed to send challenge transaction: { %s }", err)
		}

		fmt.Printf("Sent challenge transaction\nTransaction Hash: 0x%s\n", tx.Hash())
		return nil
	},
}
