package cachemulti

import (
	"github.com/pokt-network/pocket-core/store/types"
)

//----------------------------------------
// CacheMultiStore

// CacheMultiStore holds many cache-wrapped stores.
// Implements MultiStore.
// NOTE: a CacheMultiStore (and MultiStores in general) should never expose the
// keys for the substores.
type CacheMultiStore struct {
	stores map[types.StoreKey]types.CacheWrap
}

var _ types.CacheMultiStore = CacheMultiStore{}

func NewMultiCache(stores map[types.StoreKey]types.CommitStore) CacheMultiStore {
	cms := CacheMultiStore{
		stores: make(map[types.StoreKey]types.CacheWrap, len(stores)),
	}

	for key, store := range stores {
		cms.stores[key] = store.CacheWrap()
	}
	return cms
}

// GetStoreType returns the type of the store.
func (cms CacheMultiStore) GetStoreType() types.StoreType {
	return types.StoreTypeMulti
}

// Write calls Write on each underlying store.
func (cms CacheMultiStore) Write() {
	for _, store := range cms.stores {
		store.Write()
	}
}

// Implements CacheWrapper.
func (cms CacheMultiStore) CacheWrap() types.CacheWrap {
	panic("CacheWrap(): can't double cache-wrap an already cachewrapped multistore")
}

// Implements MultiStore.
func (cms CacheMultiStore) CacheMultiStore() types.CacheMultiStore {
	panic("CacheMultiStore(): can't double cache-wrap an already cachewrapped multistore")
}

// GetStore returns an underlying CacheMultiStore by key.
func (cms CacheMultiStore) GetStore(key types.StoreKey) types.Store {
	return cms.stores[key].(types.Store)
}

// GetKVStore returns an underlying KVStore by key.
func (cms CacheMultiStore) GetKVStore(key types.StoreKey) types.KVStore {
	return cms.stores[key].(types.KVStore)
}
