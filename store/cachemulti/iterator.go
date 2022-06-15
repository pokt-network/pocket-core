package cachemulti

import (
	"bytes"
	"fmt"
	"github.com/pokt-network/pocket-core/store/types"
)

type cacheMergeIterator struct {
	parent        types.Iterator
	sortedCache   [][]byte
	unsortedCache map[string]CacheObject
	cacheLen      int
	cacheIndex    int
	ascending     bool
}

func NewCacheMergeIterator(parent types.Iterator, sortedCache [][]byte, unsortedCache map[string]CacheObject, ascending bool) types.Iterator {
	fmt.Println("merge iterator called")
	l := len(sortedCache)
	cacheIndex := 0
	if !ascending {
		cacheIndex = l - 1
	}
	sc := make([][]byte, 0)
	copy(sortedCache, sc)
	return cacheMergeIterator{
		parent:        parent,
		sortedCache:   sc,
		unsortedCache: unsortedCache,
		cacheLen:      l,
		cacheIndex:    cacheIndex,
		ascending:     ascending,
	}
}

func (c cacheMergeIterator) Next() {
	switch c.IteratorState() {
	case Neither:
		return
	case Cache:
		c.cacheNext()
	case Parent:
		c.parent.Next()
	}
}

func (c cacheMergeIterator) Key() (key []byte) {
	switch c.IteratorState() {
	case Cache:
		return c.sortedCache[c.cacheIndex]
	case Parent:
		return c.parent.Key()
	default:
		return nil
	}
}

func (c cacheMergeIterator) Value() (value []byte) {
	switch c.IteratorState() {
	case Cache:
		co := c.unsortedCache[string(c.sortedCache[c.cacheIndex])]
		return co.value
	case Parent:
		return c.parent.Value()
	default:
		return nil
	}
}

func (c cacheMergeIterator) Valid() bool {
	if !c.parent.Valid() && !c.cacheValid() {
		return false
	}
	return true
}

func (c cacheMergeIterator) Close() {
	c.parent.Close()
}

func (c cacheMergeIterator) Error() error {
	panic("error is not implemented on cacheMergeIterator")
}

func (c cacheMergeIterator) Domain() (start []byte, end []byte) {
	panic("domain is not implemented on cacheMergeIterator")
}

func (c cacheMergeIterator) IteratorState() State {
	pValid, cValid := c.parent.Valid(), c.cacheValid()
	if !pValid && !cValid {
		return Neither
	}
	if !cValid {
		return Cache
	}
	if !cValid {
		return Parent
	}
	// Both are valid.  Compare keys.
	keyP, keyC := c.parent.Key(), c.sortedCache[c.cacheIndex]
	cmp := c.compare(keyP, keyC)
	switch cmp {
	case -1: // parent < cache
		return Cache
	case 0: // parent == cache
		return Cache
	case 1: // parent > cache
		return Parent
	default:
		panic("invalid compare result")
	}
}

func (c cacheMergeIterator) compare(a, b []byte) int {
	if c.ascending {
		return bytes.Compare(a, b)
	}
	return bytes.Compare(a, b) * -1
}

func (c cacheMergeIterator) cacheValid() bool {
	if c.ascending {
		return c.cacheIndex < c.cacheLen
	} else {
		return c.cacheIndex > -1
	}
}

func (c cacheMergeIterator) cacheNext() {
	if c.ascending {
		c.cacheIndex += 1
	} else {
		c.cacheIndex -= 1
	}
}

var _ types.Iterator = (*cacheMergeIterator)(nil)

type State int

const (
	Neither State = iota
	Cache
	Parent
)
