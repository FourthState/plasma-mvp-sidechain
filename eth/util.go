package eth

import (
	"bytes"
	"math/big"
)

const (
	// prefixes
	prefixSeperator       = "::"
	depositPrefix         = "deposit"
	transactionExitPrefix = "txExit"
	depositExitPrefix     = "depositExit"

<<<<<<< HEAD
	// keys
	lastCommittedBlock = "lastCommittedBlock"

=======
>>>>>>> b3167013cb609ec55bd2a944e44a4d169ed332c9
	// constants
	blockIndexFactor = 1000000
	txIndexFactor    = 10
)

func prefixKey(prefix string, key []byte) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte(prefix))
	buffer.Write([]byte(prefixSeperator))
	buffer.Write(key)
	return buffer.Bytes()
}

// [blockNumber, txIndex, outputIndex]
func calcPriority(position [3]*big.Int) *big.Int {
	bFactor := big.NewInt(blockIndexFactor)
	tFactor := big.NewInt(txIndexFactor)

	bFactor = bFactor.Mul(bFactor, position[0])
	tFactor = tFactor.Mul(tFactor, position[1])

<<<<<<< HEAD
	return new(big.Int).Add(bFactor, tFactor).Add(position[2], big.NewInt(0))
=======
	temp := new(big.Int).Add(bFactor, tFactor)
	return temp.Add(temp, position[2])
>>>>>>> b3167013cb609ec55bd2a944e44a4d169ed332c9
}
