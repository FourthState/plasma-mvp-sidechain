package subcmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	return versionCmd
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the plasma client",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Plasma Client v0.3.0")
	},
}
