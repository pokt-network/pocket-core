package types

import (
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/codec/types"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth/exported"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface("x.auth.ModuleAccount", (*exported.ModuleAccountI)(nil), &ModuleAccountEncodable{})
	cdc.RegisterInterface("x.auth.Account", (*exported.Account)(nil), &BaseAccountEncodable{}, &ModuleAccountEncodable{})
	cdc.RegisterInterface("x.auth.Supply", (*exported.SupplyI)(nil), &Supply{})
	cdc.RegisterStructure(&BaseAccount{}, "posmint/Account")
	cdc.RegisterStructure(LegacyStdTx{}, "posmint/StdTx")
	cdc.RegisterStructure(&Supply{}, "posmint/Supply")
	cdc.RegisterStructure(&ModuleAccount{}, "posmint/ModuleAccount")
	cdc.RegisterImplementation((*sdk.Tx)(nil), &StdTx{})
}

// module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.NewCodec(types.NewInterfaceRegistry())
	RegisterCodec(ModuleCdc)
	crypto.RegisterAmino(ModuleCdc.AminoCodec().Amino)
}
