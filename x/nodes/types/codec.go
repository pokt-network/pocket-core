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
	cdc.RegisterStructure(MsgProtoStake{}, "pos/MsgProtoStake")
	cdc.RegisterStructure(MsgBeginUnstake{}, "pos/MsgBeginUnstake")
	cdc.RegisterStructure(MsgUnjail{}, "pos/MsgUnjail")
	cdc.RegisterStructure(MsgSend{}, "pos/Send")
	cdc.RegisterStructure(MsgStake{}, "pos/MsgStake")
	cdc.RegisterImplementation((*sdk.ProtoMsg)(nil), &MsgUnjail{}, &MsgBeginUnstake{}, &MsgSend{}, &MsgStake{})
	cdc.RegisterImplementation((*sdk.Msg)(nil), &MsgUnjail{}, &MsgBeginUnstake{}, &MsgSend{}, &MsgStake{})
	cdc.RegisterInterface("nodes/validatorI", (*exported.ValidatorI)(nil), &ValidatorProto{})
	ModuleCdc = cdc
}

var ModuleCdc *codec.Codec // generic sealed codec to be used throughout this module

func init() {
	ModuleCdc = codec.NewCodec(types.NewInterfaceRegistry())
	RegisterCodec(ModuleCdc)
	crypto.RegisterAmino(ModuleCdc.AminoCodec().Amino)
	ModuleCdc.AminoCodec().Seal()
}
