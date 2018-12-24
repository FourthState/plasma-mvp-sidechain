package plasma

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// Output represents the outputs of a transaction
type Output struct {
	Owner  common.Address `json:"Owner"`
	Amount *big.Int       `json:"Amount"`
}

func NewOutput(owner common.Address, amount *big.Int) Output {
	return Output{
		Owner:  owner,
		Amount: amount,
	}
}

func (o *Output) Bytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(o.Owner.Bytes())
	buffer.Write(o.Amount.Bytes())

	return buffer.Bytes()
}
