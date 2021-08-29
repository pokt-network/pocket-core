package heightcache

import (
	"github.com/pokt-network/pocket-core/store/types"
)

var _ types.MultiStoreCache = &MultiStoreMemoryCache{}

type MultiStoreMemoryCache struct {
	capacity int64
	stores   map[types.StoreKey]*MemoryCache
}

func NewMultiStoreMemoryCache(capacity int64) types.MultiStoreCache {
	stores := make(map[types.StoreKey]*MemoryCache)
	return &MultiStoreMemoryCache{
		capacity: capacity,
		stores:   stores,
	}
}

func (m MultiStoreMemoryCache) InitializeSingleStoreCache(height int64, storeKey types.StoreKey) error {
	if m.stores[storeKey] == nil {
		m.stores[storeKey] = NewMemoryCache(m.capacity)
	}
	return m.stores[storeKey].InitializeStoreCache(height)
}

func (m MultiStoreMemoryCache) GetSingleStoreCache(storeKey types.StoreKey) types.SingleStoreCache {
	if m.stores[storeKey] == nil {
		m.InitializeSingleStoreCache(-1, storeKey)
	}
	return m.stores[storeKey]
}
