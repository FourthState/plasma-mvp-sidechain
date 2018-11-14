package client

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/bgentry/speakeasy"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	isatty "github.com/mattn/go-isatty"
	"github.com/pkg/errors"
)

const (
	// Minimum acceptable password length
	MinPassLength = 8

	// Flags
	FlagNode      = "node"
	FlagHeight    = "height"
	FlagTrustNode = "trust-node"
	FlagAddress   = "address"
)

var ks *keystore.KeyStore

// Allows for reading prompts for stdin
func BufferStdin() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

// Build SpendMsg
func BuildMsg(inaddr0, inaddr1, addr0, addr1 common.Address, position0, position1 types.PlasmaPosition, confirmSigs0, confirmSigs1 [2][65]byte, amount0, amount1, fee uint64) types.SpendMsg {
	return types.SpendMsg{
		Blknum0:      position0.Blknum,
		Txindex0:     position0.TxIndex,
		Oindex0:      position0.Oindex,
		DepositNum0:  position0.DepositNum,
		Owner0:       inaddr0,
		ConfirmSigs0: confirmSigs0,
		Blknum1:      position1.Blknum,
		Txindex1:     position1.TxIndex,
		Oindex1:      position1.Oindex,
		DepositNum1:  position1.DepositNum,
		Owner1:       inaddr1,
		ConfirmSigs1: confirmSigs1,
		Newowner0:    addr0,
		Amount0:      amount0,
		Newowner1:    addr1,
		Amount1:      amount1,
		FeeAmount:    fee,
	}
}

// initialize a keystore in the specified directory
func GetKeyStore(dir string) *keystore.KeyStore {
	if ks == nil {
		ks = keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	}
	return ks
}

// Prompts for a password one-time
// Enforces minimum password length
func GetPassword(prompt string, buf *bufio.Reader) (pass string, err error) {
	if inputIsTty() {
		pass, err = speakeasy.Ask(prompt)
	} else {
		pass, err = readLineFromBuf(buf)
	}
	if err != nil {
		return "", err
	}
	if len(pass) < MinPassLength {
		return "", fmt.Errorf("Password must be at least %d characters", MinPassLength)
	}
	return pass, nil
}

// Prompts for a password twice to verify they match
func GetCheckPassword(prompt, prompt1 string, buf *bufio.Reader) (string, error) {
	if !inputIsTty() {
		return GetPassword(prompt, buf)
	}

	pass, err := GetPassword(prompt, buf)
	if err != nil {
		return "", err
	}
	pass1, err := GetPassword(prompt1, buf)
	if err != nil {
		return "", err
	}
	if pass != pass1 {
		return "", errors.New("Passphrases did not match")
	}
	return pass, nil
}

// value in a position defaults to 0 if not provided
func ParsePositions(posStr string) (position [2]types.PlasmaPosition, err error) {
	for i, v := range strings.Split(posStr, "::") {
		var pos [4]uint64
		for k, number := range strings.Split(v, ".") {
			pos[k], err = strconv.ParseUint(strings.TrimSpace(number), 0, 64)
			if err != nil {
				return [2]types.PlasmaPosition{}, err
			}
		}
		position[i] = types.NewPlasmaPosition(pos[0], uint16(pos[1]), uint8(pos[2]), uint64(pos[3]))
	}
	return position, nil
}

// Amounts will default to 0 if not provided
func ParseAmounts(amtStr string) (amount [3]uint64, err error) {
	for i, v := range strings.Split(amtStr, ",") {
		amount[i], err = strconv.ParseUint(strings.TrimSpace(v), 0, 64)
		if err != nil {
			return [3]uint64{}, err
		}
	}
	return amount, nil

}

// Convert string to Ethereum Address
func StrToAddress(addrStr string) (common.Address, error) {
	if !common.IsHexAddress(strings.TrimSpace(addrStr)) {
		return common.Address{}, errors.New("invalid address provided, please use hex format")
	}
	return common.HexToAddress(addrStr), nil
}

// Returns true iff we have an interactive prompt
func inputIsTty() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
}

// reads one line from stdin
func readLineFromBuf(buf *bufio.Reader) (string, error) {
	pass, err := buf.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(pass), nil
}
