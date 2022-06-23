package slim

import (
	"github.com/pokt-network/pocket-core/store/cachemulti"
	"github.com/pokt-network/pocket-core/store/iavl"
	"github.com/pokt-network/pocket-core/store/slim/dedup"
	"github.com/pokt-network/pocket-core/store/slim/memdb"
	"github.com/pokt-network/pocket-core/store/types"
	db "github.com/tendermint/tm-db"
)

type Store struct {
	Cache     *dedup.Store
	Dedup     *dedup.Store
	IAVLStore iavl.Store
}

func NewStoreWithIAVL(d *db.GoLevelDB, cDB *memdb.PocketMemDB, height int64, prefix string, commitID types.CommitID) *Store {
	iavlStore, err := iavl.NewStore(db.NewPrefixDB(d, []byte("s/k:"+prefix+"/")), commitID)
	if err != nil {
		panic("iavl store failed to load for height: %s prefix: %s")
	}
	return &Store{
		Cache:     dedup.NewStore(height, prefix, cDB, true),
		Dedup:     dedup.NewStore(height, prefix, d, false),
		IAVLStore: *iavlStore,
	}
}

func NewStoreWithoutIAVL(db *db.GoLevelDB, cDB *memdb.PocketMemDB, latestHeight, height int64, prefix string) *Store {
	var cache *dedup.Store
	// if cache contains this height
	if height > latestHeight-dedup.DefaultCacheKeepHeights+1 {
		cache = dedup.NewStore(height, prefix, cDB, true)
	}
	return &Store{
		Cache: cache,
		Dedup: dedup.NewStore(height, prefix, db, false),
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

// writes (both stores)

func (s *Store) Set(key, value []byte) error {
	if err := s.Cache.Set(key, value); err != nil {
		return err
	}
	//if err := s.Dedup.Set(key, value); err != nil {
	//	return err
	//}
	return s.IAVLStore.Set(key, value)
}

func (s *Store) Delete(key []byte) error {
	if err := s.Cache.Delete(key); err != nil {
		return err
	}
	//if err := s.Dedup.Delete(key); err != nil {
	//	return err
	//}
	return s.IAVLStore.Delete(key)
}

// lifecycle operations (special)

func (s *Store) CommitBatch(b db.Batch) (types.CommitID, db.Batch) {
	// commit both stores, but only return commitID from IAVL
	_, _ = s.Cache.CommitBatch(b, s.Dedup)
	//_, b = s.Dedup.CommitBatch(b)
	return s.IAVLStore.CommitBatch(b)
}

func (s *Store) PreloadCache(latestHeight int64) { s.Dedup.PreloadCache(latestHeight, s.Cache) }
func (s *Store) CacheWrap() types.CacheWrap      { return cachemulti.NewStoreCache(s) }
func (s *Store) LastCommitID() types.CommitID    { return s.IAVLStore.LastCommitID() }

// unused below

func (s *Store) Commit() types.CommitID {
	panic("Commit() called in store; should use commitBatch for atomic safety")
}
