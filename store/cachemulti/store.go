package cachemulti

import (
	"bytes"
	"fmt"
	"github.com/pokt-network/pocket-core/store/types"
	"sort"
	"sync"
)

var _ types.CacheKVStore = (*StoreCache)(nil)

type Operation int64

const (
	Delete Operation = iota
	Set
)

type CacheObject struct {
	key       []byte
	value     []byte
	operation Operation
}

// Store wraps an in-memory cache around an underlying types.KVStore.
type StoreCache struct {
	mtx           sync.Mutex
	unsortedCache map[string]CacheObject // used for reading
	sortedCache   [][]byte               // used for writing and iterating
	parent        types.KVStore          // the parent this store is caching
}

func NewStoreCache(parent types.KVStore) *StoreCache {
	return &StoreCache{
		mtx:           sync.Mutex{},
		unsortedCache: make(map[string]CacheObject),
		sortedCache:   make([][]byte, 0),
		parent:        parent,
	}
}

func (i *StoreCache) Get(key []byte) ([]byte, error) {
	if cacheObj, ok := i.unsortedCache[string(key)]; ok {
		if cacheObj.operation == Delete {
			return nil, nil
		} else {
			return cacheObj.value, nil
		}
	}
	return i.parent.Get(key)
}

func (i *StoreCache) Has(key []byte) (bool, error) {
	if cacheObj, ok := i.unsortedCache[string(key)]; ok {
		if cacheObj.operation == Delete {
			return false, nil
		} else {
			return true, nil
		}
	}
	return i.parent.Has(key)
}

func (i *StoreCache) Set(key, value []byte) error {
	_, found := i.unsortedCache[string(key)]
	i.unsortedCache[string(key)] = CacheObject{
		key:       key,
		value:     value,
		operation: Set,
	}
	// TODO naive
	if found {
		fmt.Println("Setting from found storeCache")
		index := sort.Search(len(i.sortedCache), func(a int) bool { return bytes.Equal(i.sortedCache[a], key) })
		i.sortedCache[index] = key
	} else {
		i.sortedCache = append(i.sortedCache, key)
		sort.Slice(i.sortedCache, func(x, y int) bool {
			return bytes.Compare(i.sortedCache[x], i.sortedCache[y]) < 0
		})
	}
	return nil
}

func (i *StoreCache) Delete(key []byte) error {
	// let's see if it found in the unsorted cache
	_, found := i.unsortedCache[string(key)]
	i.unsortedCache[string(key)] = CacheObject{
		key:       key,
		value:     nil,
		operation: Delete,
	}
	// TODO naive
	// if found, let's delete it from the sorted cache too
	if found {
		fmt.Println("deleting from found storeCache")
		l := len(i.sortedCache)
		index := sort.Search(l, func(a int) bool { return bytes.Equal(i.sortedCache[a], key) })
		if index+1 == l {
			i.sortedCache = i.sortedCache[:index]
		} else {
			i.sortedCache = append(i.sortedCache[:index], i.sortedCache[index+1:]...)
		}
	}
	return nil
}

func (i *StoreCache) Iterator(start, end []byte) (types.Iterator, error) {
	fmt.Println("iterating on StoreCache")
	parent, err := i.parent.Iterator(start, end)
	if err != nil {
		return nil, err
	}
	return NewCacheMergeIterator(parent, i.sortedCache, i.unsortedCache, false), nil
}

func (i *StoreCache) ReverseIterator(start, end []byte) (types.Iterator, error) {
	fmt.Println("reverse-iterating on StoreCache")
	parent, err := i.parent.ReverseIterator(start, end)
	if err != nil {
		return nil, err
	}
	return NewCacheMergeIterator(parent, i.sortedCache, i.unsortedCache, false), nil
}

func (i *StoreCache) Write() {
	for _, k := range i.sortedCache {
		co, _ := i.unsortedCache[string(k)]
		switch co.operation {
		case Set:
			_ = i.parent.Set(co.key, co.value)
		case Delete:
			_ = i.parent.Delete(co.key)
		}
	}
	// Clear the cache
	i.unsortedCache = make(map[string]CacheObject)
	i.sortedCache = make([][]byte, 0)
}

func (i *StoreCache) GetStoreType() types.StoreType {
	panic("GetStoreType not implemented for StoreCache")
}

func (i *StoreCache) CacheWrap() types.CacheWrap {
	panic("CacheWrap not implemented for StoreCache")
}
