package crypto

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

type byter interface {
	Bytes() []byte
}

func checkAminoBinary(t *testing.T, src, dst interface{}, size int) {
	// Marshal to binary bytes.
	bz, err := cdc.MarshalBinaryBare(src)
	require.Nil(t, err, "%+v", err)
	if byterSrc, ok := src.(byter); ok {
		// Make sure this is compatible with current (Bytes()) encoding.
		require.Equal(t, byterSrc.Bytes(), bz, "Amino binary vs Bytes() mismatch")
	}
	// Make sure we have the expected length.
	if size != -1 {
		require.Equal(t, size, len(bz), "Amino binary size mismatch")
	}
	// Unmarshal
	err = cdc.UnmarshalBinaryBare(bz, dst)
	require.Nil(t, err, "%+v", err)
}

func checkAminoJSON(t *testing.T, src interface{}, dst interface{}, isNil bool) {
	// Marshal to JSON bytes.
	js, err := cdc.MarshalJSON(src)
	require.Nil(t, err, "%+v", err)
	if isNil {
		require.Equal(t, string(js), `null`)
	} else {
		require.Contains(t, string(js), `"type":`)
		require.Contains(t, string(js), `"value":`)
	}
	// Unmarshal.
	err = cdc.UnmarshalJSON(js, dst)
	require.Nil(t, err, "%+v", err)
}

func checkJSONMarshalUnMarshal(t *testing.T, src interface{}, dst interface{}) {
	// Marshal to JSON bytes.
	jbytes, err := json.Marshal(src)
	require.Nil(t, err, "%+v", err)
	// Unmarshal.
	err = json.Unmarshal(jbytes, dst)
	require.Nil(t, err, "%+v", err)
}

func TestKeyEncodings(t *testing.T) {
	cases := []struct {
		privKey           PrivateKey
		privSize, pubSize int // binary sizes with the amino overhead
	}{
		{
			privKey:  PrivateKey(Ed25519PrivateKey{}).GenPrivateKey(),
			privSize: 69,
			pubSize:  37,
		},
		{
			privKey:  PrivateKey(Secp256k1PrivateKey{}).GenPrivateKey(),
			privSize: 37,
			pubSize:  38,
		},
	}

	for _, tc := range cases {

		// Check (de/en)codings of PrivKeys.
		var priv2, priv3 PrivateKey
		checkAminoBinary(t, tc.privKey, &priv2, tc.privSize)
		require.EqualValues(t, tc.privKey, priv2)
		checkAminoJSON(t, tc.privKey, &priv3, false) // TODO also check Prefix bytes.
		require.EqualValues(t, tc.privKey, priv3)

		// Check (de/en)codings of Sigs.
		var sig1, sig2 []byte
		sig1, err := tc.privKey.Sign([]byte("something"))
		require.NoError(t, err)
		checkAminoBinary(t, sig1, &sig2, -1) // Signature size changes for Secp anyways.
		require.EqualValues(t, sig1, sig2)

		// Check (de/en)codings of PubKeys.
		pubKey, err := PubKeyToPublicKey(tc.privKey.PubKey())
		assert.Nil(t, err)
		var pub2, pub3 PublicKey
		checkAminoBinary(t, pubKey, &pub2, tc.pubSize)
		require.EqualValues(t, pubKey, pub2)
		checkAminoJSON(t, pubKey, &pub3, false) // TODO also check Prefix bytes.
		require.EqualValues(t, pubKey, pub3)
	}
}

func TestEd25519PublicKeyCustomMarshalling(t *testing.T) {

	//Get a Random PubKey
	pub := getRandomPubKey(t)
	var pub2 Ed25519PublicKey
	//Do Marshalling and Unmarshalling
	checkJSONMarshalUnMarshal(t, pub, &pub2)
	//Pub2 should have the same value if everything is ok.
	require.EqualValues(t, pub, pub2)
}

func TestSecp256k1PublicKeyCustomMarshalling(t *testing.T) {

	//Get a Random PubKey
	pub := getRandomPubKeySecp(t)
	var pub2 Secp256k1PublicKey
	//Do Marshalling and Unmarshalling
	checkJSONMarshalUnMarshal(t, pub, &pub2)
	//Pub2 should have the same value if everything is ok.
	require.EqualValues(t, pub, pub2)
}

func TestNilEncodings(t *testing.T) {

	// Check nil Signature.
	var a, b []byte
	checkAminoJSON(t, &a, &b, true)
	require.EqualValues(t, a, b)

	// Check nil PublicKey.
	var c, d PublicKey
	checkAminoJSON(t, &c, &d, true)
	require.EqualValues(t, c, d)

	// Check nil PrivKey.
	var e, f PrivateKey
	checkAminoJSON(t, &e, &f, true)
	require.EqualValues(t, e, f)

}
