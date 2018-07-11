package cmd

import (
	"errors"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		addrStr := args[0]
		addr, err := client.StrToAddress(addrStr)
		if err != nil {
			return err
		}

		buf := client.BufferStdin()
		pass, err := client.GetPassword("DANGER - enter passphrase to permanetly delete key:", buf)
		if err != nil {
			return err
		}

		dir := viper.GetString(FlagHomeDir)
		ks := client.GetKeyStore(dir)
		if err != nil {
			return err
		}

		if !ks.HasAddress(addr) {
			return errors.New("the account trying to be deleted does not exist")
		}

		acc := accounts.Account{
			Address: addr,
		}

		err = ks.Delete(acc, pass)
		if err != nil {
			return err
		}
		fmt.Println("Account deleted forever")
		return nil
	},
}
