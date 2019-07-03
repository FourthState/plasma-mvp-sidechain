package main

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/cmd/plasmacli/subcmd"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"os"
)

const (
	FlagHome = tmcli.HomeFlag
)

func main() {
	rootCmd := subcmd.RootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
