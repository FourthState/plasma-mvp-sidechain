package keys

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/store"
	"github.com/spf13/cobra"
	"os"
)

// ExportCmd returns the keys export command
func ExportCmd() *cobra.Command {
	return exportcmd
}

var exportcmd = &cobra.Command{
	Use:   "export <name> <location>",
	Short: "Export a private key",
	Long: `Exports a private key to a specified location (must be absolute path).

Usage:
	plasmacli export <name> <location>

The account is saved in encrypted format, you are prompted for a passphrase.
You must remember this passphrase to unlock your account to export.
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		location := args[1]

		fd, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return fmt.Errorf("error opening file: %s", err)
		}
		defer fd.Close()

		accountjson, err := store.Export(name)

		if numwritten, err := fd.Write(accountjson); err != nil {
			return fmt.Errorf("error writing to file: %s, bytes written: %s, total bytes to write: %s", err, numwritten, len(accountjson))
		}

		fmt.Printf("Successfully exported %s to %s", name, location)
		return nil
	},
}

