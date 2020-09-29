package types

import (
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/codec/types"
	sdk "github.com/pokt-network/pocket-core/types"
)

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.NewCodec(types.NewInterfaceRegistry())
	RegisterCodec(ModuleCdc)
	ModuleCdc.AminoCodec().Seal()
}

// RegisterCodec registers all necessary param module types with a given codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterStructure(MsgChangeParam{}, "gov/msg_change_param")
	cdc.RegisterStructure(MsgDAOTransfer{}, "gov/msg_dao_transfer")
	cdc.RegisterStructure(MsgUpgrade{}, "gov/msg_upgrade")
	cdc.RegisterInterface("x.interface.nil", (*interface{})(nil))
	cdc.RegisterStructure(ACL{}, "gov/non_map_acl")
	cdc.RegisterStructure(Upgrade{}, "gov/upgrade")
	cdc.RegisterImplementation((*sdk.ProtoMsg)(nil), &MsgChangeParam{}, &MsgDAOTransfer{}, &MsgUpgrade{})
	cdc.RegisterImplementation((*sdk.Msg)(nil), &MsgChangeParam{}, &MsgDAOTransfer{}, &MsgUpgrade{})
	ModuleCdc = cdc
}
