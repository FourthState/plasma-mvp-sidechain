package eth

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/store"
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
	challengeCmd.Flags().StringP(gasLimitF, "g", "300000", "gas limit for ethereum transaction")
	challengeCmd.Flags().String(ownerF, "", "owner of the challenging transaction, required if different from the specified account")
	challengeCmd.Flags().String(proofF, "", "merkle proof of inclusion")
	challengeCmd.Flags().StringP(sigsF, "S", "", "confirmation signatures for the challenging transaction")
	challengeCmd.Flags().BoolP(trustNodeF, "t", false, "trust connected full node")
	challengeCmd.Flags().StringP(txBytesF, "b", "", "bytes of the challenging transaction")
}

var challengeCmd = &cobra.Command{
	Use:   "challenge <exiting position> <challenging position> <account>",
	Short: "Challenge an existing exit",
	Long: `Challenge a pending exit. If the trust-node flag is set, 
the necessary information will be retrieved from the connected full node. 
Otherwise, the transaction bytes, merkle proof, and confirmation signatures must be given. 
Usage of flags override information retrieved from full node. 

Usage:
	plasmacli eth challenge <exiting position> <challenging position> <account> --trust-node --gas-limit 30000
	plasmacli eth cahllenge <exiting position> <challenging position> <account> --proof <proof> --signatures <confirm signatures> --txBytes <challenge transaction bytes>`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
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

		key, err := store.GetKey(args[2])
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

		if len(confirmSignatures) == 0 {
			sigs, err := store.GetSig(challengingPos)
			if err == nil {
				confirmSignatures = sigs
			}
		}

		txBytes, proof, confirmSignatures, err = parseProof(txBytes, proof, confirmSignatures)
		if err != nil {
			return err
		}

		exitPos := [4]*big.Int{exitingPos.BlockNum, big.NewInt(int64(exitingPos.TxIndex)), big.NewInt(int64(exitingPos.OutputIndex)), exitingPos.DepositNonce}
		challengePos := [2]*big.Int{challengingPos.BlockNum, big.NewInt(int64(challengingPos.TxIndex))}
		tx, err := rc.contract.ChallengeExit(transactOpts, exitPos, challengePos, txBytes, proof, confirmSignatures)
		if err != nil {
			return fmt.Errorf("failed to send challenge transaction: { %s }", err)
		}

		fmt.Printf("Sent challenge transaction\nTransaction Hash: 0x%x\n", tx.Hash())
		return nil
	},
}
