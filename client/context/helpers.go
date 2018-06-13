package context

import (
	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/ethereum/go-ethereum/common"
	rlp "github.com/ethereum/go-ethereum/rlp"
	crypto "github.com/tendermint/go-crypto"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/types"
)

// Broadcast the transaction bytes to Tendermint
func (ctx ClientContext) BroadcastTx(tx []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxCommit(tx)
	if err != nil {
		return res, err
	}

	if res.CheckTx.Code != uint32(0) {
		return res, errors.Errorf("CheckTx failed: (%d) %s",
			res.CheckTx.Code, res.CheckTx.Log)
	}
	if res.DeliverTx.Code != uint32(0) {
		return res, errors.Errorf("DeliverTx failed: (%d) %s",
			res.DeliverTx.Code, res.DeliverTx.Log)
	}
	return res, err

}

// sign and build the transaction from the msg
func (ctx ClientContext) SignBuildBroadcast(name string, msg sdk.Msg) (res *ctypes.ResultBroadcastTxCommit, err error) {

	passphrase, err := ctx.GetPassphraseFromStdin(name)
	if err != nil {
		return nil, err
	}

	txBytes, err := ctx.SignAndBuild(name, passphrase, msg)
	if err != nil {
		return nil, err
	}

	return ctx.BroadcastTx(txBytes)
}

// Get the from address from the name flag
func (ctx ClientContext) GetFromAddress() (from common.Address, err error) {

	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	name := ctx.FromAddressName
	if name == "" {
		return nil, errors.Errorf("must provide a from address name")
	}

	info, err := keybase.Get(name)
	if err != nil {
		return nil, errors.Errorf("No key for: %s", name)
	}

	return info.PubKey.Address(), nil
}

// Sign and build the transaction
func (ctx ClientContext) SignAndBuild(name, passphrase string, msg types.SpendMsg) ([]byte, error) {

	keybase, err := client.GetKeyBase()
	if err != nil {
		return nil, err
	}

	bz := msg.GetSignBytes()

	sig, pubkey, err := keybase.Sign(name, passphrase, bz)
	if err != nil {
		return nil, err
	}

	sigs := [2]crypto.Signature{sig, sig}

	tx := types.NewBaseTx(msg, sigs)

	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Get passphrase from std input
func (ctx ClientContext) GetPassphraseFromStdin(name string) (pass string, err error) {
	buf := client.BufferStdin()
	prompt := fmt.Sprintf("Password to sign with '%s':", name)
	return client.GetPassword(prompt, buf)
}

// Prepares a simple rpc.Client
func (ctx ClientContext) GetNode() (rpcclient.Client, error) {
	if ctx.Client == nil {
		return nil, errors.New("Must define node URI")
	}
	return ctx.Client, nil
}
