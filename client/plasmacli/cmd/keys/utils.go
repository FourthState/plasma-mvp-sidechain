package keys

import (
	cli "github.com/FourthState/plasma-mvp-sidechain/client"
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
