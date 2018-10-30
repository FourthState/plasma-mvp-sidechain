package context

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	rlp "github.com/ethereum/go-ethereum/rlp"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/FourthState/plasma-mvp-sidechain/client"
	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
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
func (ctx ClientContext) SignBuildBroadcast(addrs [2]common.Address, msg types.SpendMsg, dir string) (res *ctypes.ResultBroadcastTxCommit, err error) {

	sig, err := ctx.GetSignature(addrs[0], msg, dir)
	if err != nil {
		return nil, err
	}
	sigs := []types.Signature{types.Signature{sig}}

	if utils.ValidAddress(addrs[1]) {
		sig, err = ctx.GetSignature(addrs[1], msg, dir)
		if err != nil {
			return nil, err
		}
		sigs = append(sigs, types.Signature{sig})
	}

	tx := types.NewBaseTx(msg, sigs)

	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	return ctx.BroadcastTx(txBytes)
}

func (ctx ClientContext) GetSignature(addr common.Address, msg utxo.SpendMsg, dir string) (sig []byte, err error) {

	passphrase, err := ctx.GetPassphraseFromStdin(addr)
	if err != nil {
		return nil, err
	}
	ks := client.GetKeyStore(dir)
	acc := accounts.Account{
		Address: addr,
	}

	acct, err := ks.Find(acc)
	if err != nil {
		return nil, err
	}

	bz := msg.GetSignBytes()
	hash := ethcrypto.Keccak256(bz)

	sig, err = ks.SignHashWithPassphrase(acct, passphrase, hash)
	if err != nil {
		return nil, err
	}
	return sig, nil

}

// Get the from address from the name flag
func (ctx ClientContext) GetInputAddresses(dir string) (from [2]common.Address, err error) {

	ks := client.GetKeyStore(dir)

	addrsStr := ctx.InputAddresses
	if addrsStr == "" {
		return [2]common.Address{}, errors.Errorf("must provide an address to send from")
	}

	addrs := strings.Split(addrsStr, ",")
	// first input address
	from[0], err = client.StrToAddress(strings.TrimSpace(addrs[0]))
	if err != nil {
		return [2]common.Address{}, err
	}
	if len(addrs) > 1 {
		// second input address
		from[1], err = client.StrToAddress(strings.TrimSpace(addrs[1]))
		if err != nil {
			return [2]common.Address{}, err
		}
	}

	if !ks.HasAddress(from[0]) {
		return [2]common.Address{}, errors.Errorf("no account for: %s", from[0].Hex())
	}
	if len(from) > 1 && !utils.ZeroAddress(from[1]) && !ks.HasAddress(from[1]) {
		return [2]common.Address{}, errors.Errorf("no account for: %s", from[1].Hex())
	}

	return from, nil
}

// Get passphrase from std input
func (ctx ClientContext) GetPassphraseFromStdin(addr common.Address) (pass string, err error) {
	buf := client.BufferStdin()
	prompt := fmt.Sprintf("Password to sign with '%s':", addr.Hex())
	return client.GetPassword(prompt, buf)
}
