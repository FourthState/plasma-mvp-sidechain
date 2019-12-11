package plasma

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

func TestDepositSerialization(t *testing.T) {
	one := big.NewInt(1)
	deposit := &Deposit{common.HexToAddress("0"), one, one}

	bytes, err := rlp.EncodeToBytes(deposit)
	require.NoError(t, err, "Error serializing deposit")

	recoveredDeposit := &Deposit{}
	err = rlp.DecodeBytes(bytes, recoveredDeposit)
	require.NoError(t, err, "Error deserializing deposit")

	require.True(t, reflect.DeepEqual(deposit, recoveredDeposit), "serialized and deserialized deposits are not deeply equal")
}
