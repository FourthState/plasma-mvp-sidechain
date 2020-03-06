package eth

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/subcmd/eth/query"
	"github.com/FourthState/plasma-mvp-sidechain/eth"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

const (
	// flags
	accountF  = "account"
	allF      = "all"
	depositsF = "deposits"
	feeF      = "fee"
	gasLimitF = "gas-limit"
	indexF    = "index"
	limitF    = "limit"
	ownerF    = "owner"
	positionF = "position"
	proofF    = "proof"
	sigsF     = "signatures"
	useNodeF  = "use-node"
	txBytesF  = "tx-bytes"

	minExitBond = 200000
)

var plasmaContract *eth.Plasma

// RootCmd returns root eth command
func RootCmd() *cobra.Command {
	ethCmd.AddCommand(
		ProveCmd(),
		ChallengeCmd(),
		ExitCmd(),
		FinalizeCmd(),
		DepositCmd(),
		WithdrawCmd(),
		client.LineBreak,

		query.RootCmd(),
	)

	return ethCmd
}

var ethCmd = &cobra.Command{
	Use:   "eth",
	Short: "Interact with the plasma smart contract",
	Long: `Configurations for interacting with the rootchain contract can be specified in <dirpath>/plasma.toml.
An eth node instance needs to be running for this command to work.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("setting up eth connection")
		plasma, err := config.GetContractConn()
		plasmaContract = plasma
		return err
	},
}

// HasTxExited -
func HasTxExited(pos plasma.Position) (bool, error) {
	conn, err := config.GetContractConn()
	if err != nil {
		return false, err
	}

	return conn.HasTxExited(nil, pos)
}
