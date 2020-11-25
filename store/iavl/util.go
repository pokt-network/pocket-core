package iavl

import "fmt"

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

// Returns a slice of the same length (big endian)
// except incremented by one.
// Appends 0x00 if bz is all 0xFF.
// CONTRACT: len(bz) > 0
func cpIncr(bz []byte) (ret []byte) {
	ret = cp(bz)
	for i := len(bz) - 1; i >= 0; i-- {
		if ret[i] < byte(0xFF) {
			ret[i]++
			return
		}
		ret[i] = byte(0x00)
		if i == 0 {
			return append(ret, 0x00)
		}
	}
	return []byte{0x00}
}

func cp(bz []byte) (ret []byte) {
	ret = make([]byte, len(bz))
	copy(ret, bz)
	return ret
}
