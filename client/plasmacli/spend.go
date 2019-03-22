package main

import (
	"encoding/hex"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/plasmacli/eth"
	clistore "github.com/FourthState/plasma-mvp-sidechain/client/store"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
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

		toAddrs, err := parseToAddresses(args[0])
		if err != nil {
			return err
		}

		amounts, fee, total, err := parseAmounts(args[1], toAddrs)
		if err != nil {
			return err
		}

		inputs, err := parseInputs()
		if err != nil {
			return err
		}

		change := new(big.Int)
		if len(inputs) == 0 {
			inputs, change, err = retrieveInputs(accs, total)
			if err != nil {
				return err
			}
		}

		// get confirmation signatures from local storage
		confirmSignatures := getConfirmSignatures(inputs)

		// override retireved signatures if provided through flags
		confirmSignatures, err = parseConfirmSignatures(confirmSignatures)
		if err != nil {
			return err
		}

		// build transaction
		// create the inputs without signatures
		tx := plasma.Transaction{}
		tx.Input0 = plasma.NewInput(inputs[0], [65]byte{}, confirmSignatures[0])
		if len(inputs) > 1 {
			tx.Input1 = plasma.NewInput(inputs[1], [65]byte{}, confirmSignatures[1])
		} else {
			tx.Input1 = plasma.NewInput(plasma.NewPosition(nil, 0, 0, nil), [65]byte{}, nil)
		}

		// generate outputs
		// use change to determine outcome of second output
		tx.Output0 = plasma.NewOutput(toAddrs[0], amounts[0])
		if len(toAddrs) > 1 {
			if change.Sign() == 0 {
				tx.Output1 = plasma.NewOutput(toAddrs[1], amounts[1])
			} else {
				return fmt.Errorf("cannot spend to two addresses since exact utxo inputs could not be found")
			}
		} else if change.Sign() == 1 {
			addr, err := clistore.GetAccount(accs[0])
			if err != nil {
				return err
			}
			tx.Output1 = plasma.NewOutput(addr, change)
		} else {
			tx.Output1 = plasma.NewOutput(ethcmn.Address{}, nil)
		}
		tx.Fee = fee

		// create and fill in the signatures
		signer := accs[0]
		txHash := utils.ToEthSignedMessageHash(tx.TxHash())
		var signature [65]byte
		sig, err := clistore.SignHashWithPassphrase(signer, txHash)
		if err != nil {
			return err
		}
		copy(signature[:], sig)
		tx.Input0.Signature = signature
		if len(inputs) > 1 {
			if len(accs) > 2 {
				signer = accs[1]
			}
			sig, err := clistore.SignHashWithPassphrase(signer, txHash)
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
			return fmt.Errorf("failed on validating transaction. If you didn't provide the inputs please open an issue on github. Error: { %s }", err)
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

// Retrieve confirmation signatures from local storage if they exist
func getConfirmSignatures(inputs []plasma.Position) (confirmSignatures [2][][65]byte) {
	for i, input := range inputs {
		sig, _ := clistore.GetSig(input)
		if sig != nil {
			var sigs [][65]byte
			switch len(sig) {
			case 65:
				var s [65]byte
				copy(s[:], sig[:])
				sigs = append(sigs, s)
			case 130:
				var s [65]byte
				copy(s[:], sig[:65])
				sigs = append(sigs, s)
				copy(s[:], sig[65:])
				sigs = append(sigs, s)
			}
			confirmSignatures[i] = sigs
		}
	}
	return confirmSignatures
}

// parses input amounts and fee
// amounts - [amount0, amount1]
func parseAmounts(amtArgs string, toAddrs []ethcmn.Address) (amounts []*big.Int, fee, total *big.Int, err error) {
	total = new(big.Int)
	amountTokens := strings.Split(strings.TrimSpace(amtArgs), ",")
	if len(amountTokens) != 1 && len(amountTokens) != 2 {
		return amounts, fee, total, fmt.Errorf("1 or 2 output amounts must be specified")
	}

	if len(amountTokens) != len(toAddrs) {
		return amounts, fee, total, fmt.Errorf("provided amounts to not match the number of outputs")
	}

	for _, token := range amountTokens {
		token = strings.TrimSpace(token)
		num, ok := new(big.Int).SetString(token, 10)
		if !ok {
			return amounts, fee, total, fmt.Errorf("failed to parsing amount: %s", token)
		}
		total.Add(total, num)
		amounts = append(amounts, num)
	}

	var ok bool
	fee, ok = new(big.Int).SetString(strings.TrimSpace(viper.GetString(feeF)), 10)
	if !ok {
		return amounts, fee, total, fmt.Errorf("failed to parse fee: %s", fee)
	}
	total.Add(total, fee)

	return amounts, fee, total, nil
}

// parse confirmation signatures passed in through flags
func parseConfirmSignatures(confirmSignatures [2][][65]byte) ([2][][65]byte, error) {
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
			return confirmSignatures, fmt.Errorf("only pass in 0, 1 or 2, confirm signatures")
		}

		var confirmSignature [][65]byte
		for _, token := range confirmSigTokens {
			token := strings.TrimSpace(token)
			sig, err := hex.DecodeString(token)
			if err != nil {
				return confirmSignatures, err
			}
			if len(sig) != 65 {
				return confirmSignatures, fmt.Errorf("signatures must be of length 65 bytes")
			}

			var signature [65]byte
			copy(signature[:], sig)
			confirmSignature = append(confirmSignature, signature)
		}

		confirmSignatures[i] = confirmSignature
	}
	return confirmSignatures, nil
}

