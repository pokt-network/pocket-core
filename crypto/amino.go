package crypto

import (
	"github.com/pokt-network/pocket-core/codec"
	"github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

var cdc = codec.NewLegacyAminoCodec()

func init() {
	RegisterAmino(cdc.Amino)
	cryptoAmino.RegisterAmino(cdc.Amino)
}

// RegisterAmino registers all go-crypto related types in the given (amino) codec.
func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterInterface((*PublicKey)(nil), nil)
	cdc.RegisterInterface((*PrivateKey)(nil), nil)
	cdc.RegisterConcrete(Ed25519PublicKey{}, "crypto/ed25519_public_key", nil)
	cdc.RegisterConcrete(Ed25519PrivateKey{}, "crypto/ed25519_private_key", nil)
	cdc.RegisterConcrete(Secp256k1PublicKey{}, "crypto/secp256k1_public_key", nil)
	cdc.RegisterConcrete(Secp256k1PrivateKey{}, "crypto/secp256k1_private_key", nil)
	cdc.RegisterInterface((*MultiSig)(nil), nil)
	cdc.RegisterInterface((*PublicKeyMultiSig)(nil), nil)
	cdc.RegisterConcrete(PublicKeyMultiSignature{}, "crypto/public_key_multi_signature", nil)
	cdc.RegisterConcrete(MultiSignature{}, "crypto/multi_signature", nil)
}
