package eth

import (
	"fmt"
	"math/big"
	"strconv"

	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	queryCmd.AddCommand(exitsCmd)
	exitsCmd.Flags().Bool(allF, false, "all pending exits will be displayed")
	exitsCmd.Flags().String(limitF, "1", "amount of exits to display")
	exitsCmd.Flags().String(indexF, "", "index to begin displaying exits from")
	exitsCmd.Flags().StringP(accountF, "a", "", "display exits for given account")
	exitsCmd.Flags().BoolP(depositsF, "D", false, "display deposit exits")
	viper.BindPFlags(exitsCmd.Flags())
}

var exitsCmd = &cobra.Command{
	Use:   "exits",
	Short: "Display pending exits",
	Long: `Display pending rootchain exits. Queries the rootchain exit queue.
Use the deposit flag to display deposit exits.

Usage:
	plasmacli eth query exits -a <account>
	plasmacli eth query exits --deposits
	plasmacli eth query exits --all
	plasmacli eth query exits --index <number> --limit <number>`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var curr int64
		var addr ethcmn.Address
		acc := viper.GetString(accountF)

		if acc != "" {
			addr, err = ks.Get(acc)
			if err != nil {
				return fmt.Errorf("failed to retrieve account: { %s }", err)
			}
		}

		// parse queue length
		var len *big.Int
		if viper.GetBool(depositsF) {
			len, err = rc.contract.DepositQueueLength(nil)
		} else {
			len, err = rc.contract.TxQueueLength(nil)
		}

		if err != nil {
			return err
		}

		lim, err := strconv.ParseInt(viper.GetString(limitF), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse limit: { %s }", err)
		}

		// adjust passed in limit to avoid error
		// when interacting with rootchain
		if lim > len.Int64() {
			lim = len.Int64()
		}

		if viper.GetBool(allF) { // print all exits
			lim = len.Int64()
		} else { // use index/limit
			index := viper.GetString(indexF)
			if index == "" {
				return fmt.Errorf("please specify one of the following flags: --all, --index, --account")
			}
			curr, err = strconv.ParseInt(viper.GetString(indexF), 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse index: { %s }", err)
			}
		}

		if viper.GetBool(depositsF) {
			displayDepositExits(curr, lim, addr)
		} else {
			displayTxExits(curr, lim, addr)
		}

		return nil
	},
}

func displayDepositExits(curr, lim int64, addr ethcmn.Address) {
	for lim > 0 {
		key, err := rc.contract.DepositExitQueue(nil, big.NewInt(curr))
		if err != nil {
			return
		}

		// Get right 128 bits for position mapping
		key = new(big.Int).SetBytes(key.Bytes()[16:])

		exit, err := rc.contract.DepositExits(nil, key)
		if err != nil {
			return
		}
		if !utils.IsZeroAddress(addr) && exit.Owner != addr {
			continue
		}

		curr++
		lim--
		state := parseState(exit.State)
		fmt.Printf("Owner: 0x%x\nAmount: %d\nState: %s\nCommitted Fee: %d\nCreated: %v\n",
			exit.Owner, exit.Amount, state, exit.CommittedFee, exit.CreatedAt)
	}
}

func displayTxExits(curr, lim int64, addr ethcmn.Address) {
	for lim > 0 {
		key, err := rc.contract.TxExitQueue(nil, big.NewInt(curr))
		if err != nil {
			return
		}

		// Get right 128 bits for position mapping
		key = new(big.Int).SetBytes(key.Bytes()[16:])

		exit, err := rc.contract.TxExits(nil, key)
		if err != nil {
			return
		}

		if !utils.IsZeroAddress(addr) && exit.Owner != addr {
			continue
		}

		curr++
		lim--
		state := parseState(exit.State)
		fmt.Printf("Owner: 0x%x\nAmount: %d\nState: %s\nCommitted Fee: %d\nCreated: %v\n",
			exit.Owner, exit.Amount, state, exit.CommittedFee, exit.CreatedAt)
	}
}

func parseState(exit uint8) (state string) {
	switch exit {
	case 0:
		state = "Nonexistent"
	case 1:
		state = "Pending"
	case 2:
		state = "Challenged"
	case 3:
		state = "Finalized"
	}
	return state
}
