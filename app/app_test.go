package app

import (
	"os"
	"testing"
	"fmt"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
	//types "plasma-mvp-sidechain/types" 
)

func newChildChain() *ChildChain {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "sdk/app")
	db := dbm.NewMemDB()
	return NewChildChain(logger, db)
}

func TestSpendMsg(t *testing.T) {
	//cc := newChildChain()
	fmt.Println("Testing has commenced")
	// Construct a SpendMsg
	//var msg = types.SpendMsg{}
}