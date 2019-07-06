package keys

import (
	"github.com/spf13/cobra"
)

// flags
const (
	nameF = "name"
	fileF = "file"
)

func KeysCmd() *cobra.Command {
	keysCmd.AddCommand(
		AddCmd(),
		DeleteCmd(),
		ImportCmd(),
		ListCmd(),
		UpdateCmd(),
	)

	return keysCmd
}

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage local private keys",
	Long:  `Keys allows you to manage your local keystore.`,
}
