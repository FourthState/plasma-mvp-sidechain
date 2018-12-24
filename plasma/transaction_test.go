package plasma

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

func TestTransactionSerialization(t *testing.T) {
	one := big.NewInt(1)
	zero := big.NewInt(0)

	// contstruct a transaction
	tx := &Transaction{}
	tx.Input0 = NewInput(one, 1, 1, one,
		common.HexToAddress("1"), [65]byte{}, [][65]byte{[65]byte{}})
	tx.Input0.Signature[1] = byte(1)
	tx.Input1 = NewInput(zero, 0, 0, zero,
		common.HexToAddress("0"), [65]byte{}, [][65]byte{[65]byte{}})
	tx.Output0 = NewOutput(common.HexToAddress("1"), one)
	tx.Output1 = NewOutput(common.HexToAddress("0"), zero)
	tx.Fee = big.NewInt(1)

	bytes, err := rlp.EncodeToBytes(tx)
	require.NoError(t, err, "Error serializing transaction")

	recoveredTx := &Transaction{}
	err = rlp.DecodeBytes(bytes, recoveredTx)
	require.NoError(t, err, "Error deserializing transaction")

	require.True(t, reflect.DeepEqual(tx, recoveredTx), "serialized and deserialized objects not deeply equal")
}
