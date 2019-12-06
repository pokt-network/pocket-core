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
	cdc.RegisterConcrete(MsgProof{}, "pocketcore/proof", nil)
}
