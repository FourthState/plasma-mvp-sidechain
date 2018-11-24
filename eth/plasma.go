package eth

import (
	"crypto/ecdsa"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	contracts "github.com/fourthstate/plasma-mvp-sidechain/contracts/wrappers"
)

type Plamsa struct {
	privateKey *ecdsa.PrivateKey
	contract   *contracts.Plasma
}

type DepositEvent struct {
	Owner common.Address
	Value *big.Int
	Nonce *big.Int
}

func InitPlasma(contractAddr string, privateKey *ecdsa.PrivateKey, client *Client) (*Plasma, error) {
	plasmaContract, err = contracts.NewRootChain(common.HexToAddress(contractAddr), client.ec)
	if err != nil {
		return nil, err
	}

	return &Plasma{privateKey, plasmaContract}, nil
}

func (plasma *Plasma) SubmitBlock(header []byte) error {
	return nil
}

// Check deposit checks the validity
func (plasma *Plamsa) CheckDeposit(nonce types.Uint) (bool, err) {

	return false, nil
}
