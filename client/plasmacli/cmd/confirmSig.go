package cmd

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/FourthState/plasma-mvp-sidechain/app"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/client/context"
	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

const (
//flagAddr = "addr"
)

func init() {
	rootCmd.AddCommand(signCmd)
	//signCmd.Flags().String(flagAddr, "", "Address to sign with")
	viper.BindPFlags(signCmd.Flags())
}

var signCmd = &cobra.Command{
	Use:   "sign <amount> <address>",
	Short: "Sign confirmation signatures for amount provided, if it exists",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.NewClientContextFromViper()
		cdc := app.MakeCodec()
		ctx = ctx.WithCodec(cdc)

		ethAddr := common.HexToAddress(args[1])

		res, err := ctx.QuerySubspace(ethAddr.Bytes(), ctx.UTXOStore)
		if err != nil {
			return err
		}

		var utxos []types.BaseUTXO
		for _, pair := range res {
			var utxo types.BaseUTXO
			err := ctx.Codec.UnmarshalBinaryBare(pair.Value, &utxo)
			if err != nil {
				return err
			}
			utxos = append(utxos, utxo)
		}

		dir := viper.GetString(FlagHomeDir)

		ks := client.GetKeyStore(dir)
		// get address to sign with
		//addrStr := viper.GetString(flagAddr)
		//addr, err := client.StrToAddress(addrStr)
		//if err != nil {
		//	return err
		//}
		acc := accounts.Account{
			Address: ethAddr,
		}
		// get account to sign with
		acct, err := ks.Find(acc)
		if err != nil {
			return err
		}

		amount, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return err
		}
		// get passphrase
		passphrase, err := ctx.GetPassphraseFromStdin(ethAddr)
		if err != nil {
			return err
		}

		index := -1
		for i, utxo := range utxos {
			if utxo.GetAmount() == amount {
				index = i
				break
			}
		}

		if index == -1 {
			fmt.Println("Sorry, the provided amount and address do not match to an existing utxo")
			return nil
		}

		blknumKey := make([]byte, binary.MaxVarintLen64)
		binary.PutUvarint(blknumKey, utxos[index].GetPosition().Get()[0].Uint64())

		blockhash, err := ctx.QuerySubspace(blknumKey, ctx.MetadataStore)
		if err != nil {
			return err
		}

		hash := ethcrypto.Keccak256(append(utxos[index].MsgHash, blockhash[0].Value...))
		sig, err := ks.SignHashWithPassphrase(acct, passphrase, hash)

		fmt.Printf("\nConfirmation Signature for utxo with\n position: %v \namount: %d\n")
		fmt.Printf("Signature:%v\n", sig)
		fmt.Printf("UTXO had %d inputs\n")
		return nil
	},
}
