package types

import (
	"github.com/pokt-network/pocket-core/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgStake{}, "pos/MsgStake", nil)
	cdc.RegisterConcrete(MsgBeginUnstake{}, "pos/MsgBeginUnstake", nil)
	cdc.RegisterConcrete(MsgUnjail{}, "pos/MsgUnjail", nil)
	cdc.RegisterConcrete(MsgSend{}, "pos/Send", nil)
}

var ModuleCdc *codec.Codec // generic sealed codec to be used throughout this module

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
