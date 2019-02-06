package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func TxDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var spendMsg SpendMsg
	if err := rlp.DecodeBytes(txBytes, &spendMsg); err != nil {
		return SpendMsg{}, sdk.ErrTxDecode(err.Error())
	}

	return spendMsg, nil
}
