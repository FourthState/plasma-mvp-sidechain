package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	ks "github.com/FourthState/plasma-mvp-sidechain/client/store"
	"github.com/FourthState/plasma-mvp-sidechain/query"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func init() {
	queryCmd.AddCommand(balanceCmd)
}

var balanceCmd = &cobra.Command{
	Use:   "balance <name>",
	Short: "Query account balance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCLIContext().WithCodec(codec.New()).WithTrustNode(true)
		name := args[0]

		addr, err := ks.GetAccount(name)
		if err != nil {
			return err
		}

		queryRoute := fmt.Sprintf("custom/utxo/balance/%s", addr.Hex())
		data, err := ctx.Query(queryRoute, nil)
		if err != nil {
			return err
		}

		var resp query.BalanceResp
		if err := json.Unmarshal(data, &resp); err != nil {
			return err
		} else if !bytes.Equal(resp.Address[:], addr[:]) {
			return fmt.Errorf("Mismatch in Account and Response Address.\nAccount: 0x%x\n Response: 0x%x\n",
				addr, resp.Address)
		}

		fmt.Printf("Address: %0x\n", resp.Address)
		fmt.Printf("Total: %s\n", resp.Total)
		return nil
	},
}
