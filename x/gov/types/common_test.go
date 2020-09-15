package types

import (
	"math/rand"

	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/codec/types"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
)

// nolint: deadcode unused
// create a codec used only for testing
func makeTestCodec() *codec.Codec {
	var cdc = codec.NewCodec(types.NewInterfaceRegistry())
	auth.RegisterCodec(cdc)
	RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	crypto.RegisterAmino(cdc.AminoCodec().Amino)
	return cdc
}

func getRandomPubKey() crypto.Ed25519PublicKey {
	var pub crypto.Ed25519PublicKey
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	return pub
}

func getRandomValidatorAddress() sdk.Address {
	return sdk.Address(getRandomPubKey().Address())
}

var testACL ACL

func createTestACL() ACL {
	if testACL == nil {
		acl := ACL(make([]ACLPair, 0))
		acl.SetOwner("auth/MaxMemoCharacters", getRandomValidatorAddress())
		acl.SetOwner("auth/TxSigLimit", getRandomValidatorAddress())
		acl.SetOwner("gov/daoOwner", getRandomValidatorAddress())
		acl.SetOwner("gov/acl", getRandomValidatorAddress())
		acl.SetOwner("gov/upgrade", getRandomValidatorAddress())
		testACL = acl
	}
	return testACL
}

func createTestAdjacencyMap() map[string]bool {
	m := make(map[string]bool)
	m["auth/MaxMemoCharacters"] = true // set
	m["auth/TxSigLimit"] = true        // set
	m["gov/daoOwner"] = true           // set
	m["gov/acl"] = true                // set
	m["gov/upgrade"] = true
	return m
}
