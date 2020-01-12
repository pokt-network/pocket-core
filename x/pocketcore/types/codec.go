package types

import (
	"github.com/pokt-network/posmint/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgClaim{}, "pocketcore/claim", nil)
	cdc.RegisterConcrete(MsgProof{}, "pocketcore/RelayProof", nil)
	cdc.RegisterConcrete(Relay{}, "pocketcore/relay", nil)
	cdc.RegisterConcrete(Session{}, "pocketcore/session", nil)
	cdc.RegisterConcrete(RelayResponse{}, "pocketcore/relay_response", nil)
}
