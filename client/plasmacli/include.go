package main

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/store"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/cosmos/cosmos-sdk/client/context"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"strings"
)

func init() {
	includeCmd.Flags().Int64P(replayF, "r", 0, "Replay Nonce that can be incremented to allow for resubmissions of include deposit messages")
	includeCmd.Flags().String(addressF, "", "address represented as hex string")
	includeCmd.Flags().Bool(asyncF, false, "wait for transaction commitment synchronously")
}

var includeCmd = &cobra.Command{
	Use:   "include-deposit <nonce> <account_name>",
	Short: "Include a deposit from <account_name> with given nonce",
	Long: `Example usage:
	plasmacli include-deposit <nonce> <account_name>
	plasmacli include-deposit <nonce> --address <address>
	plasmacli include-deposit <nonce> --address <address> -r 3`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx := context.NewCLIContext()

		// validate addresses
		var address ethcmn.Address
		var err error
		if len(args) == 2 {
			address, err = store.GetAccount(args[1])
			if err != nil {
				return fmt.Errorf("Could not retrieve account: %s", args[1])
			}
		} else {
			addrToken := viper.GetString(addressF)
			addrToken = strings.TrimSpace(addrToken)
			if !ethcmn.IsHexAddress(addrToken) {
				return fmt.Errorf("invalid address provided. please use hex format")
			}
			address := ethcmn.HexToAddress(addrToken)
			if utils.IsZeroAddress(address) {
				return fmt.Errorf("cannot include deposit from the zero address")
			}
		}

		nonce, ok := new(big.Int).SetString(strings.TrimSpace(args[0]), 10)
		if !ok {
			return fmt.Errorf("could not parse deposit nonce. Please resubmit with nonce in base 10")
		}

		replay := viper.GetInt(replayF)

		msg := msgs.IncludeDepositMsg{
			DepositNonce: nonce,
			Owner:        address,
			ReplayNonce:  uint64(replay),
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
			fmt.Printf("Committed at block %d. Hash %s\n", res.Height, res.TxHash)
		}

		return nil
	},
}
