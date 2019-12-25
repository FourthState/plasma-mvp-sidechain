package subcmd

import (
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var versionCmd = version.VersionCmd

// VersionCmd -
func VersionCmd() *cobra.Command {
	viper.BindPFlags(versionCmd.LocalFlags())
	return versionCmd
}
