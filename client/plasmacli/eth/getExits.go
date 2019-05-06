package eth

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/client/store"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
	"strconv"
	"time"
)

func init() {
	queryCmd.AddCommand(getExitsCmd)
	getExitsCmd.Flags().StringP(accountF, "a", "", "display exits for given account")
	getExitsCmd.Flags().Bool(allF, false, "all pending exits will be displayed")
	getExitsCmd.Flags().BoolP(depositsF, "D", false, "display deposit exits")
	getExitsCmd.Flags().String(indexF, "0", "index to begin displaying exits from")
	getExitsCmd.Flags().String(limitF, "1", "amount of exits to display")
	getExitsCmd.Flags().StringP(positionF, "p", "", "display exit status for specified position")
}

var getExitsCmd = &cobra.Command{
	Use:   "exit",
	Short: "Display pending exits",
	Long: `Display pending rootchain exits. Queries the rootchain exit queue.
Use the deposit flag to display deposit exits.

Usage:
	plasmacli eth query exit -a <account>
	plasmacli eth query exit --deposits
	plasmacli eth query exit --all
	plasmacli eth query exit --index <number> --limit <number>
	plasmacli eth query exit --position <position>`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		viper.BindPFlags(cmd.Flags())

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
			addr, err = store.GetAccount(acc)
			if err != nil {
				return fmt.Errorf("failed to retrieve account: { %s }", err)
			}
		}

		// parse queue length
		var queueLength *big.Int
		if viper.GetBool(depositsF) {
			queueLength, err = rc.contract.DepositQueueLength(nil)
		} else {
			queueLength, err = rc.contract.TxQueueLength(nil)
		}

		if err != nil {
			return fmt.Errorf("failed to retrieve exit queue length: { %s }", err)
		}

		lim, err := strconv.ParseInt(viper.GetString(limitF), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse limit: { %s }", err)
		}

		// adjust passed in limit to avoid error
		// when interacting with rootchain
		if lim > queueLength.Int64() {
			lim = queueLength.Int64()
		}

		if viper.GetBool(allF) { // print all exits
			lim = queueLength.Int64()
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
		EthBlockNum  *big.Int
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
	fmt.Printf("Owner: 0x%x\nAmount: %d\nState: %s\nCommitted Fee: %d\nCreated: %v\nEthBlockNum: %v\n\n",
		exit.Owner, exit.Amount, state, exit.CommittedFee, time.Unix(exit.CreatedAt.Int64(), 0), exit.EthBlockNum)
	if state == "Pending" {
		timeLeft := time.Until(time.Unix(exit.CreatedAt.Int64(), 0).Add(time.Hour * oneWeek))
		if timeLeft > 0 {
			fmt.Printf("Exit will be finalized in about: %v hours\n\n", timeLeft.Hours())
		} else {
			fmt.Printf("Exit is ready to be finalized!\n\n")
		}
	}

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
