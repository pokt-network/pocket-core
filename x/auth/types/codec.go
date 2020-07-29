package types

import (
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/x/auth/exported"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "posmint/Account", nil)
	cdc.RegisterConcrete(StdTx{}, "posmint/StdTx", nil)
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&ModuleAccount{}, "posmint/ModuleAccount", nil)
	cdc.RegisterConcrete(&Supply{}, "posmint/Supply", nil)
}

// module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
