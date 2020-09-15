package keys

import (
	"github.com/pokt-network/pocket-core/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"

	"github.com/pokt-network/pocket-core/codec"
)

var cdc *codec.LegacyAmino

func init() {
	cdc = codec.NewLegacyAminoCodec()
	cryptoAmino.RegisterAmino(cdc.Amino)
	crypto.RegisterAmino(cdc.Amino)
	cdc.RegisterConcrete(KeyPair{}, "crypto/keys/keypair", nil)
	cdc.Seal()
}
