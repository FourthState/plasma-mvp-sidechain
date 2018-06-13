package client

import (
	"strings"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/go-crypto/keys/words"
	dbm "github.com/tendermint/tmlibs/db"
)

// ethKeybase implements the keybase interface from tendermint's go-crypto
// Provides a full-featured key manager with encryption
type ethKeybase struct {
	kb crypto.dbKeyBase
}

func New(db dbm.DB, codec words.Codec) ethKeybase {
	return ethKeybase{
		kb: crypto.dbKeyBase{
			db:    db,
			codec: codec,
		},
	}
}

var _ Keybase = ethKeybase{}

// Create generates a new key using Ethereum's crypto GenerateKey() function
// and persists it to storage, encrypted using the passphrase. It returns
// the generated seedphrase and the key Info. It returns an error if it fails
// to generate a key for the given algo type, or if another key is
// already stored under the same name.
func (kb dbKeybase) CreateMnemonic(name, passphrase string, algo keys.CryptoAlgo) (keys.Info, string, error) {
	// 32 byte secret corresponds to 24 BIP39 words.
	secret := crypto.CRandBytes(32)
	priv, err := ethcrypto.ToECDSA(secret)
	if err != nil {
		return crypto.Info{}, "", err
	}

	// encrypt and persist the key
	info := ethKB.kb.writeKey(priv, name, passphrase)

	// return the mnemonic phrase
	words, err := ethKB.kb.codec.BytesToWords(secret)
	seed := strings.Join(words, " ")
	return info, seed, err
}

// Recover converts a seedphrase to an eth private key and persists it,
// encypted with the given passphrase. Functions like Create, but
// seedphrase is input not output.
func (ethKB ethKeybase) Recover(name, passphrase, seedphrase string) (keys.Info, error) {
	words := strings.Split(strings.TrimSpace(seedphrase, " "))
	secret, err := ethKB.kb.codec.WordsToBytes(words)
	if err != nil {
		return crypto.Info{}, err
	}

	priv, err := ethcrypto.ToECDSA(secret)
	if err != nil {
		return crypto.Info{}, err
	}

	// encrypt and persist key.
	public := ethKB.kb.writeKey(priv, name, passphrase)
	return public, err
}

func (ethKB ethKeybase) CreateLedger(name string, path crypto.DerivationPath, algo keys.SignAlgo) (keys.Info, error) {
	return ethKB.kb.CreateLedger(name, path, algo)
}

func (ethKB ethKeybase) CreateOffline(name string, pub crypto.PubKey) (keys.Info, error) {
	return ethKB.kb.CreateOffline(name, pub)
}

// List returns the keys from storage in alphabetical order.
func (ethKB ethKeybase) List() ([]keys.Info, error) {
	return ethKB.kb.List()
}

// Get returns the public information about one key.
func (ethKB ethKeybase) Get(name string) (keys.Info, error) {
	return ethKB.kb.Get(name)
}

// Sign signs the msg with the named key.
// It returns an error if the key doesn't exist or the decyption fails.
func (ethKB ethKeybase) Sign(name, passphrase string, msg []byte) (sig crypto.Signature, pub crypto.Signature, err error) {
	return ethKB.kb.Sign(name, passphrase, msg)
}

// Export exports the key corresponding to the name provided.
// Returns an error if no such key exists.
func (ethKB ethKeybase) Export(name string) (armor string, err error) {
	return ethKB.kb.Export(name)
}

// ExportPubKey returns public keys in ASCII armored format.
// Retrieve a Info object by its name and return the public key in
// a portable format.
func (ethKB ethKeybase) ExportPubKey(name string) (armor string, err error) {
	return ethKB.kb.ExportPubKey(name)
}

func (ethKB ethKeybase) Import(name string, armor string) (err error) {
	return ethKB.kb.Import(name, armor)
}

// ImportPubKey imports AsCII-armored public keys.
// Store a new Info object holding a public key only, i.e. it will
// not be possible to sign with it as it lacks the secret key.
func (ethKB ethKeybase) ImportPubKey(name string, armor string) (err error) {
	return ethKB.kb.ImportPubKey(name, armor)
}

// Delete removes the key forever, but the proper passphrase must
// be presented before it is deleted
func (ethKB ethKeybase) Delete(name, passphrase string) error {
	return ethKB.kb.Delete(name, passphrase)
}

// Update changes the passphrase with which an already stored key is
// encrypted.
// oldpass must be the current passphrase used for encryption
// newpass will be the only valid passphrase from this time forward.
func (ethKB ethKeybase) Update(name, oldpass, newpass string) error {
	return ethKB.kb.Update(name, oldpass, newpass)
}
