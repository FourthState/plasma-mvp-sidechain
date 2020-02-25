package keys

import (
	"github.com/spf13/cobra"
)

// flags
const (
	nameF = "name"
	fileF = "file"
	encryptF = "encrypted"
)

// RootCmd returns the keys command
func RootCmd() *cobra.Command {
	keysCmd.AddCommand(
		AddCmd(),
		DeleteCmd(),
		ListCmd(),
		UpdateCmd(),
		ImportCmd(),
		ExportCmd(),
	)

	return keysCmd
}

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage local private keys",
	Long:  `Keys allows you to manage your local keystore.`,
}
