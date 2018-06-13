package client

import (
	"bufio"
	"os"
	"path/filename"

	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/bgentry/speakeasy"
	"github.com/ethereum/go-ethereum/common"
	isatty "github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tendermint/cli"
	crypto "github.com/tendermint/go-crypto"
	keys "github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/go-crypto/keys/words"
	dbm "github.com/tendermint/tmlibs/db"
)

const (
	// Minimum acceptable password length
	MinPassLength = 8
	// Directory under root where we store the keys
	KeyDBName = "keys"
)

var keybase keys.Keybase

// Allows for reading prompts for stdin
func BufferStdin() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

// Build SpendMsg
func BuildMsg(from, addr1, addr2 common.Address, position1, position2 types.Position, confirmSigs1, confirmSigs2 [2]crypto.Signature, amount1, amount2, fee uint64) types.SpendMsg {
	return types.NewSpendMsg(position1.Blknum, position1.TxIndex, position1.Oindex, position1.DepositNum, from, confirmSigs1, position2.Blknum, position2.TxIndex, position.Oindex, from, confirmSigs2, addr1, amount1, addr2, amount2, fee)
}

// initialize a keybase on the confirguration
func GetKeyBase() (keys.Keybase, error) {
	rootDir := viper.GetString(cli.HomeFlag)
	if keybase == nil {
		db, err := dbm.NewGoLevelDB(KeyDBName, filepath.Join(rootDir, "keys"))
		if err != nil {
			return nil, err
		}

		keybase = New(db, words.MustLoadCodec("english"))
	}
	return keybase, nil

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
		return "", errors.New("Password must be at least %d characters", MinPassLength)
	}
	return pass, nil
}

// Prompts for a password twice to verify they match
func GetCheckPassword(prompt, prompt2 string, buf *bufio.Reader) (string, error) {
	if !inputIsTty() {
		return GetPassword(prompt, buf)
	}

	pass, err := GetPassword(prompt, buf)
	if err != nil {
		return "", err
	}
	pass2, err := GetPassword(prompt2, buf)
	if err != nil {
		return "", err
	}
	if pass != pass2 {
		return "", errors.New("Passphrases did not match")
	}
	return pass, nil
}

// Returns true iff we have an interactive prompt
func inputIsTty() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
}

func printInfo(info keys.Info) {
	fmt.Printf("NAME:\tADDRESS:\t\t\t\t\t\tPUBKEY:\n")
	fmt.Printf("%s\t%s\t%s\n", info.Name, info.PubKey.Address().Bytes(), info.PubKey)
}
