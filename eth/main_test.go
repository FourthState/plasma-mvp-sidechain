package eth

import (
	"testing"
)

const clientAddr = "http://127.0.0.1:8545"

func TestConnection(t *testing.T) {
	t.Logf("Connecting to remote client: %s", clientAddr)
	client, err := InitEthConn(clientAddr)
	if err != nil {
		t.Error("Connection Error:", err)
	}

	_, err = client.accounts()
	if err != nil {
		t.Error("Error Retrieving Accounts:", err)
	}
}
