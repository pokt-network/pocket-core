package crypto

import "testing"

func getRandomPrivateKey(t *testing.T) Ed25519PrivateKey {
	return Ed25519PrivateKey{}.GenPrivateKey().(Ed25519PrivateKey)
}

func getRandomPrivateKeySecp(t *testing.T) Secp256k1PrivateKey {
	return Secp256k1PrivateKey{}.GenPrivateKey().(Secp256k1PrivateKey)
}

func getRandomPubKey(t *testing.T) Ed25519PublicKey {
	pk := getRandomPrivateKey(t)
	return pk.PublicKey().(Ed25519PublicKey)
}

func getRandomPubKeySecp(t *testing.T) Secp256k1PublicKey {
	pk := getRandomPrivateKeySecp(t)
	return pk.PublicKey().(Secp256k1PublicKey)
}
