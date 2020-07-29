package types

import (
	"github.com/pokt-network/pocket-core/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAppStake{}, "apps/MsgAppStake", nil)
	cdc.RegisterConcrete(MsgBeginAppUnstake{}, "apps/MsgAppBeginUnstake", nil)
	cdc.RegisterConcrete(MsgAppUnjail{}, "apps/MsgAppUnjail", nil)
}

var ModuleCdc *codec.Codec // generic sealed codec to be used throughout this module

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
