package store

import (
	"bytes"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/pokt-network/pocket-core/store/types"
)

// Gets the first item.
func First(st KVStore, start, end []byte) (kp kv.Pair, ok bool) {
	iter, _ := st.Iterator(start, end)
	if !iter.Valid() {
		return kp, false
	}
	defer iter.Close()

	return kv.Pair{Key: iter.Key(), Value: iter.Value()}, true
}

// Gets the last item.  `end` is exclusive.
func Last(st KVStore, start, end []byte) (kp kv.Pair, ok bool) {
	iter, _ := st.ReverseIterator(end, start)
	if !iter.Valid() {
		if v, _ := st.Get(start); v != nil {
			return kv.Pair{Key: types.Cp(start), Value: types.Cp(v)}, true
		}
		return kp, false
	}
	defer iter.Close()

	if bytes.Equal(iter.Key(), end) {
		// Skip this one, end is exclusive.
		iter.Next()
		if !iter.Valid() {
			return kp, false
		}
	}

	return kv.Pair{Key: iter.Key(), Value: iter.Value()}, true
}
