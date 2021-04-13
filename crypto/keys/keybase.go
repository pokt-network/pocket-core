package keys

import (
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"

	"github.com/pkg/errors"

	"github.com/pokt-network/pocket-core/crypto/keys/mintkey"
	"github.com/pokt-network/pocket-core/types"

	dbm "github.com/tendermint/tm-db"
)

var _ Keybase = &dbKeybase{}

// dbKeybase combines encryption and storage implementation to provide
// a full-featured key manager
type dbKeybase struct {
	db       dbm.DB
	coinbase KeyPair
}

// newDbKeybase creates a new keybase instance using the passed DB for reading and writing keys.
func newDbKeybase(db dbm.DB) Keybase {
	return &dbKeybase{
		db: db,
	}
}

// NewInMemory creates a transient keybase on top of in-memory storage
// instance useful for testing purposes and on-the-fly key generation.
func NewInMemory() Keybase { return &dbKeybase{db: dbm.NewMemDB()} }

func (kb *dbKeybase) GetCoinbase() (KeyPair, error) {
	if kb.coinbase.PrivKeyArmor == "" {
		kps, err := kb.List()
		if err != nil {
			return KeyPair{}, err
		}
		if len(kps) == 0 {
			return KeyPair{}, fmt.Errorf("0 keypairs in the keybase, so could not get a coinbase")
		}
		kb.coinbase = kps[0]
	}
	return kb.coinbase, nil
}

func (kb *dbKeybase) SetCoinbase(address types.Address) error {
	kp, err := kb.Get(address)
	if err != nil {
		return err
	}
	kb.coinbase = kp
	return nil
}

// List returns the keys from storage in alphabetical order.
func (kb dbKeybase) List() ([]KeyPair, error) {
	var res []KeyPair
	iter, _ := kb.db.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		kp, err := readKeyPair(iter.Value())
		if err != nil {
			return nil, err
		}
		res = append(res, kp)
	}
	return res, nil
}

// Get returns the public information about one key.
func (kb dbKeybase) Get(address types.Address) (KeyPair, error) {
	ik, _ := kb.db.Get(addrKey(address))
	if len(ik) == 0 {
		return KeyPair{}, fmt.Errorf("key with address %s not found", address)
	}
	return readKeyPair(ik)
}

// Delete removes key forever, but we must present the
// proper passphrase before deleting it (for security).
// It returns an error if the key doesn't exist or
// passphrases don't match.
func (kb dbKeybase) Delete(address types.Address, passphrase string) error {
	// verify we have the key in the keybase
	kp, err := kb.Get(address)
	if err != nil {
		return err
	}

	// Verify passphrase matches
	if _, err = mintkey.UnarmorDecryptPrivKey(kp.PrivKeyArmor, passphrase); err != nil {
		return err
	}

	return kb.db.DeleteSync(addrKey(kp.GetAddress()))
}

// Delete without passphrase verification
func (kb *dbKeybase) UnsafeDelete(address types.Address) error {
	// verify we have the key in the keybase
	kp, err := kb.Get(address)
	if err != nil {
		return err
	}
	return kb.db.DeleteSync(addrKey(kp.GetAddress()))
}

// Update changes the passphrase with which an already stored key is
// encrypted.
//
// oldpass must be the current passphrase used for encryption,
// getNewpass is a function to get the passphrase to permanently replace
// the current passphrase
func (kb dbKeybase) Update(address types.Address, oldpass string, newpass string) error {
	kp, err := kb.Get(address)
	if err != nil {
		return err
	}

	privKey, err := mintkey.UnarmorDecryptPrivKey(kp.PrivKeyArmor, oldpass)
	if err != nil {
		return err
	}

	_, err = kb.writeLocalKeyPair(privKey, newpass, "")
	if err != nil {
		return err
	}
	return nil
}

// Sign signs the msg with the named key.
// It returns an error if the key doesn't exist or the decryption fails.
func (kb dbKeybase) Sign(address types.Address, passphrase string, msg []byte) ([]byte, crypto.PublicKey, error) {
	kp, err := kb.Get(address)
	if err != nil {
		return nil, nil, err
	}

	if kp.PrivKeyArmor == "" {
		err = fmt.Errorf("private key not available")
		return nil, nil, err
	}

	priv, err := mintkey.UnarmorDecryptPrivKey(kp.PrivKeyArmor, passphrase)
	if err != nil {
		return nil, nil, err
	}

	sig, err := priv.Sign(msg)
	if err != nil {
		return nil, nil, err
	}

	pub := priv.PublicKey()
	return sig, pub, nil
}

