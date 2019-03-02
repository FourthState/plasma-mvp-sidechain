package main

import (
	"fmt"
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
	rootCmd.AddCommand(includeCmd)
	includeCmd.Flags().Int64P(flagReplay, "r", 0, "Replay Nonce that can be incremented to allow for resubmissions of include deposit messages")
}

var includeCmd = &cobra.Command{
	Use:   "include-deposit <nonce> <account>",
	Short: "Include a deposit from <account> with given nonce",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext()

		// validate addresses
		addrToken := strings.TrimSpace(args[1])

		if !ethcmn.IsHexAddress(addrToken) {
			return fmt.Errorf("invalid address provided. please use hex format")
		}
		addr := ethcmn.HexToAddress(addrToken)
		if utils.IsZeroAddress(addr) {
			return fmt.Errorf("cannot include deposit from the zero address")
		}

		nonce, ok := new(big.Int).SetString(strings.TrimSpace(args[0]), 10)
		if !ok {
			return fmt.Errorf("could not parse deposit nonce. Please resubmit with nonce in base 10")
		}

		replay := viper.GetInt(flagReplay)

		msg := msgs.IncludeDepositMsg{
			DepositNonce: nonce,
			Owner:        addr,
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
