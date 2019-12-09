package config

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AddPersistentTMFlags adds the tendermint flags as persistent flags
func AddPersistentTMFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(client.FlagNode, "", "<host>:<port> to tendermint rpc interface for this chain")
	cmd.PersistentFlags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	cmd.PersistentFlags().String(client.FlagChainID, "", "id of the chain. Required if --trust-node=false")
	viper.BindPFlag(client.FlagNode, cmd.PersistentFlags().Lookup(client.FlagNode))
	viper.BindPFlag(client.FlagTrustNode, cmd.PersistentFlags().Lookup(client.FlagTrustNode))
	viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID))
}
