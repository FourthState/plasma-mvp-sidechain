package msgs

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// TxDecoder attempts to RLP decode the transaction bytes into a SpendMsg or
// IncludeDepositMsg otherwise returns an error.
func TxDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var spendMsg SpendMsg
	if err := rlp.DecodeBytes(txBytes, &spendMsg); err != nil {
		var depositMsg IncludeDepositMsg
		if err2 := rlp.DecodeBytes(txBytes, &depositMsg); err2 != nil {
			return nil, sdk.ErrTxDecode(fmt.Sprintf("Decode to SpendMsg error: { %s } Decode to DepositMsg error: { %s }",
				err.Error(), err2.Error()))
		}
		return depositMsg, nil
	}

	return spendMsg, nil
}
