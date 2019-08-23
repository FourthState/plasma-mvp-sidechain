package store

import (
	"crypto/ecdsa"
	"fmt"
	cosmoscli "github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"io/ioutil"
	"path/filepath"
)

const (
	MinPasswordLength = 8

	NewPassphrasePrompt       = "Enter new passphrase for your key:"
	NewPassphrasePromptRepeat = "Repeat passphrase:"

	PassphrasePrompt = "Enter passphrase:"

	accountsDir = "data/accounts.ldb"
	keysDir     = "keys"
)

var (
	home string
	ks   *keystore.KeyStore
)

// InitKeystore initializes a keystore in the specified directory
func InitKeystore(homeDir string) {
	home = homeDir

	dir := getDir(keysDir)
	if ks == nil {
		ks = keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	}
}

// AccountIterator returns an iterator for accounts.
// CONTRACT: Caller is responsible for closing db after use.
func AccountIterator() (iterator.Iterator, *leveldb.DB) {
	dir := getDir(accountsDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		fmt.Printf("leveldb: %s", err)
		return nil, nil
	}

	return db.NewIterator(nil, nil), db
}

// AddAccount adds a new account to the keystore
func AddAccount(name string) (ethcmn.Address, error) {
	dir := getDir(accountsDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return ethcmn.Address{}, fmt.Errorf("leveldb: %s", err)
	}
	defer db.Close()

	key := []byte(name)
	if _, err = db.Get(key, nil); err == nil {
		return ethcmn.Address{}, fmt.Errorf("you are trying to override an existing private key name. Please delete it first")
	}

	buf := cosmoscli.BufferStdin()
	password, err := cosmoscli.GetCheckPassword(NewPassphrasePrompt, NewPassphrasePromptRepeat, buf)
	if err != nil {
		return ethcmn.Address{}, err
	}

	acc, err := ks.NewAccount(password)
	if err != nil {
		return ethcmn.Address{}, fmt.Errorf("keystore: %s", err)
	}

	if err = db.Put(key, acc.Address.Bytes(), nil); err != nil {
		return ethcmn.Address{}, fmt.Errorf("leveldb: %s", err)
	}

	return acc.Address, nil
}

// GetAccount retrieves the address of an account.
func GetAccount(name string) (ethcmn.Address, error) {
	dir := getDir(accountsDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return ethcmn.Address{}, fmt.Errorf("leveldb: %s", err)
	}
	defer db.Close()

	addr, err := db.Get([]byte(name), nil)
	if err == leveldb.ErrNotFound {
		return ethcmn.Address{}, fmt.Errorf("account does not exist")
	} else if err != nil {
		return ethcmn.Address{}, fmt.Errorf("leveldb: %s", err)
	}

	return ethcmn.BytesToAddress(addr), nil
}

// DeleteAccount removes an account from keystore
func DeleteAccount(name string) error {
	dir := getDir(accountsDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return fmt.Errorf("leveldb: %s", err)
	}
	defer db.Close()

	addr, err := db.Get([]byte(name), nil)
	if err == leveldb.ErrNotFound {
		return fmt.Errorf("account does not exist")
	} else if err != nil {
		return fmt.Errorf("leveldb: %s", err)
	}

	buf := cosmoscli.BufferStdin()
	password, err := cosmoscli.GetPassword(PassphrasePrompt, buf)
	if err != nil {
		return err
	}

	acc := accounts.Account{
		Address: ethcmn.BytesToAddress(addr),
	}

	if err = ks.Delete(acc, password); err != nil {
		return fmt.Errorf("keystore: %s", err)
	}

	return nil
}

// UpdateAccount updates either the name of an account or the passphrase for
// an account.
func UpdateAccount(name string, updatedName string) (msg string, err error) {
	dir := getDir(accountsDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return msg, fmt.Errorf("leveldb: %s", err)
	}
	defer db.Close()

	key := []byte(name)
	addr, err := db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return msg, fmt.Errorf("account does not exist")
	} else if err != nil {
		return msg, fmt.Errorf("leveldb: %s", err)
	}

	if updatedName != "" {
		// Update key name
		if err = db.Delete(key, nil); err != nil {
			return msg, fmt.Errorf("leveldb: %s", err)
		}

		if err = db.Put([]byte(updatedName), addr, nil); err != nil {
			return msg, fmt.Errorf("leveldb: %s", err)
		}
		msg = "Account name has been updated."
	} else {
		// Update passphrase
		buf := cosmoscli.BufferStdin()
		password, err := cosmoscli.GetPassword(PassphrasePrompt, buf)
		updatedPassword, err := cosmoscli.GetCheckPassword(NewPassphrasePrompt, NewPassphrasePromptRepeat, buf)
		if err != nil {
			return msg, err
		}

		acc := accounts.Account{
			Address: ethcmn.BytesToAddress(addr),
		}
		if err = ks.Update(acc, password, updatedPassword); err != nil {
			return msg, fmt.Errorf("keystore: %s", err)
		}
		msg = "Account passphrase has been updated."
	}

	return msg, nil
}

// ImportECDSA imports a private key with associated an account name.
func ImportECDSA(name string, pk *ecdsa.PrivateKey) (ethcmn.Address, error) {
	dir := getDir(accountsDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return ethcmn.Address{}, fmt.Errorf("leveldb: %s", err)
	}
	defer db.Close()

	buf := cosmoscli.BufferStdin()
	password, err := cosmoscli.GetCheckPassword(NewPassphrasePrompt, NewPassphrasePromptRepeat, buf)
	if err != nil {
		return ethcmn.Address{}, err
	}

	acct, err := ks.ImportECDSA(pk, password)
	if err != nil {
		return ethcmn.Address{}, fmt.Errorf("keystore: %s", err)
	}

	if err = db.Put([]byte(name), acct.Address.Bytes(), nil); err != nil {
		return ethcmn.Address{}, fmt.Errorf("leveldb: %s", err)
	}

	return acct.Address, nil

}

// SignHashWithPassphrase will sign over the provided hash if the the passphrase
// provided through user interaction is correct.
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

	var sig []byte
	sig, err = ks.SignHashWithPassphrase(acc, password, hash)
	if err != nil {
		return nil, fmt.Errorf("keystore: %s", err)
	}

	return sig, nil
}

// GetKey returns the private key mapped to the provided key name.
func GetKey(name string) (*ecdsa.PrivateKey, error) {
	dir := getDir(accountsDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, fmt.Errorf("leveldb: %s", err)
	}
	defer db.Close()

	addr, err := db.Get([]byte(name), nil)
	if err == leveldb.ErrNotFound {
		return nil, fmt.Errorf("account does not exist")
	} else if err != nil {
		return nil, fmt.Errorf("leveldb: %s", err)
	}

	acc, err := ks.Find(
		accounts.Account{
			Address: ethcmn.BytesToAddress(addr),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("keystore: %s", err)
	}

	bz, err := ioutil.ReadFile(acc.URL.Path)
	if err != nil {
		return nil, fmt.Errorf("ioutil: %s", err)
	}

	buf := cosmoscli.BufferStdin()
	pass, err := cosmoscli.GetPassword(PassphrasePrompt, buf)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(bz, pass)
	if err != nil {
		return nil, fmt.Errorf("keystore: %s", err)
	}
	return key.PrivateKey, nil
}

// returns the directory specified by the --directory flag
// with the passed in string appended to the end
func getDir(location string) string {
	return filepath.Join(home, location)
}
