package cmd

import (
	"errors"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func init() {
	rootCmd.AddCommand(deleteKeyCmd)
}

var deleteKeyCmd = &cobra.Command{
	Use:   "delete <address>",
	Short: "Delete the given address",
	Long: `Deletes the account from the keystore matching the address provided, if the passphrase
			is correct.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		addrStr := strings.TrimSpace(args[0])
		if !common.IsHexAddress(addrStr) {
			return fmt.Errorf("Invalid address provided. please use hex format")
		}

		// delete from the keystore
		dir := viper.GetString(client.FlagHomeDir)
		keystore.InitKeystore(dir)
		if err := keystore.Delete(common.HexToAddress(addrStr)); err != nil {
			return err
		}

		fmt.Println("Account deleted.")
		return nil
	},
}
