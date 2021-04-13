package keys

import (
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
)

// SigningAlgo defines an algorithm to derive key-pairs which can be used for cryptographic signing.
type SigningAlgo string

// Keybase exposes operations on a generic keystore
// Keybase only supports Ed25519 key pairs
// Optimization: Merge Keybase interface with LazyKeybase and Keybase impl into a single type
type Keybase interface {
	// CRUD on the keystore
	List() ([]KeyPair, error)
	Get(address types.Address) (KeyPair, error)
	Delete(address types.Address, passphrase string) error
	UnsafeDelete(address types.Address) error
	Update(address types.Address, oldpass string, newpass string) error
	GetCoinbase() (KeyPair, error)
	SetCoinbase(address types.Address) error
	// Sign some bytes, looking up the private key to use
	Sign(address types.Address, passphrase string, msg []byte) ([]byte, crypto.PublicKey, error)

	// Create a new KeyPair and encrypt it to disk using encryptPassphrase
	Create(encryptPassphrase string) (KeyPair, error)

	// ImportPrivKey using Armored private key string. Decrypts armor with decryptPassphrase, and stores locally using encryptPassphrase
	ImportPrivKey(armor, decryptPassphrase, encryptPassphrase string) (KeyPair, error)

	// ExportPrivKeyArmor using Armored private key string. Decrypts armor with decryptPassphrase, and encrypts result armor using the encryptPassphrase
	ExportPrivKeyEncryptedArmor(address types.Address, decryptPassphrase, encryptPassphrase, hint string) (armor string, err error)

	// ImportPrivateKeyObject using the raw unencrypted privateKey string and encrypts it to disk using encryptPassphrase
	ImportPrivateKeyObject(privateKey [64]byte, encryptPassphrase string) (KeyPair, error)

	// ExportPrivateKeyObject exports raw PrivKey object.
	ExportPrivateKeyObject(address types.Address, passphrase string) (crypto.PrivateKey, error)

	// CloseDB closes the database.
	CloseDB()
}

// KeyPair is the public information about a locally stored key
type KeyPair struct {
	PublicKey    crypto.PublicKey `json:"pubkey"`
	PrivKeyArmor string           `json:"privkey.armor"`
}

// NewKeyPair with the given public key and priv armor key
func NewKeyPair(pub crypto.PublicKey, privArmor string) KeyPair {
	return KeyPair{
		PublicKey:    pub,
		PrivKeyArmor: privArmor,
	}
}

// GetAddress for the given KeyPair
func (kp KeyPair) GetAddress() types.Address {
	return kp.PublicKey.Address().Bytes()
}

// encoding info
func writeKeyPair(kp KeyPair) []byte {
	bz, err := cdc.MarshalBinaryLengthPrefixed(kp)
	if err != nil {
		panic(err)
	}
	return bz
}

// decoding info
func readKeyPair(bz []byte) (kp KeyPair, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(bz, &kp)
	return
}
