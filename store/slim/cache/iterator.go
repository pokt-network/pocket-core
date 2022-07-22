package cache

import (
	"github.com/pokt-network/pocket-core/store/types"
)

var _ types.Iterator = &Iterator{}

type Iterator struct {
	parent    *CacheDB
	height    int64
	prefix    string
	startKey  string
	endKey    string
	index     int
	lastKey   string
	key       string
	value     []byte
	valid     bool
	isReverse bool
}

func NewIterator(parent *CacheDB, height int64, prefix string, startKey, endKey []byte, isReverse bool) *Iterator {
	start := PrefixKey(prefix, startKey)
	end := PrefixKey(prefix, endKey)
	if endKey == nil {
		end = string(types.PrefixEndBytes([]byte(prefix)))
	}
	i := &Iterator{
		parent:    parent,
		height:    height,
		prefix:    prefix,
		startKey:  start,
		endKey:    end,
		isReverse: isReverse,
		valid:     true,
	}
	i.Next()
	return i
}

func (i *Iterator) Next() {
	i.parent.L.RLock()
	defer i.parent.L.RUnlock()
	kvSlice := i.parent.M[i.height]
	if kvSlice.size == 0 {
		i.valid = false
		return
	}
	if i.index >= kvSlice.size {
		i.index = kvSlice.size - 1
	}
	// first seek
	if i.key == "" {
		if i.isReverse {
			i.index, _ = kvSlice.Search(i.endKey)
			i.index--
		} else {
			i.index, _ = kvSlice.Search(i.startKey)
		}
		// if the set is manipulated during iteration
		// note that this matches the goleveldb behavior
		// of changing the set during iteration
	} else if i.key != kvSlice.S.KVSlice[i.index].Key {
		var found bool
		i.index, found = kvSlice.Search(i.key)
		if i.isReverse {
			i.index--
		} else if found {
			i.index++
		}
	} else { // standard next
		if i.isReverse {
			i.index--
		} else {
			i.index++
		}
	}
	if i.index < 0 || i.index >= kvSlice.size {
		i.valid = false
		return
	}
	keyValue := kvSlice.S.KVSlice[i.index]
	if keyValue.Key < i.startKey || i.endKey <= keyValue.Key {
		i.valid = false
		return
	}
	i.lastKey = i.key
	i.key = keyValue.Key
	i.value = keyValue.Value
}

func (i *Iterator) Key() (key []byte)                  { return RemovePrefix(i.prefix, i.key) }
func (i *Iterator) Value() (value []byte)              { return i.value }
func (i *Iterator) Valid() bool                        { return i.valid }
func (i *Iterator) Error() error                       { panic("error not implemented in cache iterator") }
func (i *Iterator) Close()                             { /* no op */ }
func (i *Iterator) Domain() (start []byte, end []byte) { return []byte(i.startKey), []byte(i.endKey) }
func PrefixKey(prefix string, key []byte) string       { return string(append([]byte(prefix), key...)) }
func RemovePrefix(prefix string, key string) []byte    { return []byte(key[len(prefix):]) }
