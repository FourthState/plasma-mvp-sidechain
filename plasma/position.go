package plasma

import (
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"strconv"
	"strings"
)

// Position -
type Position struct {
	BlockNum     *big.Int
	TxIndex      uint16
	OutputIndex  uint8
	DepositNonce *big.Int
}

const (
	blockIndexFactor = 1000000
	txIndexFactor    = 10
)

// NewPosition
func NewPosition(blkNum *big.Int, txIndex uint16, oIndex uint8, depositNonce *big.Int) Position {
	if depositNonce == nil {
		depositNonce = big.NewInt(0)
	}
	if blkNum == nil {
		blkNum = big.NewInt(0)
	}

	return Position{
		BlockNum:     blkNum,
		TxIndex:      txIndex,
		OutputIndex:  oIndex,
		DepositNonce: depositNonce,
	}
}

// Bytes serializes `p`
func (p Position) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(&p)
	return bytes
}

// ValidateBasic ensures deposit and child chain positions are mutually exclusive.
// Block numbering starts at 1
func (p Position) ValidateBasic() error {
	if p.IsNilPosition() {
		return fmt.Errorf("nil position is not a valid position")
	}

	// deposit position
	if p.IsDeposit() {
		if p.BlockNum.Sign() > 0 || p.TxIndex > 0 || p.OutputIndex > 0 {
			return fmt.Errorf("chain position must be all zero if a deposit nonce is specified. (0.0.0.nonce)")
		}
	} else {
		if p.BlockNum.Sign() == 0 {
			return fmt.Errorf("block numbering starts at 1")
		}
		if p.OutputIndex > 1 {
			return fmt.Errorf("output index must be 0 or 1")
		}
	}

	return nil
}

// IsDeposit -
func (p Position) IsDeposit() bool {
	return p.DepositNonce.Sign() != 0
}

// IsFee -
func (p Position) IsFee() bool {
	return p.TxIndex == 1<<16-1
}

// IsNilPosition -
func (p Position) IsNilPosition() bool {
	return p.BlockNum.Sign() == 0 && p.DepositNonce.Sign() == 0
}

// Priority
func (p Position) Priority() *big.Int {
	if p.IsDeposit() {
		return p.DepositNonce
	}

	bFactor := big.NewInt(blockIndexFactor)
	tFactor := big.NewInt(txIndexFactor)

	bFactor = bFactor.Mul(bFactor, p.BlockNum)
	tFactor = tFactor.Mul(tFactor, big.NewInt(int64(p.TxIndex)))

	temp := new(big.Int).Add(bFactor, tFactor)
	return temp.Add(temp, big.NewInt(int64(p.OutputIndex)))
}

// ToBigIntArray -
func (p Position) ToBigIntArray() [4]*big.Int {
	return [4]*big.Int{p.BlockNum, big.NewInt(int64(p.TxIndex)), big.NewInt(int64(p.OutputIndex)), p.DepositNonce}
}

func (p Position) String() string {
	if p.BlockNum == nil {
		p.BlockNum = big.NewInt(0)
	}
	if p.DepositNonce == nil {
		p.DepositNonce = big.NewInt(0)
	}
	return fmt.Sprintf("(%s.%d.%d.%s)",
		p.BlockNum, p.TxIndex, p.OutputIndex, p.DepositNonce)
}

// FromPositionString constructs a `Position` from the string representation
// "(blockNumber, txIndex, outputIndex, depositNonce)". The returned error will
// also reflect an invalid position object in addition to any deserialization errors
func FromPositionString(posStr string) (Position, error) {
	posStr = strings.TrimSpace(posStr)
	if string(posStr[0]) != "(" || string(posStr[len(posStr)-1]) != ")" {
		return Position{}, fmt.Errorf("position must be enclosed in parens. (blockNum,txIndex,oIndex,depositNonce)")
	}

	// remove the parens
	posStr = posStr[1 : len(posStr)-1]

	blkNum := new(big.Int)
	depositNonce := new(big.Int)

	var txIndex uint16
	var oIndex uint8

	tokens := strings.Split(posStr, ".")
	if len(tokens) != 4 {
		return Position{},
			fmt.Errorf("invalid position. positions follow (blockNum.txIndex.oIndex.depositNonce). ex: (1.0.0.0)")
	}

	var err error
	var ok bool
	var num uint64
	for i, token := range tokens {
		token = strings.TrimSpace(token)
		if i == 0 {
			blkNum, ok = blkNum.SetString(token, 10)
			if !ok {
				return Position{}, fmt.Errorf("error parsing the block number")
			}
		} else if i == 1 {
			num, err = strconv.ParseUint(token, 0, 16)
			txIndex = uint16(num)
		} else if i == 2 {
			num, err = strconv.ParseUint(token, 0, 8)
			oIndex = uint8(num)
		} else {
			depositNonce, ok = depositNonce.SetString(token, 10)
			if !ok {
				return Position{}, fmt.Errorf("error parsing the deposit nonce")
			}
		}

		if err != nil {
			return Position{}, err
		}
	}

	pos := NewPosition(blkNum, txIndex, oIndex, depositNonce)
	return pos, pos.ValidateBasic()
}

// FromExitKey creates the position of a utxo given the exit key
func FromExitKey(key *big.Int, deposit bool) Position {
	if deposit {
		return NewPosition(big.NewInt(0), 0, 0, key)
	}

	return NewPosition(
		new(big.Int).Div(key, big.NewInt(blockIndexFactor)),
		uint16(key.Int64()%blockIndexFactor/txIndexFactor),
		uint8(key.Int64()%2),
		big.NewInt(0),
	)
}
