package codec

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/pokt-network/pocket-core/codec/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

type Codec struct {
	protoCdc        *ProtoCodec
	legacyCdc       *LegacyAmino
	afterUpgradeMod bool
}

func NewCodec(anyUnpacker types.AnyUnpacker) *Codec {
	return &Codec{
		protoCdc:  NewProtoCodec(anyUnpacker),
		legacyCdc: NewLegacyAminoCodec(),
	}
}

var NotProtoCompatibleInterfaceError = errors.New("the interface passed for encoding does not implement proto marshaller")

func (cdc *Codec) RegisterStructure(o interface{}, name string) {
	cdc.legacyCdc.RegisterConcrete(o, name, nil)
}

func (cdc *Codec) RegisterInterface(name string, iface interface{}, impls ...proto.Message) {
	res, ok := cdc.protoCdc.anyUnpacker.(types.InterfaceRegistry)
	if !ok {
		panic("unable to convert protocodec.anyUnpacker into types.InterfaceRegistry")
	}
	res.RegisterInterface(name, iface, impls...)
	cdc.legacyCdc.Amino.RegisterInterface(iface, nil)
}

func (cdc *Codec) RegisterImplementation(iface interface{}, impls ...proto.Message) {
	res, ok := cdc.protoCdc.anyUnpacker.(types.InterfaceRegistry)
	if !ok {
		panic("unable to convert protocodec.anyUnpacker into types.InterfaceRegistry")
	}
	res.RegisterImplementations(iface, impls...)
}

func (cdc *Codec) MarshalBinaryBare(o interface{}) ([]byte, error) {
	p, ok := o.(ProtoMarshaler)
	if !ok {
		if cdc.afterUpgradeMod {
			return nil, NotProtoCompatibleInterfaceError
		}
		return cdc.legacyCdc.MarshalBinaryBare(o)
	} else {
		if cdc.afterUpgradeMod {
			return cdc.protoCdc.MarshalBinaryBare(p)
		}
		return cdc.legacyCdc.MarshalBinaryBare(p)
	}
}

func (cdc *Codec) MarshalBinaryLengthPrefixed(o interface{}) ([]byte, error) {
	p, ok := o.(ProtoMarshaler)
	if !ok {
		if cdc.afterUpgradeMod {
			return nil, NotProtoCompatibleInterfaceError
		}
		return cdc.legacyCdc.MarshalBinaryLengthPrefixed(o)
	} else {
		if cdc.afterUpgradeMod {
			return cdc.protoCdc.MarshalBinaryLengthPrefixed(p)
		}
		return cdc.legacyCdc.MarshalBinaryLengthPrefixed(p)
	}
}

func (cdc *Codec) UnmarshalBinaryBare(bz []byte, ptr interface{}) error {
	p, ok := ptr.(ProtoMarshaler)
	if !ok {
		if cdc.afterUpgradeMod {
			return NotProtoCompatibleInterfaceError
		}
		return cdc.legacyCdc.UnmarshalBinaryBare(bz, ptr)
	}
	err := cdc.protoCdc.UnmarshalBinaryBare(bz, p)
	if err == nil {
		return nil
	}
	return cdc.legacyCdc.UnmarshalBinaryBare(bz, ptr)
}

func (cdc *Codec) UnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}) error {
	p, ok := ptr.(ProtoMarshaler)
	if !ok {
		if cdc.afterUpgradeMod {
			return NotProtoCompatibleInterfaceError
		}
		return cdc.legacyCdc.UnmarshalBinaryLengthPrefixed(bz, ptr)
	}
	err := cdc.protoCdc.UnmarshalBinaryLengthPrefixed(bz, p)
	if err == nil {
		return nil
	}
	return cdc.legacyCdc.UnmarshalBinaryLengthPrefixed(bz, ptr)
}

func (cdc *Codec) ProtoMarshalBinaryBare(o ProtoMarshaler) ([]byte, error) {
	return cdc.protoCdc.MarshalBinaryBare(o)
}

func (cdc *Codec) LegacyMarshalBinaryBare(o interface{}) ([]byte, error) {
	return cdc.legacyCdc.MarshalBinaryBare(o)
}

func (cdc *Codec) ProtoUnmarshalBinaryBare(bz []byte, ptr ProtoMarshaler) error {
	return cdc.protoCdc.UnmarshalBinaryBare(bz, ptr)
}

func (cdc *Codec) LegacyUnmarshalBinaryBare(bz []byte, ptr interface{}) error {
	return cdc.legacyCdc.UnmarshalBinaryBare(bz, ptr)
}

func (cdc *Codec) ProtoMarshalBinaryLengthPrefixed(o ProtoMarshaler) ([]byte, error) {
	return cdc.protoCdc.MarshalBinaryLengthPrefixed(o)
}

func (cdc *Codec) LegacyMarshalBinaryLengthPrefixed(o interface{}) ([]byte, error) {
	return cdc.legacyCdc.MarshalBinaryLengthPrefixed(o)
}

func (cdc *Codec) ProtoUnmarshalBinaryLengthPrefixed(bz []byte, ptr ProtoMarshaler) error {
	return cdc.protoCdc.UnmarshalBinaryLengthPrefixed(bz, ptr)
}

func (cdc *Codec) LegacyUnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}) error {
	return cdc.legacyCdc.UnmarshalBinaryLengthPrefixed(bz, ptr)
}

func (cdc *Codec) MarshalJSONIndent(o interface{}, prefix string, indent string) ([]byte, error) {
	return cdc.legacyCdc.MarshalJSONIndent(o, prefix, indent)
}

func (cdc *Codec) MarshalJSON(o interface{}) ([]byte, error) {
	return cdc.legacyCdc.MarshalJSON(o)
}

func (cdc *Codec) UnmarshalJSON(bz []byte, o interface{}) error {
	return cdc.legacyCdc.UnmarshalJSON(bz, o)
}

func (cdc *Codec) MustMarshalJSON(o interface{}) []byte {
	bz, err := cdc.MarshalJSON(o)
	if err != nil {
		panic(err)
	}
	return bz
}

func (cdc *Codec) MustUnmarshalJSON(bz []byte, ptr interface{}) {
	err := cdc.UnmarshalJSON(bz, ptr)
	if err != nil {
		panic(err)
	}
}

func RegisterEvidences(legacy *LegacyAmino, _ *ProtoCodec) {
	tmTypes.RegisterEvidences(legacy.Amino)
}

func (cdc *Codec) AminoCodec() *LegacyAmino {
	return cdc.legacyCdc
}

func (cdc *Codec) ProtoCodec() *ProtoCodec {
	return cdc.protoCdc
}

func (cdc *Codec) SetAfterUpgradeMod(value bool) {
	if cdc.afterUpgradeMod != value {
		cdc.afterUpgradeMod = value
	}
}

func (cdc *Codec) IsAfterUpgrade() bool {
	return cdc.afterUpgradeMod
}
