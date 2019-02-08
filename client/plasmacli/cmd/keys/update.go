package keys

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"strings"
)

var updateCmd = &cobra.Command{
	Use:   "update <address>",
	Short: "Update account passphrase",
	Long:  `Update local encrypted private keys to be encrypted with the new passphrase`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		addrStr := strings.TrimSpace(args[0])
		if !ethcmn.IsHexAddress(addrStr) {
			return fmt.Errorf("invalid address provided, please use hex format")
		}
		addr := ethcmn.HexToAddress(addrStr)
		if err := keystore.Update(addr); err != nil {
			return err
		}

		fmt.Printf("Account %v passphrase has been updated.\n", addr)
		return nil
	},
}
