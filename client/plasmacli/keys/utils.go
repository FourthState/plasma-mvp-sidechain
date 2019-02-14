package keys

import (
	"fmt"
	cli "github.com/FourthState/plasma-mvp-sidechain/client"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"path/filepath"
)

// Return directory in which account data is stored
func AccountDir() string {
	dir := viper.GetString(cli.DirFlag)
	if dir[len(dir)-1] != '/' {
		dir = filepath.Join(dir, "/")
	}
	return os.ExpandEnv(filepath.Join(dir, accountDir))
}

func printAccount(name string, address ethcmn.Address) {
	fmt.Printf("%s\t\t%v\n", name, address.Hex())
}

// Open leveldb using command line flag, attempt to retrieve account address
// If an error occurs, close the db
// Return db so that it can be closed later if successful
func OpenAndGet(name string) (*leveldb.DB, ethcmn.Address, error) {
	dir := AccountDir()
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, ethcmn.Address{}, err
	}

	key, err := rlp.EncodeToBytes(name)
	if err != nil {
		db.Close()
		return nil, ethcmn.Address{}, err
	}

	address, err := db.Get(key, nil)
	if err != nil {
		db.Close()
		return nil, ethcmn.Address{}, err
	}

	return db, ethcmn.BytesToAddress(address), nil
}
