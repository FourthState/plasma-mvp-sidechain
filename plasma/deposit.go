package plasma

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

type Deposit struct {
	Owner       common.Address `json:"Owner"`
	Amount      *big.Int       `json:"Amount"`
	EthBlockNum *big.Int       `json:"EthBlockNum"`
}

type deposit struct {
	Owner       common.Address
	Amount      []byte
	EthBlockNum []byte
}

func NewDeposit(owner common.Address, amount *big.Int, ethBlockNum *big.Int) Deposit {
	return Deposit{
		Owner:       owner,
		Amount:      amount,
		EthBlockNum: ethBlockNum,
	}
}

func (d *Deposit) EncodeRLP(w io.Writer) error {
	deposit := &deposit{d.Owner, d.Amount.Bytes(), d.EthBlockNum.Bytes()}

	return rlp.Encode(w, deposit)
}

func (d *Deposit) DecodeRLP(s *rlp.Stream) error {
	var dep deposit
	if err := s.Decode(&dep); err != nil {
		return err
	}

	d.Owner = dep.Owner
	d.Amount = new(big.Int).SetBytes(dep.Amount)
	d.EthBlockNum = new(big.Int).SetBytes(dep.EthBlockNum)

	return nil
}
