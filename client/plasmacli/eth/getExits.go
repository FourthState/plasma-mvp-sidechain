package eth

import (
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/keystore"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"strconv"
)

func init() {
	queryCmd.AddCommand(exitsCmd)
	exitsCmd.Flags().Bool(allF, false, "all pending exits will be displayed")
	exitsCmd.Flags().String(limitF, "1", "amount of exits to display")
	exitsCmd.Flags().String(indexF, "", "index to begin displaying exits from")
	exitsCmd.Flags().StringP(accountF, "a", "", "display exits for given account")
	exitsCmd.Flags().BoolP(depositsF, "D", false, "display deposit exits")
	exitsCmd.Flags().StringP(positionF, "p", "", "display exit status for specified position")
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

		// check if position specified
		if viper.GetString(positionF) != "" {
			pos, err := plasma.FromPositionString(viper.GetString(positionF))
			if err != nil {
				return fmt.Errorf("failed to parse position: { %s }", err)
			}
			return displayExit(pos.Priority(), addr, pos.IsDeposit())
		}

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

		if err := displayExits(curr, lim, addr, viper.GetBool(depositsF)); err != nil {
			return fmt.Errorf("failure occured while querying exits: { %s }", err)
		}

		return nil
	},
}

func displayExits(curr, lim int64, addr ethcmn.Address, deposits bool) (err error) {
	for lim > 0 {
		var key *big.Int

		if deposits {
			key, err = rc.contract.DepositExitQueue(nil, big.NewInt(curr))
		} else {
			key, err = rc.contract.TxExitQueue(nil, big.NewInt(curr))
		}
		if err != nil {
			return err
		}

		// Get right 128 bits for position mapping
		key = new(big.Int).SetBytes(key.Bytes()[16:])

		if err := displayExit(key, addr, deposits); err != nil {
			return err
		}

		curr++
		lim--
	}
	return nil
}

// display a single exit given the position key in big.Int format
func displayExit(key *big.Int, addr ethcmn.Address, deposits bool) (err error) {
	var exit struct {
		Amount       *big.Int
		CommittedFee *big.Int
		CreatedAt    *big.Int
		Owner        ethcmn.Address
		State        uint8
	}

	if deposits {
		exit, err = rc.contract.DepositExits(nil, key)
	} else {
		exit, err = rc.contract.TxExits(nil, key)
	}

	if err != nil {
		return err
	}

	if !utils.IsZeroAddress(addr) && exit.Owner != addr {
		return nil
	}

	state := parseState(exit.State)
	fmt.Printf("Owner: 0x%x\nAmount: %d\nState: %s\nCommitted Fee: %d\nCreated: %v\n",
		exit.Owner, exit.Amount, state, exit.CommittedFee, exit.CreatedAt)
	return nil
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