// parse inputs passed in through flags
// Split will return a slice of at least length 1
func parseInputs() (inputs []plasma.Position, err error) {
	positions := strings.Split(strings.TrimSpace(viper.GetString(positionF)), "::")
	if len(positions) == 1 && len(positions[0]) == 0 {
		return inputs, err
	}

	if len(positions) > 2 {
		return inputs, fmt.Errorf("only pass in 1 or 2 positions")
	}

	for _, token := range positions {
		token = strings.TrimSpace(token)
		position, err := plasma.FromPositionString(token)
		if err != nil {
			return inputs, err
		}
		inputs = append(inputs, position)
	}

	return inputs, nil
}

// parse the passed in addresses that will be sent to
func parseToAddresses(addresses string) (toAddrs []ethcmn.Address, err error) {
	toAddrTokens := strings.Split(strings.TrimSpace(addresses), ",")
	if len(toAddrTokens) == 0 || len(toAddrTokens) > 2 {
		return toAddrs, fmt.Errorf("1 or 2 outputs must be specified")
	}

	for _, token := range toAddrTokens {
		token := strings.TrimSpace(token)
		if !ethcmn.IsHexAddress(token) {
			return toAddrs, fmt.Errorf("invalid address provided. please use hex format")
		}

		addr := ethcmn.HexToAddress(token)
		if utils.IsZeroAddress(addr) {
			return toAddrs, fmt.Errorf("cannot spend to the zero address")
		}
		toAddrs = append(toAddrs, addr)
	}

	return toAddrs, nil
}

// attempt to retrieve inputs to generate a valid spend transaction
// returns inputs and sum(inputs) - total
func retrieveInputs(accs []string, total *big.Int) (inputs []plasma.Position, change *big.Int, err error) {
	ctx := context.NewCLIContext().WithCodec(codec.New()).WithTrustNode(true)
	change = total
	// must specifiy inputs if using two accounts
	if len(accs) > 1 {
		return inputs, change, nil
	}

	addr, err := clistore.GetAccount(accs[0])
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
			return nil, nil, err
		}

		exitted, err := eth.HasTxExitted(utxo0.Position)
		if err != nil {
			return nil, nil, fmt.Errorf("Must connect full eth node or specify inputs using flags. Error encountered: %s", err)
		}
		if exitted {
			continue
		}
		// check if first utxo satisfies transfer amount
		if utxo0.Output.Amount.Cmp(total) == 0 {
			inputs = append(inputs, utxo0.Position)
			return inputs, big.NewInt(0), nil
		}
		for k, inner := range res {
			// do not pair an input with itself
			if i == k {
				continue
			}

			exitted, err := eth.HasTxExitted(utxo0.Position)
			if err != nil {
				return nil, nil, fmt.Errorf("Must connect full eth node or specify inputs using flags. Error encountered: %s", err)
			}
			if exitted {
				continue
			}

			if err := rlp.DecodeBytes(inner.Value, &utxo1); err != nil {
				return nil, nil, err
			}

			// check if only utxo1 satisfies transfer amount
			if utxo0.Output.Amount.Cmp(total) == 0 {
				inputs = append(inputs, utxo0.Position)
				return inputs, big.NewInt(0), nil
			}
			// check for exact match
			sum := new(big.Int).Add(utxo0.Output.Amount, utxo1.Output.Amount)
			if sum.Cmp(total) == 0 {
				inputs = append(inputs, utxo0.Position)
				inputs = append(inputs, utxo1.Position)
				return inputs, big.NewInt(0), nil
			}

			diff := new(big.Int).Sub(sum, total)
			if diff.Sign() == 1 && diff.Cmp(optimalChange) == -1 {
				optimalChange = diff
				position0 = utxo0.Position
				position1 = utxo1.Position
			}

		}
	}

	// check if a pairing was found
	if optimalChange.Cmp(total) == 0 {
		return inputs, optimalChange, nil
	}

	inputs = append(inputs, position0)
	inputs = append(inputs, position1)
	return inputs, optimalChange, nil
}
