package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func TxDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var spendMsg SpendMsg
	if err := rlp.DecodeBytes(txBytes, &spendMsg); err != nil {
		var depositMsg IncludeDepositMsg
		if err2 := rlp.DecodeBytes(txBytes, &depositMsg); err2 != nil {
			return nil, sdk.ErrTxDecode(err.Error())
		}
		return depositMsg, nil
	}

	return spendMsg, nil
}
