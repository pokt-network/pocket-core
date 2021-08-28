package heightcache

import (
	"errors"
	"github.com/pokt-network/pocket-core/store/types"
)

var _ types.SingleStoreCache = &MemoryCache{}

type MemoryCache struct {
	capacity    int64
	pastHeights []*StoreAtHeight
	current     *StoreAtHeight
}

func NewStoreMemoryCache(size int64) *MemoryCache {
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

func (c *MemoryCache) InitializeStoreCache(height int64) error {
	c.current.height = height
	if c.current.data == nil {
		c.current.data = map[string]string{}
	}
	return nil
}

func (m MemoryCache) Get(height int64, key []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (m MemoryCache) Has(height int64, key []byte) (bool, error) {
	return false, errors.New("not implemented")
}

func (m MemoryCache) Set(key []byte, value []byte) {
}

func (m MemoryCache) Remove(key []byte) error {
	return errors.New("not implemented")
}

func (m MemoryCache) Iterator(height int64, start, end []byte) (types.Iterator, error) {
	return nil, errors.New("not implemented")
}

func (m MemoryCache) ReverseIterator(height int64, start, end []byte) (types.Iterator, error) {
	return nil, errors.New("not implemented")
}

func (m MemoryCache) Commit(height int64) {
	panic("implement me")
}

func (m MemoryCache) Initialize(currentData map[string]string, version int64) {
	panic("implement me")
}
