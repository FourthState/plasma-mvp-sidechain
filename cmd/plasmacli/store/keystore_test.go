package store

import (
	"bufio"
	cosmoscli "github.com/cosmos/cosmos-sdk/client"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

// Add, Delete, Update and iterate through accounts
func TestAccounts(t *testing.T) {
	// setup testing env
	os.Mkdir("testing", os.ModePerm)
	InitKeystore("./testing")

	// cleanup
	defer func() {
		viper.Reset()
		os.RemoveAll("testing")
	}()

	cases := []string{
		"mykey",
		"another-key",
		"!aADS_AS@#$%^&*()",
		"    last     key",
	}

	// Check adding and getting accounts
	for i, n := range cases {
		cleanUp := cosmoscli.OverrideStdin(bufio.NewReader(strings.NewReader("test1234\ntest1234\n")))
		defer cleanUp()
		addr1, err := AddAccount(n)
		require.NoErrorf(t, err, "case %d: failed to add account %s", i, n)

		addr2, err := GetAccount(n)
		require.NoError(t, err, "case %d: failed to get account %s", i, n)
		require.Equal(t, addr1, addr2, "case %d: address added and retrieved not equal for account %s", i, n)
	}

	// Check that iterator matches
	iter, db := AccountIterator()
	i := 0
	var actualNames []string
	var actualAddrs []ethcmn.Address
	for iter.Next() {
		actualNames = append(actualNames, string(iter.Key()))
		actualAddrs = append(actualAddrs, ethcmn.BytesToAddress(iter.Value()))
		i++
	}
	iter.Release()
	db.Close()

	for i, n := range actualNames {
		expectedAddr, err := GetAccount(string(n))
		require.NoError(t, err, "case %d: failed to get account %s", i, n)
		require.Equal(t, actualAddrs[i], expectedAddr, "case %d: address returned from iterator does not match actual address for account %s", i, n)
		i++
	}

	updatedNames := []string{
		"key",
		"k",
		"sdfghjk",
		"plasma",
	}

	// Update Account name and password, exports keys, imports them, then Delete Accounts
	for i, n := range cases {
		_, err := UpdateAccount(cases[i], updatedNames[i])
		require.NoError(t, err, "case %d: failed to update account name %s", i, n)

		_, err = GetAccount(cases[i])
		require.Error(t, err, "case %d: retireved account for mapping that should not exist, account %s", i, n)
		_, err = GetAccount(updatedNames[i])
		require.NoError(t, err, "case %d: failed to retrieve account after updating, account %s", i, updatedNames[i])

		// Unsuccessful Export
		cleanUp := cosmoscli.OverrideStdin(bufio.NewReader(strings.NewReader("wrongpassword\nwrongpassword\n")))
		defer cleanUp()
		accjson, err := Export(updatedNames[i])
		assert.Error(t, err)

		// Export
		cleanUp = cosmoscli.OverrideStdin(bufio.NewReader(strings.NewReader("test1234\ntest1234\n")))
		defer cleanUp()
		accjson, err = Export(updatedNames[i])
		require.NoError(t, err, "case %d: failed to export account for account %s", i, updatedNames[i])

		// Delete
		cleanUp = cosmoscli.OverrideStdin(bufio.NewReader(strings.NewReader("test1234\n")))
		defer cleanUp()
		err = DeleteAccount(updatedNames[i])
		require.NoError(t, err, "case %d: failed to delete account %s", i, updatedNames[i])

		// Unsuccessful Import
		cleanUp = cosmoscli.OverrideStdin(bufio.NewReader(strings.NewReader("wrongpass\nwrongpass\n")))
		defer cleanUp()
		_, err = Import(updatedNames[i], accjson)
		assert.Error(t, err)

		// Import
		cleanUp = cosmoscli.OverrideStdin(bufio.NewReader(strings.NewReader("test1234\ntest1234\n")))
		defer cleanUp()
		_, err = Import(updatedNames[i], accjson)
		require.NoError(t, err, "case %d: failed to import account for account %s", i, updatedNames[i])

		// Update
		cleanUp = cosmoscli.OverrideStdin(bufio.NewReader(strings.NewReader("test1234\nnewpass1234\n")))
		defer cleanUp()
		_, err = UpdateAccount(updatedNames[i], "")
		require.NoError(t, err, "case %d: failed to update account passphrase for account %s", i, updatedNames[i])

		cleanUp = cosmoscli.OverrideStdin(bufio.NewReader(strings.NewReader("newpass1234\n")))
		defer cleanUp()
		err = DeleteAccount(updatedNames[i])
		require.NoError(t, err, "case %d: failed to delete account %s", i, updatedNames[i])
	}

}
