package iavl

import (
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/codec/types"
)

var cdc = codec.NewCodec(types.NewInterfaceRegistry())
