package store

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	cosmoscli "github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

const (
	MinPasswordLength = 8

	NewPassphrasePrompt       = "Enter new passphrase for your key:"
	NewPassphrasePromptRepeat = "Repeat passphrase:"

	PassphrasePrompt = "Enter passphrase:"

	accountDir = "data/accounts.ldb"
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
func AddAccount(name string) (ethcmn.Address, error) {
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

	buf := cosmoscli.BufferStdin()
	password, err := cosmoscli.GetCheckPassword(NewPassphrasePrompt, NewPassphrasePromptRepeat, buf)
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
func GetAccount(name string) (ethcmn.Address, error) {
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
func DeleteAccount(name string) error {
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

	buf := cosmoscli.BufferStdin()
	password, err := cosmoscli.GetPassword(PassphrasePrompt, buf)
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
func UpdateAccount(name string, updatedName string) error {
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
		buf := cosmoscli.BufferStdin()
		password, err := cosmoscli.GetPassword(PassphrasePrompt, buf)
		updatedPassword, err := cosmoscli.GetPassword(NewPassphrasePrompt, buf)
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

	buf := cosmoscli.BufferStdin()
	password, err := cosmoscli.GetCheckPassword(NewPassphrasePrompt, NewPassphrasePromptRepeat, buf)
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
	addr, err := GetAccount(signer)
	if err != nil {
		return nil, err
	}

	acc := accounts.Account{
		Address: addr,
	}

	buf := cosmoscli.BufferStdin()
	password, err := cosmoscli.GetPassword(PassphrasePrompt, buf)
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
