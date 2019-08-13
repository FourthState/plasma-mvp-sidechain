package eth

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	ks "github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcmn "github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tm "github.com/tendermint/tendermint/rpc/core/types"
	"math/big"
	"strconv"
)

// ExitCmd returns the eth exit command
func ExitCmd() *cobra.Command {
	config.AddPersistentTMFlags(exitCmd)
	exitCmd.Flags().String(feeF, "0", "fee committed in an unfinalized spend of the input")
	exitCmd.Flags().StringP(gasLimitF, "g", "300000", "gas limit for ethereum transaction")
	exitCmd.Flags().String(proofF, "", "merkle proof of inclusion")
	exitCmd.Flags().String(sigsF, "", "confirmation signatures for exiting utxo")
	exitCmd.Flags().Bool(useNodeF, false, "retrieve information from connected full node")
	exitCmd.Flags().String(txBytesF, "", "bytes of the transaction that created the utxo ")
	return exitCmd
}

var exitCmd = &cobra.Command{
	Use:   "exit <account> <position>",
	Short: "Start an exit for the given position",
	Long: `Starts an exit for the given position. If the trust-node flag is set, 
the necessary information will be retrieved from the connected full node. 
Otherwise, the transaction bytes, merkle proof, and confirmation signatures must be given. 
Usage of flags override information retrieved from full node. 

Deposit/Fee Exit Usage:
	plasmacli exit <account> <position>
	
Transaction Exit Usage:
	plasmacli exit <account> <position> --trust-node --gas-limit 30000
	plasmacli exit <account> <position> -t --fee <amount>
	plasmacli exit <account> <position> -b <tx-bytes> --proof <merkle-proof> -S <confirmation-signatures> --fee <amount>`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		viper.BindPFlags(cmd.Flags())
		var tx *eth.Transaction

		// parse position
		position, err := plasma.FromPositionString(args[1])
		if err != nil {
			return err
		}

		fee, err := strconv.ParseInt(viper.GetString(feeF), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse fee: { %s }", err)
		}

		gasLimit, err := strconv.ParseUint(viper.GetString(gasLimitF), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse gas limit: { %s }", err)
		}

		// retrieve account key
		key, err := ks.GetKey(args[0])
		if err != nil {
			return fmt.Errorf("failed to retrieve account key: { %s }", err)
		}

		// bind key, generate transact opts
		auth := bind.NewKeyedTransactor(key)
		transactOpts := &bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: gasLimit,
			Value:    big.NewInt(minExitBond), // minExitBond
		}

		// send fee exit
		if position.IsFee() {
			tx, err = plasmaContract.StartFeeExit(transactOpts, position.BlockNum, big.NewInt(fee))
			if err != nil {
				return fmt.Errorf("failed to start fee exit: { %s }", err)
			}
			fmt.Printf("Sent fee exit transaction\nTransaction Hash: 0x%x\n", tx.Hash())
			return nil
		}

		// send deposit exit
		if position.IsDeposit() {
			tx, err := plasmaContract.StartDepositExit(transactOpts, position.DepositNonce, big.NewInt(fee))
			if err != nil {
				return fmt.Errorf("failed to start deposit exit: { %s }", err)
			}
			fmt.Printf("Sent deposit exit transaction\nTransaction Hash: 0x%x\n", tx.Hash())
			return nil
		}

		// retrieve information necessary for transaction exit
		var txBytes, proof, confirmSignatures []byte
		if viper.GetBool(useNodeF) { // query full node
			var result *tm.ResultTx
			ctx := context.NewCLIContext()
			result, confirmSignatures, err = getProof(ctx, position)
			if err != nil {
				return fmt.Errorf("failed to retrieve exit information: { %s }", err)
			}

			txBytes = result.Tx

			// flatten proof
			for _, aunt := range result.Proof.Proof.Aunts {
				proof = append(proof, aunt...)
			}
		}

		if len(confirmSignatures) == 0 {
			sigs, err := ks.GetSig(position)
			if err == nil {
				confirmSignatures = sigs
			}
		}

		txBytes, proof, confirmSignatures, err = parseProof(txBytes, proof, confirmSignatures)
		if err != nil {
			return err
		}

		txPos := [3]*big.Int{position.BlockNum, big.NewInt(int64(position.TxIndex)), big.NewInt(int64(position.OutputIndex))}
		tx, err = plasmaContract.StartTransactionExit(transactOpts, txPos, txBytes, proof, confirmSignatures, big.NewInt(fee))
		if err != nil {
			return fmt.Errorf("failed to start transaction exit: { %s }", err)
		}
		fmt.Printf("Sent exit transaction\nTransaction Hash: 0x%x\n", tx.Hash())
		return nil
	},
}

// Parses flags related to proving exit/challenge
// Flags override full node information
// All necessary exit/challenge information is returned, or error is thrown
func parseProof(txBytes, proof, confirmSignatures []byte) ([]byte, []byte, []byte, error) {
	if viper.GetString(txBytesF) != "" {
		txBytes = ethcmn.FromHex(viper.GetString(txBytesF))
	}

	if viper.GetString(proofF) != "" {
		proof = ethcmn.FromHex(viper.GetString(proofF))
	}

	if viper.GetString(sigsF) != "" {
		confirmSignatures = ethcmn.FromHex(viper.GetString(sigsF))
	}

	// return error if information is missing
	if len(txBytes) != 811 {
		return txBytes, proof, confirmSignatures, fmt.Errorf("please provide txBytes with a length of 811 bytes. Current length: %d", len(txBytes))
	}

	if len(proof)%32 != 0 {
		return txBytes, proof, confirmSignatures, fmt.Errorf("please provide a merkle proof of inclusion for the given position. Proof must consist of 32 byte hashes")
	}

	if len(confirmSignatures)%65 != 0 {
		return txBytes, proof, confirmSignatures, fmt.Errorf("please provde confirmation signatures for the given position. Signatures must be 65 bytes in length")
	}

	if len(proof) == 0 {
		fmt.Println("Warning: No proof was found or provided. If the exiting transaction was not the only transaction included in the block, this transaction will fail.")
	}

	return txBytes, proof, confirmSignatures, nil
}
