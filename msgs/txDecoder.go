package msgs

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// NOTE: This is fucking stupid

func TxDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var spendMsg SpendMsg
	if err := rlp.DecodeBytes(txBytes, &spendMsg); err != nil {
		var depositMsg IncludeDepositMsg
		if err2 := rlp.DecodeBytes(txBytes, &depositMsg); err2 != nil {
			var initiatePresenceClaimMessage InitiatePresenceClaimMsg
			if err3 := rlp.DecodeBytes(txBytes, &initiatePresenceClaimMessage); err3 != nil {
				var postLogsMsg PostLogsMsg
				if err4 := rlp.DecodeBytes(txBytes, &postLogsMsg); err4 != nil {
					var createZoneMsg CreateZoneMsg
					if err5 := rlp.DecodeBytes(txBytes, &createZoneMsg); err5 != nil {
						return nil, sdk.ErrTxDecode(fmt.Sprintf("Decode to SpendMsg error: { %s } Decode to DepositMsg error: { %s } Decode to InitiatePresenceClaimMsg error: { %s }  Decode to PostLogsMsg error: { %s } Decode to CreateZoneMsg error: { %s }",
							err.Error(), err2.Error(), err3.Error(), err4.Error(), err5.Error()))

					}
					return createZoneMsg, nil

				}
				return postLogsMsg, nil
			}
			return initiatePresenceClaimMessage, nil
		}
		return depositMsg, nil
	}

	return spendMsg, nil
}
