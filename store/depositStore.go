package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

type Deposit struct {
	Deposit plasma.Deposit
	Spent   bool
	Spender []byte
}

/* Deposit Store */
type DepositStore struct {
	kvStore
}

func NewDepositStore(ctxKey sdk.StoreKey) DepositStore {
	return DepositStore{
		kvStore: NewKVStore(ctxKey),
	}
}

func (store DepositStore) GetDeposit(ctx sdk.Context, nonce *big.Int) (Deposit, bool) {
	data := store.Get(ctx, nonce.Bytes())
	if data == nil {
		return Deposit{}, false
	}

	var deposit Deposit
	if err := rlp.DecodeBytes(data, &deposit); err != nil {
		panic(fmt.Sprintf("deposit store corrupted: %s", err))
	}

	return deposit, true
}

func (store DepositStore) HasDeposit(ctx sdk.Context, nonce *big.Int) bool {
	return store.Has(ctx, nonce.Bytes())
}

func (store DepositStore) StoreDeposit(ctx sdk.Context, nonce *big.Int, deposit Deposit) {
	data, err := rlp.EncodeToBytes(&deposit)
	if err != nil {
		panic(fmt.Sprintf("error marshaling deposit with nonce %s: %s", nonce, err))
	}

	store.Set(ctx, nonce.Bytes(), data)
}

func (store DepositStore) SpendDeposit(ctx sdk.Context, nonce *big.Int, spender []byte) sdk.Result {
	deposit, ok := store.GetDeposit(ctx, nonce)
	if !ok {
		return ErrOutputDNE(DefaultCodespace, fmt.Sprintf("deposit with nonce %s does not exist", nonce)).Result()
	} else if deposit.Spent {
		return ErrOutputSpent(DefaultCodespace, fmt.Sprintf("deposit with nonce %s is already spent", nonce)).Result()
	}

	deposit.Spent = true
	deposit.Spender = spender

	store.StoreDeposit(ctx, nonce, deposit)

	return sdk.Result{}
}
