package msgs

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"math/big"
)

var (
	privKey, _ = crypto.GenerateKey()
	addr       = crypto.PubkeyToAddress(privKey.PublicKey)
)

func TestDepositMsgValidate(t *testing.T) {
	type depositCase struct {
		Nonce *big.Int
		Address      ethcmn.Address
	}

	invalidCases := []depositCase{
		{big.NewInt(-1), addr},
		{big.NewInt(0), addr},
		{big.NewInt(1), ethcmn.Address{}},
	}

	for i, c := range invalidCases {
		depositMsg := IncludeDepositMsg{
			DepositNonce: c.Nonce,
			Owner: c.Address,
		}
		require.NotNil(t, depositMsg.ValidateBasic(), fmt.Sprintf("Testcase %d failed", i))
	}
}

func TestDepositMsgSerialization(t *testing.T) {
	msg := IncludeDepositMsg{
		DepositNonce: big.NewInt(3),
		Owner: addr,
		ReplayNonce: 2,
	}

	bytes, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err, "serialization error")

	tx, err := TxDecoder(bytes)
	
	require.NoError(t, err, "deserialization error")

	require.True(t, reflect.DeepEqual(msg, tx), "serialized and deserialized msgs not equal")
}
