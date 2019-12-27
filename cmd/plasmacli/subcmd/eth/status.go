package eth

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/config"
	"github.com/FourthState/plasma-mvp-sidechain/eth"
	"github.com/spf13/cobra"
)
// StatusCmd returns the current state of the eth connection (syncing, crashed etc)
func StatusCmd() *cobra.Command {
	return statuscmd
}

var statuscmd = &cobra.Command {
	Use:   "status",
	Short: "check state of eth connection",
	Long: "returns current state of eth connection (syncing, crashed, etc)",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Get Conf
		conf, err := config.ParseConfigFromViper()
		if err != nil {
			return fmt.Errorf("error retrieving config: %s", err)
		}

		// Initialize Eth Connection
		client, err := eth.InitEthConn(conf.EthNodeURL)
		if err != nil {
			return fmt.Errorf("error initializing connection with client: %s", err)
		}

		// Check whether client synced
		synced, err := client.Synced()
		if err != nil {
			return fmt.Errorf("error checking synced status: %s", err)
		}
		if synced {
			num, err := client.LatestBlockNum()
			if err != nil {
				return fmt.Errorf("synced with eth node: %s \n error retrieving latest block height of eth endpoint: %s\n", conf.EthNodeURL , err)
			}
			fmt.Printf("synced with eth node: %s \nlatest block height of the eth endpoint: %d\n", conf.EthNodeURL, num)

		} else {
			fmt.Printf("could not sync with eth node: %s", conf.EthNodeURL)
		}
		return nil
	},
}
