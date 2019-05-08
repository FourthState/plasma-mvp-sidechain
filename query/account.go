package query

import (
	"errors"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func queryBalance(ctx sdk.Context, utxoStore store.UTXOStore, addr common.Address) (*big.Int, error) {
	iter := sdk.KVStorePrefixIterator(utxoStore.KVStore(ctx), addr.Bytes())
	total := big.NewInt(0)
	for ; iter.Valid(); iter.Next() {
		utxo, ok := utxoStore.GetUTXOWithKey(ctx, iter.Key())
		if !ok {
			return nil, errors.New("invalid key retrieved")
		}

		if !utxo.Spent {
			total = total.Add(total, utxo.Output.Amount)
		}
	}

	return total, nil
}

func queryInfo(ctx sdk.Context, utxoStore store.UTXOStore, addr common.Address) ([]store.UTXO, error) {
	var utxos []store.UTXO
	iter := sdk.KVStorePrefixIterator(utxoStore.KVStore(ctx), addr.Bytes())
	for ; iter.Valid(); iter.Next() {
		utxo, ok := utxoStore.GetUTXOWithKey(ctx, iter.Key())
		if !ok {
			return nil, errors.New("invalid key retrieved")
		}

		utxos = append(utxos, utxo)
	}

	return utxos, nil
}
