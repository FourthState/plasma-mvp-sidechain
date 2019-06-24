package config

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AddPersistentTMFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(client.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	cmd.PersistentFlags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	cmd.PersistentFlags().String(client.FlagChainID, "", "id of the chain. Required if --trust-node=false")
	viper.BindPFlag(client.FlagTrustNode, cmd.Flags().Lookup(client.FlagTrustNode))
	viper.BindPFlag(client.FlagChainID, cmd.Flags().Lookup(client.FlagChainID))
	viper.BindPFlag(client.FlagNode, cmd.Flags().Lookup(client.FlagNode))
}
