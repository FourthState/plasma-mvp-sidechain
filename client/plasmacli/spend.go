package main

import (
	"encoding/hex"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/store"
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
	spendCmd.Flags().String(positionF, "", "UTXO Positions to be spent, format: (blknum0.txindex0.oindex0.depositnonce0)::(blknum1.txindex1.oindex1.depositnonce1)")
	spendCmd.Flags().StringP(confirmSigs0F, "0", "", "Input Confirmation Signatures for first input to be spent (separated by commas)")
	spendCmd.Flags().StringP(confirmSigs1F, "1", "", "Input Confirmation Signatures for second input to be spent (separated by commas)")

	spendCmd.Flags().String(feeF, "0", "Fee to be spent")

	spendCmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	spendCmd.Flags().Bool(asyncF, false, "broadcast transactions asynchronously")
}

var spendCmd = &cobra.Command{
	Use:   "spend <to> <amount> <account>",
	Short: "Send a transaction spending utxos",
	Long: `Send a transaction spending from the specified account. If sending to multiple addresses, the account specified must contain exact utxo values.
If a single spending account is specified, leftover value from spending the utxo will be sent back to the account. 
In the case that the spending account does not have a large enough single utxo, two input utxos will be used. User can override retireved data with position and confirm signature flags.
<to> in the following usage is the address being sent the utxo amounts.

Usage:
	plasmacli <to> <amount> <account>
	plasmacli <to,to> <amount,amount> <account,account> --fee <fee>
	plasmacli <to> <amount> <account> --confirmSigs0 <signature> --confirmSig1 <signature>`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx := context.NewCLIContext()

		// parse accounts
		var accs []string
		names := args[2]

		accTokens := strings.Split(strings.TrimSpace(names), ",")
		if len(accTokens) == 0 || len(accTokens) > 2 {
			return fmt.Errorf("1 or 2 accounts must be specified")
		}
		for _, token := range accTokens {
			accs = append(accs, strings.TrimSpace(token))
		}

		// parse to addresses
		var toAddrs []ethcmn.Address
		toAddrTokens := strings.Split(strings.TrimSpace(args[0]), ",")
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

		// parse amounts
		var amounts []*big.Int // [amount0, amount1]
		var total int64
		amountTokens := strings.Split(strings.TrimSpace(args[1]), ",")
		if len(amountTokens) != 1 && len(amountTokens) != 2 {
			return fmt.Errorf("1 or 2 output amounts must be specified")
		}
		if len(amountTokens) != len(toAddrs) {
			return fmt.Errorf("provided amounts to not match the number of outputs")
		}
		for _, token := range amountTokens {
			token = strings.TrimSpace(token)
			num, ok := new(big.Int).SetString(token, 10)
			if !ok {
				return fmt.Errorf("failed to parsing amount: %s", token)
			}
			total += num.Int64()
			amounts = append(amounts, num)
		}

		// parse fee
		fee, ok := new(big.Int).SetString(strings.TrimSpace(viper.GetString(feeF)), 10)
		if !ok {
			return fmt.Errof("failed to parse fee: %s", fee)
		}

		total += fee.Int64()
		inputs, change, err := getInputs(total)
		if err != nil {
			return err
		}

		// get confirmation signatures from local storage
		var confirmSignatures [2][][65]byte
		for i, input := range inputs {
			sig, err := store.GetSig(input)
			if sig != nil {
				confirmSignatures[i] = sig
			}
		}

		confirmSignatures, inputs, overriden, err := parseSpend(confirmSignatures, inputs)
		if err != nil {
			return err
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
			if change == 0 || overriden {
				tx.Output1 = plasma.NewOutput(toAddrs[1], amounts[1])
			} else {
				return fmt.Errorf("cannot spend to two addresses since exact utxo inputs could not be found")
			}
		} else if change > 0 && !overriden {
			addr, err := store.GetAccount(accs[0])
			if err != nil {
				return err
			}
			tx.Output1 = plasma.NewOutput(addr, big.NewInt(change))
		} else {
			tx.Output1 = plasma.NewOutput(ethcmn.Address{}, nil)
		}
		tx.Fee = fee

		// create and fill in the signatures
		signer := accs[0]
		txHash := utils.ToEthSignedMessageHash(tx.TxHash())
		var signature [65]byte
		sig, err := store.SignHashWithPassphrase(signer, txHash)
		if err != nil {
			return err
		}
		copy(signature[:], sig)
		tx.Input0.Signature = signature
		if len(inputs) > 1 {
			if len(accs) > 2 {
				signer = accs[1]
			}
			sig, err := store.SignHashWithPassphrase(signer, txHash)
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
		if viper.GetBool(asyncF) {
			if _, err := ctx.BroadcastTxAsync(txBytes); err != nil {
				return err
			}
		} else {
			res, err := ctx.BroadcastTxAndAwaitCommit(txBytes)
			if err != nil {
				return err
			}
			fmt.Printf("Committed at block %d. Hash 0x%x\n", res.Height, res.TxHash)
		}

		return nil
	},
}

