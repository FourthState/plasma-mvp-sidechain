package store

import (
	"bufio"
	//	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

var currentStdin *bufio.Reader

func init() {
	currentStdin = bufio.NewReader(os.Stdin)
}

// overrideStdin allows to temporarily override stdin
func overrideStdin(newStdin *bufio.Reader) (cleanUp func()) {
	prevStdin := currentStdin
	currentStdin = newStdin
	cleanUp = func() {
		currentStdin = prevStdin
	}
	return cleanUp
}

// Add, Delete, Update and iterate through accounts
func TestAccounts(t *testing.T) {
	// setup testing env
	os.Mkdir("testing", os.ModePerm)
	viper.Set(DirFlag, os.ExpandEnv("./testing"))

	// cleanup
	defer func() {
		viper.Reset()
		os.RemoveAll("testing")
	}()

	InitKeystore()

	cases := []string{
		"mykey",
		"another-key",
		"!aADS_AS@#$%^&*()",
		"    last     key",
	}

	for i, n := range cases {
		cleanUp := overrideStdin(bufio.NewReader(strings.NewReader("test1234\ntest1234\n")))
		defer cleanUp()
		_, err := AddAccount(n)
		require.NoErrorf(t, err, "case %d: failed to add account %s", i, n)
	}

}
