package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	signatureDir = "data/signatures.ldb"
)

// Saves confirmation signatures
// If prepend is set to true, sig will be prepended to the currently stored signatures
// Otherwise it will be appended
// This ordering should be determined by input order in the transaction
// If the length of the currently stored signatures is 130 an error is returned
func SaveSig(position plasma.Position, sig []byte, prepend bool) error {
	if len(sig) != 65 {
		return fmt.Errorf("signature must have a length of 65 bytes")
	}

	signatures, err := GetSig(position)
	if len(signatures) == 130 {
		return fmt.Errorf("two signatures already exist for the given position")
	}

	if prepend {
		signatures = append(sig, signatures...)
	} else {
		signatures = append(signatures, sig...)
	}

	dir := getDir(signatureDir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return fmt.Errorf("failed to open db for signatures: { %s }", err)
	}
	defer db.Close()

	k := getSigKey(position)
	if err := db.Put(k, signatures, nil); err != nil {
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
