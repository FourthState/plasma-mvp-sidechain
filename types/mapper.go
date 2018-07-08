package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// UTXOMapper stores and retrieves UTXO's from stores
// retrieved from the context.
type UTXOMapper interface {
	GetUTXO(ctx sdk.Context, addr common.Address, position Position) UTXO
	GetAllUTXOsForAddress(ctx sdk.Context, addr common.Address) []UTXO
	AddUTXO(ctx sdk.Context, utxo UTXO)
	DeleteUTXO(ctx sdk.Context, addr common.Address, position Position)
}
