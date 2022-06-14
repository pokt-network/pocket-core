package dbadapter

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/pokt-network/pocket-core/store/cachekv"
	"github.com/pokt-network/pocket-core/store/types"
)

// Wrapper type for dbm.Db with implementation of KVStore
type Store struct {
	dbm.DB
}

// GetStoreType returns the type of the store.
func (Store) GetStoreType() types.StoreType {
	return types.StoreTypeDB
}

// CacheWrap cache wraps the underlying store.
func (dsa Store) CacheWrap() types.CacheWrap {
	return cachekv.NewStore(dsa)
}

// dbm.DB implements KVStore so we can CacheKVStore it.
var _ types.KVStore = Store{}
