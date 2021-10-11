package heightcache

import "github.com/pokt-network/pocket-core/store/types"

var _ types.MultiStoreCache = &MultiStoreInvalidCache{}

type MultiStoreInvalidCache struct {
	invalidCache *InvalidCache
}

func NewMultiStoreInvalidCache() types.MultiStoreCache {
	return &MultiStoreInvalidCache{invalidCache: &InvalidCache{}}
}

func (m MultiStoreInvalidCache) InitializeSingleStoreCache(height int64, storeKey types.StoreKey) error {
	return nil
}

func (m MultiStoreInvalidCache) GetSingleStoreCache(storekey types.StoreKey) types.SingleStoreCache {
	return &InvalidCache{}
}
