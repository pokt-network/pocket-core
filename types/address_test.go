package types_test

import (
	"encoding/hex"
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/pokt-network/pocket-core/types"
)

var invalidStrs = []string{
	"hello, world!",
	"0xAA",
	"AAA",
}

func testMarshal(t *testing.T, original interface{}, res interface{}, marshal func() ([]byte, error), unmarshal func([]byte) error) {
	bz, err := marshal()
	require.Nil(t, err)
	err = unmarshal(bz)
	require.Nil(t, err)
	require.Equal(t, original, res)
}

func TestEmptyAddresses(t *testing.T) {
	require.Equal(t, (types.Address{}).String(), "")

	Addr, err := types.AddressFromHex("")
	require.True(t, Addr.Empty())
	require.Nil(t, err)

}

func TestYAMLMarshalers(t *testing.T) {
	addr := crypto.GenerateSecp256k1PrivKey().PubKey().Address()

	address := types.Address(addr)

	got, _ := yaml.Marshal(&address)
	require.Equal(t, address.String()+"\n", string(got))
}

func TestRandHexAddrConsistency(t *testing.T) {
	var pub ed25519.PubKeyEd25519

	for i := 0; i < 1000; i++ {
		_, err := rand.Read(pub[:])
		if err != nil {
			_ = err
		}

		acc := types.Address(pub.Address())
		res := types.Address{}

		testMarshal(t, &acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		testMarshal(t, &acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res, err = types.AddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)

		str = hex.EncodeToString(acc)
		res, err = types.AddressFromHex(str)
		require.Nil(t, err)
		require.Equal(t, acc, res)
	}

	for _, str := range invalidStrs {
		_, err := types.AddressFromHex(str)
		require.NotNil(t, err)

		_, err = types.AddressFromHex(str)
		require.NotNil(t, err)

		err = (*types.Address)(nil).UnmarshalJSON([]byte("\"" + str + "\""))
		require.NotNil(t, err)
	}
}

func TestAddressInterface(t *testing.T) {
	var pub ed25519.PubKeyEd25519
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}

	addrs := []types.AddressI{
		types.Address(pub.Address()),
	}

	for _, addr := range addrs {
		switch addr := addr.(type) {
		case types.Address:
			_, err := types.AddressFromHex(addr.String())
			require.Nil(t, err)
		default:
			t.Fail()
		}
	}

}

func TestCustomAddressVerifier(t *testing.T) {
	// Create a 10 byte address
	addr := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	address := types.Address(addr).String()
	// Verifiy that the default logic rejects this 10 byte address
	err := types.VerifyAddressFormat(addr)
	require.NotNil(t, err)
	_, err = types.AddressFromHex(address)
	require.NotNil(t, err)

	// Set a custom address verifier that accepts 10 or 20 byte Addresses
	types.GetConfig().SetAddressVerifier(func(bz []byte) error {
		n := len(bz)
		if n == 10 || n == types.AddrLen {
			return nil
		}
		return fmt.Errorf("incorrect address length %d", n)
	})

	// Verifiy that the custom logic accepts this 10 byte address
	err = types.VerifyAddressFormat(addr)
	require.Nil(t, err)
	_, err = types.AddressFromHex(address)
	require.Nil(t, err)
}