// attempt to retrieve to generate a valid spend transaction
// returns sum(inputs) - total
func getInputs(accs []string, total *big.Int) (inputs []plasma.Position, change *big.Int, err error) {
	ctx := context.NewCLIContext().WithCodec(codec.New()).WithTrustNode(true)
	change = 0 - total
	// must specifiy inputs if using two accounts
	if len(accs) > 1 {
		return inputs, change, nil
	}

	addr, err := store.GetAccount(accs[0])
	if err != nil {
		return inputs, change, err
	}

	res, err := ctx.QuerySubspace(addr.Bytes(), "utxo")
	if err != nil {
		return inputs, change, err
	}

	var optimalChange = total
	var position0, position1 plasma.Position
	// iterate through utxo's looking for optimal input pairing
	// return input pairing if input + input == total
	utxo0 := store.UTXO{}
	utxo1 := store.UTXO{}
	for i, outer := range res {
		if err := rlp.DecodeBytes(outer.Value, &utxo0); err != nil {
			return err
		}

		for k, inner := range res {
			// do not pair an input with itself
			if i == k {
				continue
			}

			if err := rlp.DecodeBytes(inner.Value, &utxo1); err != nil {
				return err
			}

			// check for exact match
			sum := new(big.Int).Add(utxo0.Output.Amount, utxo1.Output.Amount)
			if sum.Cmp(total) == 0 {
				inputs = append(utxo0.Position, utxo1.Position)
				return inputs, big.NewInt(0), nil
			}

			diff := new(big.Int).Sub(sum, total)
			if diff.Int64() > 0 && diff.Cmp(optimalChange) == -1 {
				optimalChange = diff
				position0 = utxo0.Position
				position1 = utxo1.Position
			}

		}
	}

	inputs = append(position0, position1)
	return inputs, optimalChange, nil
}

// Parse flags related to spending
// Flags override locally retrieved information
// bool returned specifies if inputs found were overriden
func parseSpend(confirmSignatures [2][][65]byte, inputs []plasma.Position) ([2][][65]byte, []plasma.Position, bool.error) {
	// validate confirm signatures
	for i := 0; i < 2; i++ {
		var flag string
		if i == 0 {
			flag = confirmSigs0F
		} else {
			flag = confirmSigs1F

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

	// parse inputs
	positions := strings.Split(strings.TrimSpace(viper.GetString(positionF)), "::")
	if len(positions) == 0 {
		if len(inputs) != 0 {
			return confirmSignatures, inputs, false, nil
		} else {
			return nil, nil, fmt.Errorf("must specifiy inputs to be used if two accounts are used")
		}
	}
	if len(positions) > 2 {
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

	return confirmSignatures, inputs, true, nil
}
