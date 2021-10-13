package heightcache

import (
	"errors"
	"github.com/pokt-network/pocket-core/store/types"
	"math"
	"sort"
)

var _ types.SingleStoreCache = &MemoryCache{}

type MemoryCache struct {
	capacity    int64
	pastHeights []*StoreAtHeight
	current     *StoreAtHeight
}

func NewMemoryCache(size int64) *MemoryCache {
	storesAtHeight := make([]*StoreAtHeight, size)
	for i := range storesAtHeight {
		storesAtHeight[i] = NewStoreAtHeight()
	}
	return &MemoryCache{
		capacity:    size,
		current:     NewStoreAtHeight(),
		pastHeights: storesAtHeight,
	}
}

func (m *MemoryCache) InitializeStoreCache(height int64) error {
	m.current.height = height
	if m.current.data == nil {
		m.current.data = map[string]string{}
	}
	return nil
}

func (m MemoryCache) isHeightSafeToRead(height int64) bool {
	if height != m.current.height && height > m.current.height-(1+m.capacity) {
		for _, v := range m.pastHeights {
			if v.height == height {
				return true
			}
		}
	}
	return false

}

func (m MemoryCache) Get(height int64, key []byte) ([]byte, error) {
	if m.isHeightSafeToRead(height) {
		for i := range m.pastHeights {
			if m.pastHeights[i].height == height {
				return []byte(m.pastHeights[i].data[string(key)]), nil
			}
		}
	}
	return nil, errors.New("invalid height for get")
}

func (m MemoryCache) Has(height int64, key []byte) (bool, error) {
	return false, errors.New("not implemented")
}

func (m MemoryCache) Set(key []byte, value []byte) {
	m.current.data[string(key)] = string(value)
}

func (m MemoryCache) Remove(key []byte) error {
	delete(m.current.data, string(key))
	return nil
}

func (m MemoryCache) Iterator(height int64, start, end []byte) (types.Iterator, error) {
	if m.isHeightSafeToRead(height) {
		for _, v := range m.pastHeights {
			if v.height == height {
				return NewMemoryHeightIterator(v.data, string(start), string(end), v.orderedKeys, true), nil
			}
		}
	}
	return nil, errors.New("invalid height for iterator")
}

func (m MemoryCache) ReverseIterator(height int64, start, end []byte) (types.Iterator, error) {
	if m.isHeightSafeToRead(height) {
		for _, v := range m.pastHeights {
			if v.height == height {
				return NewMemoryHeightIterator(v.data, string(start), string(end), v.orderedKeys, false), nil
			}
		}
	}
	return nil, errors.New("invalid height for iterator")
}

func (m MemoryCache) Commit(height int64) {
	lowestHeight := int64(math.MaxInt64)
	lowestIdx := -1

	for idx, v := range m.pastHeights {
		if v.height < lowestHeight {
			lowestHeight = v.height
			lowestIdx = idx
		}
	}
	m.current.height = height

	m.pastHeights[lowestIdx].height = m.current.height
	m.pastHeights[lowestIdx].data = map[string]string{}

	orderedKeys := make([]string, len(m.current.data))
	for k, v := range m.current.data {
		m.pastHeights[lowestIdx].data[k] = v
		orderedKeys = append(orderedKeys, k)
	}
	sort.Strings(orderedKeys)
	m.pastHeights[lowestIdx].orderedKeys = orderedKeys
}

func (m MemoryCache) Initialize(currentData map[string]string, version int64) {
	m.current.data = currentData
	m.current.height = version
}

func (m MemoryCache) IsValid() bool {
	return true
}
