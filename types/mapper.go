package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UTXOMapper stores and retrieves UTXO's from stores
// retrieved from the context.
type UTXOMapper interface {
	GetUTXO(ctx sdk.Context, position Position) UTXO
	AddUTXO(ctx sdk.Context, utxo UTXO)
	DeleteUTXO(ctx sdk.Context, position Position)
}
