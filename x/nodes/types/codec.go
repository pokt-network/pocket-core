package types

import (
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/codec/types"
	"github.com/pokt-network/pocket-core/crypto"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterStructure(LegacyMsgProtoStake{}, "pos/MsgProtoStake")
	cdc.RegisterStructure(LegacyMsgBeginUnstake{}, "pos/MsgBeginUnstake")
	cdc.RegisterStructure(LegacyMsgUnjail{}, "pos/MsgUnjail")
	cdc.RegisterStructure(MsgSend{}, "pos/Send")
	cdc.RegisterStructure(LegacyMsgStake{}, "pos/MsgStake")
	cdc.RegisterStructure(MsgUnjail{}, "pos/8.0MsgUnjail")
	cdc.RegisterStructure(MsgBeginUnstake{}, "pos/8.0MsgBeginUnstake")
	cdc.RegisterStructure(MsgProtoStake{}, "pos/8.0MsgProtoStake")
	cdc.RegisterStructure(MsgStake{}, "pos/8.0MsgStake")
	cdc.RegisterImplementation((*sdk.ProtoMsg)(nil), &MsgUnjail{}, &MsgBeginUnstake{}, &MsgSend{}, &MsgStake{},
		&LegacyMsgUnjail{}, &LegacyMsgBeginUnstake{}, &LegacyMsgStake{})
	cdc.RegisterImplementation((*sdk.Msg)(nil), &MsgUnjail{}, &MsgBeginUnstake{}, &MsgSend{}, &MsgStake{},
		&LegacyMsgUnjail{}, &LegacyMsgBeginUnstake{}, &LegacyMsgStake{})
	cdc.RegisterInterface("nodes/validatorI", (*exported.ValidatorI)(nil), &Validator{}, &LegacyValidator{})
	ModuleCdc = cdc
}

var ModuleCdc *codec.Codec // generic sealed codec to be used throughout this module

func init() {
	ModuleCdc = codec.NewCodec(types.NewInterfaceRegistry())
	RegisterCodec(ModuleCdc)
	crypto.RegisterAmino(ModuleCdc.AminoCodec().Amino)
	ModuleCdc.AminoCodec().Seal()
}
