package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/ethereum/go-ethereum/accounts"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

const (
	flagAddr = "addr"
)

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().String(flagAddr, "", "Address to sign with")
	viper.BindPFlags(signCmd.Flags())
}

var signCmd = &cobra.Command{
	Use:   "sign <position>",
	Short: "Sign positions to create confirmation signatures, format: 0.0.0.0::0.0.0.0",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// get position
		posStr := args[0]
		pos, err := client.ParsePositions(posStr)
		if err != nil {
			return err
		}

		ctx := context.NewClientContextFromViper()

		dir := viper.GetString(FlagHomeDir)

		ks := client.GetKeyStore(dir)
		// get address to sign with
		addrStr := viper.GetString(flagAddr)
		addr, err := client.StrToAddress(addrStr)
		if err != nil {
			return err
		}
		acc := accounts.Account{
			Address: addr,
		}
		// get account to sign with
		acct, err := ks.Find(acc)
		if err != nil {
			return err
		}

		// get passphrase
		passphrase, err := ctx.GetPassphraseFromStdin(addr)
		if err != nil {
			return err
		}

		// sign positions
		signBytes1 := pos[0].GetSignBytes()
		hash1 := ethcrypto.Keccak256(signBytes1)
		sig1, err := ks.SignHashWithPassphrase(acct, passphrase, hash1)

		var sig2 []byte
		if len(pos) > 1 {
			signBytes2 := pos[1].GetSignBytes()
			hash2 := ethcrypto.Keccak256(signBytes2)
			sig2, err = ks.SignHashWithPassphrase(acct, passphrase, hash2)
			if err != nil {
				return err
			}
		}
		fmt.Println()
		fmt.Println("Sig 1:")
		fmt.Printf("%X", sig1)
		fmt.Println()
		fmt.Println("Sig 2:")
		fmt.Printf("%X", sig2)
		fmt.Println()
		return nil
	},
}
