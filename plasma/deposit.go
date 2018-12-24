package plasma

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Deposit struct {
	Owner       common.Address `json:"Owner"`
	Amount      *big.Int       `json:"Amount"`
	EthBlockNum *big.Int       `json:"EthBlockNum"`
}

func NewDeposit(owner common.Address, amount *big.Int, ethBlockNum *big.Int) Deposit {
	return Deposit{
		Owner:       owner,
		Amount:      amount,
		EthBlockNum: ethBlockNum,
	}
}
