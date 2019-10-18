package blockchain

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/pokt-network/pocket-core/types"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(types.Node{}, "blockchain/Node", nil)
	cdc.RegisterConcrete(types.Application{}, "blockchain/Application", nil)
}
