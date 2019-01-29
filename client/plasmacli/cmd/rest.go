package cmd

import (
	"github.com/FourthState/plasma-mvp-sidechain/rest"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/codec"
)

func init() {
	rootCmd.AddCommand(lcd.ServeCommand(codec.New(), registerRoutes))
}

func registerRoutes(rs *lcd.RestServer) {
	rest.RegisterRoutes(rs.CliCtx, rs.Mux)
}
