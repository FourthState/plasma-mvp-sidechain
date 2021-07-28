package msgs

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// TxDecoder attempts to RLP decode the transaction bytes into a SpendMsg first
// then to a IncludeDepositMsg then to a ConfirmSigMsg otherwise returns an error.
func TxDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var spendMsg SpendMsg
	if err := rlp.DecodeBytes(txBytes, &spendMsg); err != nil {
		var depositMsg IncludeDepositMsg
		if err2 := rlp.DecodeBytes(txBytes, &depositMsg); err2 != nil {
			var confirmSigMsg ConfirmSigMsg
			if err3 := rlp.DecodeBytes(txBytes, &confirmSigMsg); err3 != nil {
				return nil, sdk.ErrTxDecode(fmt.Sprintf("Decode to SpendMsg error: %s , Decode to DepositMsg error: %s, Decode to ConfirmSigMsg error: %s",
					err.Error(), err2.Error(), err3.Error()))
			}
			return confirmSigMsg, nil
		}
		return depositMsg, nil
	}
	return spendMsg, nil
}
