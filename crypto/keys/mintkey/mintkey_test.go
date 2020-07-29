package mintkey_test

import (
	"github.com/pokt-network/pocket-core/crypto"
	"testing"

	"github.com/pokt-network/pocket-core/crypto/keys"
	"github.com/pokt-network/pocket-core/crypto/keys/mintkey"
	"github.com/stretchr/testify/require"
)

func TestArmorUnarmorPrivKey(t *testing.T) {
	priv := crypto.Ed25519PrivateKey{}.GenPrivateKey()
	armor, _ := mintkey.EncryptArmorPrivKey(priv, "passphrase", "")
	_, err := mintkey.UnarmorDecryptPrivKey(armor, "wrongpassphrase")
	require.Error(t, err)
	decrypted, err := mintkey.UnarmorDecryptPrivKey(armor, "passphrase")
	require.NoError(t, err)
	require.True(t, priv.Equals(decrypted))
}

func TestArmorUnarmorPrivKeySecp(t *testing.T) {
	priv := crypto.Secp256k1PrivateKey{}.GenPrivateKey()
	armor, _ := mintkey.EncryptArmorPrivKey(priv, "passphrase", "")
	_, err := mintkey.UnarmorDecryptPrivKey(armor, "wrongpassphrase")
	require.Error(t, err)
	decrypted, err := mintkey.UnarmorDecryptPrivKey(armor, "passphrase")
	require.NoError(t, err)
	require.True(t, priv.Equals(decrypted))
}

func TestArmorUnarmorPubKey(t *testing.T) {
	// Select the encryption and storage for your cryptostore
	cstore := keys.NewInMemory()
	// Add keys and see they return in alphabetical order
	kp, err := cstore.Create("passphrase")
	require.NoError(t, err)
	armor := mintkey.ArmorPubKeyBytes(kp.PublicKey.RawBytes())
	pubBytes, err := mintkey.UnarmorPubKeyBytes(armor)
	require.NoError(t, err)
	pub, err := crypto.NewPublicKeyBz(pubBytes)
	require.NoError(t, err)
	require.True(t, pub.Equals(kp.PublicKey))
}
