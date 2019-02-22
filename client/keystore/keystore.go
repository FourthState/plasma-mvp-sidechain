package keystore

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/bgentry/speakeasy"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	isatty "github.com/mattn/go-isatty"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"os"
	"path/filepath"
)

const (
	MinPasswordLength = 8

	NewPassphrasePrompt       = "Enter new passphrase for your key:"
	NewPassphrasePromptRepeat = "Repeat passphrase:"

	PassphrasePrompt = "Enter passphrase:"

	accountDir = "accounts.ldb"
	keysDir    = "keys"

	DirFlag = "directory"
)

var ks *keystore.KeyStore

// initialize a keystore in the specified directory
func InitKeystore() {
	dir := getDir(keysDir)
	if ks == nil {
		ks = keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	}
}

// Return iterator over accounts
// returns db so db.close can be called
func AccountIterator() (iterator.Iterator, *leveldb.DB) {
	dir := getDir(accountDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		fmt.Printf("FAILURE: %s", err)
		return nil, nil
	}

	return db.NewIterator(nil, nil), db
}

// Add a new account to the keystore
// Add account name and address to leveldb
func Add(name string) (ethcmn.Address, error) {
	dir := getDir(accountDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return ethcmn.Address{}, err
	}
	defer db.Close()

	key := []byte(name)
	if _, err = db.Get(key, nil); err == nil {
		return ethcmn.Address{}, errors.New("you are trying to override an existing private key name. Please delete it first")
	}

	password, err := promptPassword(NewPassphrasePrompt, NewPassphrasePromptRepeat)
	if err != nil {
		return ethcmn.Address{}, err
	}

	acc, err := ks.NewAccount(password)
	if err != nil {
		return ethcmn.Address{}, err
	}

	if err = db.Put(key, acc.Address.Bytes(), nil); err != nil {
		return ethcmn.Address{}, err
	}

	return acc.Address, nil
}

// Retrieve the address of an account
func Get(name string) (ethcmn.Address, error) {
	dir := getDir(accountDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return ethcmn.Address{}, err
	}
	defer db.Close()

	addr, err := db.Get([]byte(name), nil)
	if err != nil {
		return ethcmn.Address{}, err
	}

	return ethcmn.BytesToAddress(addr), nil
}

// Remove an account from the local keystore
// and the leveldb
func Delete(name string) error {
	dir := getDir(accountDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	addr, err := db.Get([]byte(name), nil)
	if err != nil {
		return err
	}

	password, err := promptPassword(PassphrasePrompt, "")
	if err != nil {
		return err
	}

	acc := accounts.Account{
		Address: ethcmn.BytesToAddress(addr),
	}

	return ks.Delete(acc, password)
}

// Update either the name of an account
// or the passphrase for an account
func Update(name string, updatedName string) error {
	dir := getDir(accountDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	key := []byte(name)
	addr, err := db.Get(key, nil)
	if err != nil {
		return err
	}

	if updatedName != "" {
		// Update key name
		if err = db.Delete(key, nil); err != nil {
			return err
		}

		key, err = rlp.EncodeToBytes(updatedName)
		if err != nil {
			return err
		}

		if err = db.Put(key, addr, nil); err != nil {
			return err
		}
		fmt.Println("Account name has been updated.")
	} else {
		// Update passphrase
		password, err := promptPassword(PassphrasePrompt, "")
		updatedPassword, err := promptPassword(NewPassphrasePrompt, "")
		if err != nil {
			return err
		}

		acc := accounts.Account{
			Address: ethcmn.BytesToAddress(addr),
		}
		err = ks.Update(acc, password, updatedPassword)
		if err != nil {
			return err
		}
		fmt.Println("Account passphrase has been updated.")
	}

	return nil
}

// Import a private key with an account name
func ImportECDSA(name string, pk *ecdsa.PrivateKey) (ethcmn.Address, error) {
	dir := getDir(accountDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return ethcmn.Address{}, err
	}
	defer db.Close()

	password, err := promptPassword(NewPassphrasePrompt, NewPassphrasePromptRepeat)
	if err != nil {
		return ethcmn.Address{}, err
	}

	acct, err := ks.ImportECDSA(pk, password)
	if err != nil {
		return ethcmn.Address{}, err
	}

	if err = db.Put([]byte(name), acct.Address.Bytes(), nil); err != nil {
		return ethcmn.Address{}, err
	}

	return acct.Address, nil

}

func SignHashWithPassphrase(signer string, hash []byte) ([]byte, error) {
	addr, err := Get(signer)
	if err != nil {
		return nil, err
	}

	acc := accounts.Account{
		Address: addr,
	}

	password, err := promptPassword(PassphrasePrompt, "")
	if err != nil {
		return nil, err
	}

	return ks.SignHashWithPassphrase(acc, password, hash)
}

// Return the directory specified by the --directory flag
// with the passed in string appended to the end
func getDir(location string) string {
	dir := viper.GetString(DirFlag)
	if dir[len(dir)-1] != '/' {
		dir = filepath.Join(dir, "/")
	}
	return os.ExpandEnv(filepath.Join(dir, location))
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
