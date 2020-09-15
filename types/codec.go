package types

import "github.com/pokt-network/pocket-core/codec"

// // Register the sdk message type
// func RegisterCodec(cdc *codec.Codec) {
// 	amino.RegisterInterface((*Msg)(nil), nil)
// 	amino.RegisterInterface((*Tx)(nil), nil)
// 	proto.Register("types/msg", (*Msg)(nil))
// 	proto.Register("types/tx", (*Tx)(nil))
// }
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface("types/msg", (*Msg)(nil))
	cdc.RegisterInterface("types/Legacymsg", (*LegacyMsg)(nil))
	cdc.RegisterInterface("types/tx", (*Tx)(nil))
}
