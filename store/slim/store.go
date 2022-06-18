package slim

import (
	"github.com/pokt-network/pocket-core/store/cachemulti"
	"github.com/pokt-network/pocket-core/store/iavl"
	"github.com/pokt-network/pocket-core/store/slim/dedup"
	"github.com/pokt-network/pocket-core/store/types"
	db "github.com/tendermint/tm-db"
)

type Store struct {
	Dedup     dedup.Store
	IAVLStore iavl.Store
}

func NewStoreWithIAVL(d *db.GoLevelDB, height int64, prefix string, commitID types.CommitID) *Store {
	iavlStore, err := iavl.NewStore(db.NewPrefixDB(d, []byte("s/k:"+prefix+"/")), commitID)
	if err != nil {
		panic("iavl store failed to load for height: %s prefix: %s")
	}
	return &Store{
		Dedup:     dedup.NewStore(height, prefix, *d),
		IAVLStore: *iavlStore,
	}
}

func NewStoreWithoutIAVL(db *db.GoLevelDB, height int64, prefix string) *Store {
	return &Store{
		Dedup: dedup.NewStore(height, prefix, *db),
	}
}

// reads (de-dup only)

func (s *Store) Get(key []byte) ([]byte, error) { return s.Dedup.Get(key) }
func (s *Store) Has(key []byte) (bool, error)   { return s.Dedup.Has(key) }
func (s *Store) Iterator(start, end []byte) (types.Iterator, error) {
	return s.Dedup.Iterator(start, end)
}
func (s *Store) ReverseIterator(start, end []byte) (types.Iterator, error) {
	return s.Dedup.ReverseIterator(start, end)
}

// writes (both stores)

func (s *Store) Set(key, value []byte) error {
	if err := s.Dedup.Set(key, value); err != nil {
		return err
	}
	return s.IAVLStore.Set(key, value)
}

func (s *Store) Delete(key []byte) error {
	if err := s.Dedup.Delete(key); err != nil {
		return err
	}
	return s.IAVLStore.Delete(key)
}

// lifecycle operations (special)

func (s *Store) CommitBatch(b db.Batch) (types.CommitID, db.Batch) {
	// commit both stores, but only return commitID from IAVL
	_, b = s.Dedup.CommitBatch(b)
	return s.IAVLStore.CommitBatch(b)
}

func (s *Store) CacheWrap() types.CacheWrap   { return cachemulti.NewStoreCache(s) }
func (s *Store) LastCommitID() types.CommitID { return s.IAVLStore.LastCommitID() }

// unused below

func (s *Store) Commit() types.CommitID {
	panic("Commit() called in store; should use commitBatch for atomic safety")
}