// Create a new KeyPair and encrypt it to disk using encryptPassphrase
func (kb dbKeybase) Create(encryptPassphrase string) (KeyPair, error) {
	privKey := crypto.PrivateKey(crypto.Ed25519PrivateKey{}).GenPrivateKey()
	kp, err := kb.writeLocalKeyPair(privKey, encryptPassphrase, "")
	if err != nil {
		return kp, err
	}
	return kp, nil
}

// ImportPrivKey imports a private key in ASCII armor format.
// It returns an error if a key with the same address exists or a wrong decryptPassphrase is
// supplied.
func (kb dbKeybase) ImportPrivKey(armor, decryptPassphrase, encryptPassphrase string) (KeyPair, error) {
	privKey, err := mintkey.UnarmorDecryptPrivKey(armor, decryptPassphrase)
	if err != nil {
		return KeyPair{}, err
	}
	Address, err := types.AddressFromHex(privKey.PubKey().Address().String())
	if err != nil {
		return KeyPair{}, err
	}
	if _, err := kb.Get(Address); err == nil {
		return KeyPair{}, errors.New("Cannot overwrite key with address: " + Address.String())
	}
	return kb.writeLocalKeyPair(privKey, encryptPassphrase, "")
}

// ExportPrivKeyEncryptedArmor finds the KeyPair by the address, decrypts the armor private key,
// and returns an encrypted armored private key string
func (kb dbKeybase) ExportPrivKeyEncryptedArmor(address types.Address, decryptPassphrase, encryptPassphrase, hint string) (armor string, err error) {
	priv, err := kb.ExportPrivateKeyObject(address, decryptPassphrase)
	if err != nil {
		return "", err
	}
	return mintkey.EncryptArmorPrivKey(priv, encryptPassphrase, hint)
}

// ImportPrivateKeyObject using the raw unencrypted privateKey string and encrypts it to disk using encryptPassphrase
func (kb dbKeybase) ImportPrivateKeyObject(privateKey [64]byte, encryptPassphrase string) (KeyPair, error) {
	ed25519PK := crypto.Ed25519PrivateKey(privateKey)
	Address, err := types.AddressFromHex(ed25519PK.PubKey().Address().String())
	if err != nil {
		return KeyPair{}, err
	}
	if _, err := kb.Get(Address); err == nil {
		return KeyPair{}, errors.New("Cannot overwrite key with address: " + Address.String())
	}
	return kb.writeLocalKeyPair(ed25519PK, encryptPassphrase, "")
}

// ExportPrivateKeyObject exports raw PrivKey object.
func (kb dbKeybase) ExportPrivateKeyObject(address types.Address, passphrase string) (crypto.PrivateKey, error) {
	kp, err := kb.Get(address)
	if err != nil {
		return nil, err
	}

	if kp.PrivKeyArmor == "" {
		err = fmt.Errorf("private key not available")
		return nil, err
	}

	priv, err := mintkey.UnarmorDecryptPrivKey(kp.PrivKeyArmor, passphrase)
	if err != nil {
		return nil, err
	}
	return priv, err
}

// CloseDB releases the lock and closes the storage backend.
func (kb dbKeybase) CloseDB() {
	_ = kb.db.Close()
}

// Private interface
func (kb dbKeybase) writeLocalKeyPair(priv crypto.PrivateKey, passphrase, hint string) (KeyPair, error) {
	// encrypt private key using passphrase
	privArmor, err := mintkey.EncryptArmorPrivKey(priv, passphrase, hint)
	if err != nil || privArmor == "" {
		fmt.Println(err)
		return KeyPair{}, err
	}
	// make Info
	pub := priv.PublicKey()
	localKeyPair := NewKeyPair(pub, privArmor)
	kb.writeKeyPair(localKeyPair)

	return localKeyPair, nil
}

func (kb dbKeybase) writeKeyPair(kp KeyPair) {
	// write the info by key
	key := addrKey(kp.GetAddress())
	serializedInfo := writeKeyPair(kp)
	_ = kb.db.SetSync(key, serializedInfo)
}

func addrKey(address types.Address) []byte {
	return []byte(address.String())
}
