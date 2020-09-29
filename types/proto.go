package types

import (
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
)

// CustomProtobufType defines the interface custom gogo proto types must implement
// in order to be used as a "customtype" extension.
//
// ref: https://github.com/gogo/protobuf/blob/master/custom_types.md
type CustomProtobufType interface {
	Marshal() ([]byte, error)
	MarshalTo(data []byte) (n int, err error)
	Unmarshal(data []byte) error
	Size() int

	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

// primitive wrappers for proto below

var (
	_ codec.ProtoMarshaler = (*Bool)(nil)
	_ codec.ProtoMarshaler = (*Int64)(nil)
)

type Bool bool

func (b *Bool) Reset() {
	*b = false
}

func (b Bool) String() string {
	return fmt.Sprintf("%v", bool(b))
}

func (b Bool) ProtoMessage() {
	p := b.ToProto()
	p.ProtoMessage()
}

func (b Bool) Marshal() ([]byte, error) {
	p := b.ToProto()
	return p.Marshal()
}

func (b Bool) MarshalTo(data []byte) (n int, err error) {
	p := b.ToProto()
	return p.MarshalTo(data)
}

func (b Bool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	p := b.ToProto()
	return p.MarshalToSizedBuffer(dAtA)
}

func (b Bool) Size() int {
	p := b.ToProto()
	return p.Size()
}

func (b *Bool) Unmarshal(data []byte) (err error) {
	var pb ProtoBool
	err = pb.Unmarshal(data)
	*b = Bool(pb.B)
	return
}

func (b Bool) ToProto() ProtoBool {
	return ProtoBool{B: bool(b)}
}

func (pb ProtoBool) FromProto() Bool {
	return Bool(pb.B)
}

type Int64 int64

func (i *Int64) Reset() {
	*i = 0
}

func (i Int64) String() string {
	return fmt.Sprintf("%d", int64(i))
}

func (i Int64) ProtoMessage() {
	p := i.ToProto()
	p.ProtoMessage()
}

func (i Int64) Marshal() ([]byte, error) {
	p := i.ToProto()
	return p.Marshal()
}

func (i Int64) MarshalTo(data []byte) (n int, err error) {
	p := i.ToProto()
	return p.MarshalTo(data)
}

func (i Int64) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	p := i.ToProto()
	return p.MarshalToSizedBuffer(dAtA)
}

func (i Int64) Size() int {
	p := i.ToProto()
	return p.Size()
}

func (i *Int64) Unmarshal(data []byte) (err error) {
	var pi ProtoInt64
	err = pi.Unmarshal(data)
	*i = pi.FromProto()
	return
}

func (i Int64) ToProto() ProtoInt64 {
	return ProtoInt64{I: int64(i)}
}

func (pi ProtoInt64) FromProto() Int64 {
	return Int64(pi.I)
}
