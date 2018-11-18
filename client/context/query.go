package context

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pkg/errors"

	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// GetNode returns an RPC client. If the context's client is not defined, an
// error is returned.
func (ctx ClientContext) GetNode() (rpcclient.Client, error) {
	if ctx.Client == nil {
		return nil, errors.New("no RPC client defined")
	}

	return ctx.Client, nil
}

// Query performs a query for information about the connected node.
func (ctx ClientContext) Query(path string, data cmn.HexBytes) (res []byte, err error) {
	return ctx.query(path, data)
}

// Query information about the connected node with a data payload
func (ctx ClientContext) QueryWithData(path string, data []byte) (res []byte, err error) {
	return ctx.query(path, data)
}

// QueryStore performs a query from a Tendermint node with the provided key and
// store name.
func (ctx ClientContext) QueryStore(key cmn.HexBytes, storeName string) (res []byte, err error) {
	return ctx.queryStore(key, storeName, "key")
}

// QuerySubspace performs a query from a Tendermint node with the provided
// store name and subspace.
func (ctx ClientContext) QuerySubspace(subspace []byte, storeName string) (res []sdk.KVPair, err error) {
	resRaw, err := ctx.queryStore(subspace, storeName, "subspace")
	if err != nil {
		return res, err
	}

	ctx.Codec.MustUnmarshalBinaryLengthPrefixed(resRaw, &res)
	return
}

// query performs a query from a Tendermint node with the provided store name
// and path.
func (ctx ClientContext) query(path string, key cmn.HexBytes) (res []byte, err error) {
	node, err := ctx.GetNode()
	if err != nil {
		return res, err
	}

	opts := rpcclient.ABCIQueryOptions{
		Height: ctx.Height,
		Prove:  !ctx.TrustNode,
	}

	result, err := node.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, errors.Errorf(resp.Log)
	}

	return resp.Value, nil
}

// queryStore performs a query from a Tendermint node with the provided a store
// name and path.
func (ctx ClientContext) queryStore(key cmn.HexBytes, storeName, endPath string) ([]byte, error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	return ctx.query(path, key)
}
