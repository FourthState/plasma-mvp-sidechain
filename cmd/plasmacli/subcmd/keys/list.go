package keys

import (
	"errors"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

// ListCmd returns the keys list command
func ListCmd() *cobra.Command {
	return listCmd
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long:  "Return a list of all account addresses stored by the local keystore",
	RunE: func(cmd *cobra.Command, args []string) error {
		iter, db := store.AccountIterator()
		if iter == nil || db == nil {
			return errors.New("unexpected error encountered when opening account data")
		}
		defer db.Close()

		w := new(tabwriter.Writer)
		// Sets tab width to 8 characters
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(w, "NAME:\tADDRESS:\t")
		for iter.Next() {
			fmt.Fprintln(w, fmt.Sprintf("%s\t0x%x\t", iter.Key(), ethcmn.BytesToAddress(iter.Value())))
		}
		fmt.Fprintln(w)
		w.Flush()
		iter.Release()

		return nil
	},
}
