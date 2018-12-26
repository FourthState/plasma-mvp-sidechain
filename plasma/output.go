package plasma

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

// Output represents the outputs of a transaction
type Output struct {
	Owner  common.Address `json:"Owner"`
	Amount *big.Int       `json:"Amount"`
}

type output struct {
	Owner  common.Address
	Amount []byte
}

func NewOutput(owner common.Address, amount *big.Int) Output {
	return Output{
		Owner:  owner,
		Amount: amount,
	}
}

func (o *Output) EncodeRLP(w io.Writer) error {
	output := output{o.Owner, o.Amount.Bytes()}

	return rlp.Encode(w, output)
}

func (o *Output) DecodeRLP(s *rlp.Stream) error {
	var output output
	if err := s.Decode(&output); err != nil {
		return err
	}

	o.Owner = output.Owner
	o.Amount = new(big.Int).SetBytes(output.Amount)

	return nil
}

func (o Output) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(&o)
	return bytes
}
