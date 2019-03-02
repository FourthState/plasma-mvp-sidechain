package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	signatureDir = "data/signatures.ldb"
)

// Saves confirmation signature
// Overrides any value previously set at the given position
func SaveSig(position plasma.Position, sig []byte) error {
	dir := getDir(signatureDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return fmt.Errorf("failed to open db for signatures: { %s }", err)
	}
	defer db.Close()

	k := getSigKey(position)
	if err := db.Put(k, sig, nil); err != nil {
		return fmt.Errorf("failed to save confirmation signature: { %s }", err)
	}

	return nil
}

// Retrieves confirmation signautres
func GetSig(position plasma.Position) ([]byte, error) {
	dir := getDir(signatureDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open db for signature: { %s }", err)
	}
	defer db.Close()

	k := getSigKey(position)
	if sig, err := db.Get(k, nil); err != nil {
		return nil, fmt.Errorf("failed to get signature: { %s }", err)
	} else {
		return sig, nil
	}
}

// return key used for confirm signature mapping
func getSigKey(pos plasma.Position) []byte {
	return append(pos.BlockNum.Bytes(), []byte(string(pos.TxIndex))...)
}
