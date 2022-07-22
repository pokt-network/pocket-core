package dedup

import (
	"bytes"
	"github.com/pokt-network/pocket-core/types"
	dbm "github.com/tendermint/tm-db"
)

var _ types.Iterator = &CascadeIterator{}
var _ types.KVStore = &Store{}

type CascadeIterator struct {
	parent             dbm.DB
	exclusiveEndHeight []byte
	it                 dbm.Iterator
	height             int64
	prefix             string
	value              []byte
}

func NewCascadeIterator(parent dbm.DB, height int64, prefix string, exclusiveEndHeight, startKey, endKey []byte, isReverse bool) (cascadeIterator *CascadeIterator, err error) {
	endPrefix := prefix
	// iterate over the 'exists' space inorder to have only 1 key per historical set
	start := KeyForExists(prefix, startKey)
	if endKey == nil {
		endPrefix = string(types.PrefixEndBytes([]byte(prefix)))
	}
	end := KeyForExists(endPrefix, endKey)
	cascadeIterator = &CascadeIterator{
		parent:             parent,
		exclusiveEndHeight: exclusiveEndHeight,
		height:             height,
		prefix:             prefix,
	}
	if isReverse {
		cascadeIterator.it, err = parent.ReverseIterator(start, end)
	} else {
		cascadeIterator.it, err = parent.Iterator(start, end)
	}
	if cascadeIterator.it.Valid() {
		cascadeIterator.Seek()
	}
	return
}

func (c *CascadeIterator) Next() {
	c.it.Next()
	c.Seek()
}

func (c *CascadeIterator) Seek() {
	// loop until we -> find a valid value OR get to the end of the iterator
	for {
		if !c.it.Valid() {
			return
		}
		latestKey, found := GetLatestHeightKeyRelativeToQuery(c.prefix, c.exclusiveEndHeight, c.Key(), c.parent)
		if !found { // it exists (we're iterating over exists space), so it must be a future height key
			c.it.Next()
			continue
		}
		c.value, _ = c.parent.Get(latestKey)
		// deleted value, so continue
		if bytes.Equal(c.value, DeletedValue) {
			c.it.Next()
			continue
		}
		return
	}
}

func (c *CascadeIterator) Value() (value []byte)              { return c.value }
func (c *CascadeIterator) Key() (key []byte)                  { return KeyFromExistsKey(c.prefix, c.it.Key()) }
func (c *CascadeIterator) Error() error                       { return c.it.Error() }
func (c *CascadeIterator) Close()                             { c.it.Close() }
func (c *CascadeIterator) Valid() bool                        { return c.it.Valid() }
func (c *CascadeIterator) Domain() (start []byte, end []byte) { return c.it.Domain() }
