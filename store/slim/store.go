package slim

import (
	"github.com/pokt-network/pocket-core/store/cachemulti"
	"github.com/pokt-network/pocket-core/store/iavl"
	"github.com/pokt-network/pocket-core/store/slim/cache"
	"github.com/pokt-network/pocket-core/store/slim/dedup"
	"github.com/pokt-network/pocket-core/store/types"
	db "github.com/tendermint/tm-db"
)

const (
	maxCacheKeepHeights = 25
)

type Store struct {
	Cache     *cache.CacheStore // optimized query store (only last N heights)
	Dedup     *dedup.Store      // query store
	IAVLStore iavl.Store        // state commitment store
}

func NewStoreWithIAVL(persisted *db.GoLevelDB, cacheDB *cache.CacheDB, height int64, prefix string, commitID types.CommitID) *Store {
	iavlStore, _ := iavl.NewStore(db.NewPrefixDB(persisted, []byte("s/k:"+prefix+"/")), commitID)
	return &Store{
		Dedup:     dedup.NewStore(height, prefix, persisted),
		Cache:     cache.NewStore(height, prefix, cacheDB),
		IAVLStore: *iavlStore,
	}
}

func NewStoreWithoutIAVL(persisted *db.GoLevelDB, cacheDB *cache.CacheDB, latestHeight, height int64, prefix string) *Store {
	if height > latestHeight-maxCacheKeepHeights {
		return &Store{
			Cache: cache.NewStore(height, prefix, cacheDB),
			Dedup: dedup.NewStore(height, prefix, persisted),
		}
	}
	return &Store{
		Dedup: dedup.NewStore(height, prefix, persisted),
	}
}

// reads (de-dup only)

func (s *Store) Get(key []byte) ([]byte, error) {
	if s.Cache != nil {
		return s.Cache.Get(key)
	}
	return s.Dedup.Get(key)
}
func (s *Store) Has(key []byte) (bool, error) {
	if s.Cache != nil {
		return s.Cache.Has(key)
	}
	return s.Dedup.Has(key)
}
func (s *Store) Iterator(start, end []byte) (types.Iterator, error) {
	if s.Cache != nil {
		return s.Cache.Iterator(start, end)
	}
	return s.Dedup.Iterator(start, end)
}
func (s *Store) ReverseIterator(start, end []byte) (types.Iterator, error) {
	if s.Cache != nil {
		return s.Cache.ReverseIterator(start, end)
	}
	return s.Dedup.ReverseIterator(start, end)
}

// writes (all stores)

func (s *Store) Set(key, value []byte) error {
	s.Cache.Set(key, value)
	if err := s.Dedup.Set(key, value); err != nil {
		return err
	}
	return s.IAVLStore.Set(key, value)
}

func (s *Store) Delete(key []byte) error {
	s.Cache.Delete(key)
	if err := s.Dedup.Delete(key); err != nil {
		return err
	}
	return s.IAVLStore.Delete(key)
}

// lifecycle operations (special)

func (s *Store) Commit() types.CommitID {
	s.Cache.Commit()
	s.Dedup.Commit()
	return s.IAVLStore.CommitIAVL()
}

func (s *Store) CacheWrap() types.CacheWrap   { return cachemulti.NewStoreCache(s) }
func (s *Store) LastCommitID() types.CommitID { return s.IAVLStore.LastCommitID() }
