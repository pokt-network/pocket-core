package iavl

import (
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/codec/types"
)

var cdc = codec.NewCodec(types.NewInterfaceRegistry())

// Options define tree options.
type Options struct {
	Sync bool
}

// DefaultOptions returns the default options for IAVL.
func DefaultOptions() *Options {
	return &Options{
		Sync: false,
	}
}

var (
	debugging = false
)

func debug(format string, args ...interface{}) {
	if debugging {
		fmt.Printf(format, args...)
	}
}

func maxInt8(a, b int8) int8 {
	if a > b {
		return a
	}
	return b
}

func cp(bz []byte) (ret []byte) {
	ret = make([]byte, len(bz))
	copy(ret, bz)
	return ret
}
