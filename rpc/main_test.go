package rpc

import (
	"context"
	"testing"
	"time"
)

const clientAddr = "http://127.0.0.1:8545"

func TestConnection(t *testing.T) {
	t.Logf("Connecting to remote client: %s", clientAddr)
	client, err := ethclient.Dial(clientAddr)
	if err != nil {
		t.Error("Connection failed,", err)
	}

}
