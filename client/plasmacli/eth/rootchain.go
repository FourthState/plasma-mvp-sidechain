package eth

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	queryCmd.AddCommand(rootchainCmd)
}

var rootchainCmd = &cobra.Command{
	Use:   "rootchain",
	Short: "Display rootchain contract information",
	Long: `Display last committed block, total contract balance, total withdraw balance, minimum exit bond, and operator address.
Total contract balance does not include total withdraw balance. The total withdraw balance are exits that have been finalized, but not transfered yet.`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, agrs []string) error {
		lastCommittedBlock, err := rc.session.LastCommittedBlock()
		if err != nil {
			return err
		}

		totalBalance, err := rc.session.PlasmaChainBalance()
		if err != nil {
			return err
		}

		withdrawBalance, err := rc.session.TotalWithdrawBalance()
		if err != nil {
			return err
		}

		minExitBond, err := rc.session.MinExitBond()
		if err != nil {
			return err
		}

		operator, err := rc.session.Operator()
		if err != nil {
			return err
		}
		fmt.Printf("Last Committed Block: %d\nContract Balance: %d\nWithdraw Balance: %d\nMinimum Exit Bond: %d\nOperator: 0x%x\n",
			lastCommittedBlock, totalBalance, withdrawBalance, minExitBond, operator)
		return nil
	},
}
