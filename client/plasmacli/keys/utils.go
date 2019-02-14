package keys

import (
	"fmt"
	cli "github.com/FourthState/plasma-mvp-sidechain/client"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// Return directory in which account data is stored
func accountDir() string {
	dir := viper.GetString(cli.DirFlag)
	if dir[len(dir)-1] != '/' {
		dir = filepath.Join(dir, "/")
	}
	return os.ExpandEnv(filepath.Join(dir, "accounts.ldb"))
}

func printAccount(name string, address ethcmn.Address) {
	fmt.Printf("%s\t\t%v\n", name, address.Hex())
}
