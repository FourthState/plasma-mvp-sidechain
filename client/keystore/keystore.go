package keystore

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/bgentry/speakeasy"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	isatty "github.com/mattn/go-isatty"
	"os"
)

const (
	MinPasswordLength = 8

	NewPassphrasePrompt       = "Enter new passphrase for your key:"
	NewPassphrasePromptRepeat = "Repeat passphrase:"

	PassphrasePrompt = "Enter passphrase:"
)

var ks *keystore.KeyStore

// initialize a keystore in the specified directory
func InitKeystore(dir string) {
	if ks == nil {
		ks = keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	}
}

func Accounts() []accounts.Account {
	return ks.Accounts()
}

// Prompts for a password one-time
// Enforces minimum password length
func promptPassword(prompt, repeatPrompt string) (string, error) {
	if !isatty.IsTerminal(os.Stdin.Fd()) && !isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		return "", fmt.Errorf("Only interactive terminals are supported")
	}

	password0, err := speakeasy.Ask(prompt)
	if err != nil {
		return "", err
	}

	if repeatPrompt != "" {
		password1, err := speakeasy.Ask(repeatPrompt)
		if err != nil {
			return "", err
		}

		if password0 != password1 {
			return "", fmt.Errorf("Passphrases do not match")
		}
	}

	if len(password0) < MinPasswordLength {
		return "", fmt.Errorf("Password must be at least %d characters", MinPasswordLength)
	}

	return password0, nil
}

func NewAccount() (common.Address, error) {
	password, err := promptPassword(NewPassphrasePrompt, NewPassphrasePromptRepeat)
	if err != nil {
		return common.Address{}, err
	}

	acc, err := ks.NewAccount(password)
	if err != nil {
		return common.Address{}, err
	}

	return acc.Address, nil
}

func Find(addr common.Address) (accounts.Account, error) {
	acc := accounts.Account{
		Address: addr,
	}

	return ks.Find(acc)
}

func Delete(addr common.Address) error {
	password, err := promptPassword(PassphrasePrompt, "")
	if err != nil {
		return err
	}

	acc := accounts.Account{
		Address: addr,
	}

	return ks.Delete(acc, password)
}

func Update(addr common.Address) error {
	password, err := promptPassword(PassphrasePrompt, "")
	newPassword, err := promptPassword(NewPassphrasePrompt, "")
	if err != nil {
		return err
	}

	acc := accounts.Account{
		Address: addr,
	}

	return ks.Update(acc, password, newPassword)
}

func ImportECDSA(key *ecdsa.PrivateKey) (accounts.Account, error) {
	password, err := promptPassword(NewPassphrasePrompt, NewPassphrasePromptRepeat)
	if err != nil {
		return accounts.Account{}, err
	}

	return ks.ImportECDSA(key, password)
}

func SignHashWithPassphrase(signer common.Address, hash []byte) ([]byte, error) {
	acc, err := Find(signer)
	if err != nil {
		return nil, err
	}

	password, err := promptPassword(PassphrasePrompt, "")
	if err != nil {
		return nil, err
	}

	return ks.SignHashWithPassphrase(acc, password, hash)
}
