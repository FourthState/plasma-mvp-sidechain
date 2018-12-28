package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	amino "github.com/tendermint/go-amino"
	"io"
)

// Wrapper around
type UTXO struct {
	InputKeys        [][]byte // keys to retrieve the inputs of this output
	ConfirmationHash [32]byte // confirmation hash of the input transaction

	Output   plasma.Output
	Spent    bool
	Position plasma.Position
}

type utxo struct {
	InputKeys        [][]byte
	ConfirmationHash [32]byte
	Output           []byte
	Spent            bool
	Position         []byte
}

func (u *UTXO) EncodeRLP(w io.Writer) error {
	utxo := utxo{u.InputKeys, u.ConfirmationHash, u.Output.Bytes(), u.Spent, u.Position.Bytes()}

	return rlp.Encode(w, &utxo)
}

func (u *UTXO) DecodeRLP(s *rlp.Stream) error {
	utxo := utxo{}
	if err := s.Decode(&utxo); err != nil {
		return err
	}
	if err := rlp.DecodeBytes(utxo.Output, &u.Output); err != nil {
		return err
	}
	if err := rlp.DecodeBytes(utxo.Position, &u.Position); err != nil {
		return err
	}

	u.InputKeys = utxo.InputKeys
	u.ConfirmationHash = utxo.ConfirmationHash
	u.Spent = utxo.Spent

	return nil
}

type UTXOStore struct {
	KVStore
	cdc *amino.Codec
}

func NewUTXOStore(ctxKey sdk.StoreKey) UTXOStore {
	return UTXOStore{
		KVStore: NewKVStore(ctxKey),
	}
}

func (store UTXOStore) GetUTXO(ctx sdk.Context, addr common.Address, pos plasma.Position) (UTXO, bool) {
	key := append(addr.Bytes(), pos.Bytes()...)

	data := store.Get(ctx, key)
	if data == nil {
		return UTXO{}, false
	}

	var utxo UTXO
	err := rlp.DecodeBytes(data, &utxo)
	if err != nil {
		panic(fmt.Sprintf("utxo store corrupted: %s", err))
	}

	return utxo, true
}

func (store UTXOStore) StoreUTXO(ctx sdk.Context, utxo UTXO) {
	key := append(utxo.Output.Owner.Bytes(), utxo.Position.Bytes()...)
	data, err := rlp.EncodeToBytes(&utxo)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling utxo: %s", err))
	}

	store.Set(ctx, key, data)
}

func (store UTXOStore) SpendUTXO(ctx sdk.Context, addr common.Address, pos plasma.Position, spenderKeys [][]byte) sdk.Error {
	utxo, ok := store.GetUTXO(ctx, addr, pos)
	if !ok {
		return sdk.ErrUnknownRequest("utxo does not exist")
	}
	if utxo.Spent {
		return sdk.ErrUnauthorized("utxo already marked as spent")
	}

	utxo.Spent = true
	utxo.InputKeys = spenderKeys

	store.StoreUTXO(ctx, utxo)

	return nil
}
