package query

import (
	"github.com/FourthState/plasma-mvp-sidechain/eth"
	"github.com/spf13/cobra"
)

var plasmaContract eth.Plasma

var (
	// flags
	addrF     = "address"
	allF      = "all"
	accountF  = "account"
	depositsF = "deposits"
	indexF    = "index"
	limitF    = "limit"
	positionF = "position"
)

func QueryCmd(plasma *eth.Plasma) *cobra.Command {
	plasmaContract = *plasma
	queryCmd.AddCommand(
		BalanceCmd(),
		BlockCmd(),
		DepositCmd(),
		ExitsCmd(),
		RootchainCmd(),
	)

	return queryCmd
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query for rootchain related information",
}
