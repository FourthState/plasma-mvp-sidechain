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
	Long: `Exports a private key to a specified location.

Usage:
	plasmacli export <name> <location>

The account is saved in encrypted format, you are prompted for a passphrase.
You must remember this passphrase to unlock your account to export.
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		location := args[1]
		cmd.SilenceUsage = true

		accountjson, err := store.Export(name)
		if err != nil {
			return fmt.Errorf("error exporting key: %s", err.Error())
		}

		fd, err := os.OpenFile(location, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return fmt.Errorf("error opening file: %s", err)
		}
		defer fd.Close()

		if numwritten, err := fd.Write(accountjson); err != nil {
			return fmt.Errorf("error writing to file: %s, bytes written: %d, total bytes to write: %d", err, numwritten, len(accountjson))
		}

		fmt.Printf("Successfully exported %s to %s", name, location)
		return nil
	},
}

