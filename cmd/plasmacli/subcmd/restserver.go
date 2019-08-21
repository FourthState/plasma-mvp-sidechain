package subcmd

import (
	"github.com/FourthState/plasma-mvp-sidechain/app"
	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	sdkCli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RestServerCmd() *cobra.Command {
	serverCmd.Flags().String(sdkCli.FlagListenAddr, "tcp://localhost:1317", "The address for the server to listen on")
	serverCmd.Flags().Bool(sdkCli.FlagTLS, false, "Enable SSL/TLS layer")
	serverCmd.Flags().String(sdkCli.FlagSSLHosts, "", "Comma-separated hostnames and IPs to generate a certificate for")
	serverCmd.Flags().String(sdkCli.FlagSSLCertFile, "", "Path to a SSL certificate file. If not supplied, a self-signed certificate will be generated.")
	serverCmd.Flags().String(sdkCli.FlagSSLKeyFile, "", "Path to a key file; ignored if a certificate file is not supplied.")
	serverCmd.Flags().String(sdkCli.FlagCORS, "", "Set the domains that can make CORS requests (* for all)")
	serverCmd.Flags().Int(sdkCli.FlagMaxOpenConnections, 1000, "The number of maximum open connections")

	config.AddPersistentTMFlags(serverCmd)
	return serverCmd
}

var serverCmd = &cobra.Command{
	Use:   "rest-server",
	Short: "Start LCD (light-client daemon), a local REST server",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		rs := lcd.NewRestServer(app.MakeCodec())
		client.RegisterRoutes(rs.CliCtx, rs.Mux)

		// Start the rest server and return error if one exists
		err := rs.Start(
			viper.GetString(sdkCli.FlagListenAddr),
			viper.GetString(sdkCli.FlagSSLHosts),
			viper.GetString(sdkCli.FlagSSLCertFile),
			viper.GetString(sdkCli.FlagSSLKeyFile),
			viper.GetInt(sdkCli.FlagMaxOpenConnections),
			viper.GetBool(sdkCli.FlagTLS))

		return err
	},
}
