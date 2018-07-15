package cmd

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/ethereum/go-ethereum/accounts"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(showKeysCmd)
}

var showKeysCmd = &cobra.Command{
	Use:   "find <address>",
	Short: "Find the account for the given address",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		addrStr := args[0]
		addr, err := client.StrToAddress(addrStr)
		if err != nil {
			return err
		}
		acct := accounts.Account{
			Address: addr,
		}

		dir := viper.GetString(FlagHomeDir)
		ks := client.GetKeyStore(dir)
		acc, err := ks.Find(acct)
		if err != nil {
			return err
		}
		fmt.Println("Your account has been found!")
		fmt.Printf("Account Address: %s\nAccount Location:%s\n", acc.Address.Hex(), acc.URL.String())
		return nil
	},
}
