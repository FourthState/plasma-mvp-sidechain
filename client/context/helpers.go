package context

import (
	"fmt"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	rlp "github.com/ethereum/go-ethereum/rlp"
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
func (ctx ClientContext) SignBuildBroadcast(addr common.Address, msg types.SpendMsg, dir string) (res *ctypes.ResultBroadcastTxCommit, err error) {

	passphrase, err := ctx.GetPassphraseFromStdin(addr)
	if err != nil {
		return nil, err
	}

	txBytes, err := ctx.SignAndBuild(addr, passphrase, msg, dir)
	if err != nil {
		return nil, err
	}

	return ctx.BroadcastTx(txBytes)
}

// Get the from address from the name flag
func (ctx ClientContext) GetFromAddress(dir string) (from common.Address, err error) {

	ks := client.GetKeyStore(dir)

	addrStr := ctx.FromAddress
	if addrStr == "" {
		return common.Address{}, errors.Errorf("must provide an address to send from")
	}
	addr, err := client.StrToAddress(addrStr)
	if err != nil {
		return common.Address{}, err
	}

	if !ks.HasAddress(addr) {
		return common.Address{}, errors.Errorf("No account for: %X", addr)
	}

	return addr, nil
}

// Sign and build the transaction
func (ctx ClientContext) SignAndBuild(addr common.Address, passphrase string, msg types.SpendMsg, dir string) ([]byte, error) {

	ks := client.GetKeyStore(dir)
	acc := accounts.Account{
		Address: addr,
	}
	acct, err := ks.Find(acc)
	if err != nil {
		return nil, err
	}

	bz := msg.GetSignBytes()
	// May need to change so hash is done in sendtx.go
	hash := ethcrypto.Keccak256(bz)

	sig, err := ks.SignHashWithPassphrase(acct, passphrase, hash)
	if err != nil {
		return nil, err
	}

	sigs := []types.Signature{types.Signature{sig}, types.Signature{sig}}

	tx := types.NewBaseTx(msg, sigs)

	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Get passphrase from std input
func (ctx ClientContext) GetPassphraseFromStdin(addr common.Address) (pass string, err error) {
	buf := client.BufferStdin()
	prompt := fmt.Sprintf("Password to sign with '%X':", addr)
	return client.GetPassword(prompt, buf)
}

// Prepares a simple rpc.Client
func (ctx ClientContext) GetNode() (rpcclient.Client, error) {
	if ctx.Client == nil {
		return nil, errors.New("must define node URI")
	}
	return ctx.Client, nil
}
