package cachemulti

import (
	"github.com/pokt-network/pocket-core/store/types"
)

var _ types.CacheMultiStore = (*CacheMultiStore)(nil)

type CacheMultiStore struct {
	cacheStores map[types.StoreKey]types.CacheWrap
}

func NewCacheMulti(stores map[types.StoreKey]types.CommitStore) types.CacheMultiStore {
	newStores := make(map[types.StoreKey]types.CacheWrap)
	for k, s := range stores {
		newStores[k] = s.CacheWrap()
	}
	return &CacheMultiStore{newStores}
}

func (s *CacheMultiStore) GetStoreType() types.StoreType {
	return types.StoreTypeMulti
}

func (s *CacheMultiStore) Write() {
	for _, store := range s.cacheStores {
		store.Write()
	}
}

func (s *CacheMultiStore) GetStore(key types.StoreKey) types.Store {
	return s.cacheStores[key].(*StoreCache)
}

func (s *CacheMultiStore) GetKVStore(key types.StoreKey) types.KVStore {
	return s.cacheStores[key].(*StoreCache)
}

func (s *CacheMultiStore) CacheWrap() types.CacheWrap {
	panic("CacheWrap(): can't double cache-wrap an already cachwrapped multistore")
}

func (s *CacheMultiStore) CacheMultiStore() types.CacheMultiStore {
	panic("CacheMultiStore(): can't double cache-wrap an already cachwrapped multistore")
}
