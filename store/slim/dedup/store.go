package dedup

import (
	"github.com/pokt-network/pocket-core/store/types"
	db "github.com/tendermint/tm-db"
)

type Store struct {
	Height   int64
	Prefix   string
	ParentDB db.DB
}

func NewStore(height int64, prefix string, parent db.DB) *Store {
	return &Store{
		Height:   height,
		Prefix:   prefix,
		ParentDB: parent,
	}
}

// reads

func (s *Store) Get(k []byte) ([]byte, error) {
	linkStoreKey := HeightKey(s.Height, s.Prefix, k)
	dataStoreKey, err := s.ParentDB.Get(linkStoreKey)
	if err != nil {
		return nil, err
	}
	return s.ParentDB.Get(dataStoreKey)
}

func (s *Store) Has(k []byte) (bool, error) {
	return s.ParentDB.Has(HeightKey(s.Height, s.Prefix, k))
}

func (s *Store) Iterator(start, end []byte) (types.Iterator, error) {
	return NewDedupIterator(s.ParentDB, s.Height, s.Prefix, start, end, false)
}

func (s *Store) ReverseIterator(start, end []byte) (types.Iterator, error) {
	return NewDedupIterator(s.ParentDB, s.Height, s.Prefix, start, end, true)
}

// writes

func (s *Store) Set(k, value []byte) error {
	linkStoreKey := HeightKey(s.Height, s.Prefix, k)
	dataStoreKey := HashKey(linkStoreKey)
	if err := s.ParentDB.Set(linkStoreKey, dataStoreKey); err != nil {
		return err
	}
	if err := s.ParentDB.Set(dataStoreKey, value); err != nil {
		return err
	}
	return nil
}

func (s *Store) Delete(k []byte) error {
	linkStoreKey := HeightKey(s.Height, s.Prefix, k)
	return s.ParentDB.Delete(linkStoreKey)
}

// lifecycle ops

// PrepareNextHeight : load the current height with the endState of the current-1 height
func (s *Store) PrepareNextHeight(b db.Batch) {
	it := NewDedupIteratorForHeight(s.ParentDB, s.Height-1, s.Prefix)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		nextHeightKey := HeightKey(s.Height, s.Prefix, it.Key())
		linkValue := it.it.Value()
		if b == nil {
			_ = s.ParentDB.Set(nextHeightKey, linkValue)
		} else {
			b.Set(nextHeightKey, linkValue)
		}
	}
}

func (s *Store) DeleteNextHeight() {
	it := NewDedupIteratorForHeight(s.ParentDB, s.Height, s.Prefix)
	keysToDelete := make([][]byte, 0)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		keysToDelete = append(keysToDelete, it.it.Key())
	}
	for _, k := range keysToDelete {
		_ = s.Delete(k)
	}
}

// CommitBatch : writes the cache height to batch, prunes a height from cache, and then prepares the next cache height
func (s *Store) CommitBatch(b db.Batch) {
	s.Height++
	// prepare the current height by copying all the keys
	s.PrepareNextHeight(b)
	return
}

// unused below

func (s *Store) Commit() types.CommitID       { panic("Commit() called in de-dup store ") }
func (s *Store) CacheWrap() types.CacheWrap   { panic("cachewrap not implemented for de-dup store") }
func (s *Store) LastCommitID() types.CommitID { panic("lastCommitID not implemented for de-dup store") }
