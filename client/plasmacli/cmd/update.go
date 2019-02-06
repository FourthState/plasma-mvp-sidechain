package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(updateKeysCmd)
}

var updateKeysCmd = &cobra.Command{
	Use:   "update <address>",
	Short: "Change the password used to protect the account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		addrStr := strings.TrimSpace(args[0])
		if !common.IsHexAddress(addrStr) {
			fmt.Errorf("invalid address provided. please use hex format")
		}

		if err := keystore.Update(common.HexToAddress(addrStr)); err != nil {
			return err
		}

		fmt.Println("Password successfully updated!")
		return nil
	},
}
