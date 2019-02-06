package plasma

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

// Output represents the outputs of a transaction
type Output struct {
	Owner  common.Address
	Amount *big.Int
}

type output struct {
	Owner  common.Address
	Amount []byte
}

func NewOutput(owner common.Address, amount *big.Int) Output {
	if amount == nil {
		amount = big.NewInt(0)
	}

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

// ValidateBasic ensures either a nil or valid output
func (o Output) ValidateBasic() error {
	if utils.IsZeroAddress(o.Owner) {
		if o.Amount.Sign() > 0 {
			return fmt.Errorf("amount specified with a nil owner")
		}
	} else {
		if o.Amount.Sign() == 0 {
			return fmt.Errorf("cannot send a zero amount")
		}
	}

	return nil
}

func (o Output) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(&o)
	return bytes
}

func (o Output) String() string {
	return fmt.Sprintf("Owner: %x, Amount: %s", o.Owner, o.Amount)
}
