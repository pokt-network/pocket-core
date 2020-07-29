package types

import (
	"github.com/pokt-network/pocket-core/codec"
)

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers all necessary param module types with a given codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgChangeParam{}, "gov/msg_change_param", nil)
	cdc.RegisterConcrete(MsgDAOTransfer{}, "gov/msg_dao_transfer", nil)
	cdc.RegisterInterface((*interface{})(nil), nil)
	cdc.RegisterConcrete(ACL{}, "gov/non_map_acl", nil)
	cdc.RegisterConcrete(Upgrade{}, "gov/upgrade", nil)
	cdc.RegisterConcrete(MsgUpgrade{}, "gov/msg_upgrade", nil)
}
