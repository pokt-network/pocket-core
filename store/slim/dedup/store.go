package dedup

import (
	"bytes"
	"fmt"
	"github.com/pokt-network/pocket-core/store/types"
	db "github.com/tendermint/tm-db"
)

// Dedup store consists of two different spaces:
//
// DATASTORE that holds the actual bytes of the data
// LINKSTORE that holds the link or alias to the data
//
// This design allows for DATA space to only be affected during writes
// while allowing for the LINKSTORE to keep track of height based states
// for historical queries. For example, the link space will have a key for
// every single height / item, but the value is just the DATASTORE key,
// whereas the DATASTORE will only have key/values for the height/item combinations
// where the item actually had a state change (was written).

// The first design should be simple:
// DATASTORE: KEY: <Hash> -> VALUE: <data-bytes>
// LINKSTORE: KEY: /link/<height>/<key>/ -> VALUE: <Hash>

// Example:
// Height 1
// <SomeHash1> -> <validatorProtoBytes>    | /link/height1/validator/addr1 -> <SomeHash1>
// Height 2
// <noStateChange>                         | /link/height2/validator/addr1 -> <SomeHash1>
// Height 3
// <SomeHash2> -> <newValidatorProtoBytes> | /link/height3/validator/addr1 -> <SomeHash2>

type Store struct {
	Height   int64
	Prefix   string
	ParentDB db.DB
	isCache  bool
}

func NewStore(height int64, prefix string, parent db.DB, isCache bool) *Store {
	return &Store{
		Height:   height,
		Prefix:   prefix,
		ParentDB: parent,
		isCache:  isCache,
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
	//if err := s.trackOrphan(linkStoreKey, Set); err != nil {
	//	return err
	//}
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
	//if err := s.trackOrphan(linkStoreKey, Del); err != nil {
	//	return err
	//}
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

// PreloadCache : Preload the last N heights into the 'cache' so the most recent lookups are not off of disk
func (s *Store) PreloadCache(latestHeight int64, cache *Store) {
	fmt.Printf("Preloading cache for %s\n", s.Prefix)
	for height := getPreloadStartHeight(latestHeight); height <= latestHeight; height++ {
		it := NewDedupIteratorForHeight(s.ParentDB, height, s.Prefix)
		cache.Height = height
		for ; it.Valid(); it.Next() {
			_ = cache.Set(it.Key(), it.Value())
		}
		cache.Height = latestHeight
		it.Close()
	}
	return
}

// CommitCache : writes the cache height to batch, prunes a height from cache, and then prepares the next cache height
func (s *Store) CommitCache(b db.Batch, dedupStore *Store) {
	// copy the cache to disk
	//s.CopyHeight(s.Height, dedupStore, b)
	// increment heights before the next logic
	//s.Height++
	dedupStore.Height++
	// prune the cache
	//s.PruneCache(s.Height - DefaultCacheKeepHeights)
	// prepare the current height by copying all of the keys
	//s.PrepareNextHeight(nil)
	dedupStore.PrepareNextHeight(b)
	return
}

// CopyHeight : copy the entire height from the cache to disk using the batch
func (s *Store) CopyHeight(height int64, dedupStore *Store, b db.Batch) {
	it := NewDedupIteratorForHeight(s.ParentDB, height, s.Prefix)
	for ; it.Valid(); it.Next() {
		// full linkstore key with height
		linkKey := it.it.Key()
		// full datastore (hash) key
		dataKey := it.it.Value()
		// batches don't have the parent.Set() logic
		b.Set(linkKey, dataKey)
		// set the data if the data (hash) key isn't already set
		if found, _ := dedupStore.ParentDB.Has(dataKey); !found {
			dataValue, _ := s.ParentDB.Get(dataKey)
			b.Set(dataKey, dataValue)
		}
	}
	it.Close()
}

// PruneCache : delete the height and orphans from cache
func (s *Store) PruneCache(height int64) {
	if height < 0 {
		height = 0
	}
	keysToDelete := make([][]byte, 0)
	it := NewDedupIteratorForHeight(s.ParentDB, height, s.Prefix)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		linkKey := it.it.Key()
		keysToDelete = append(keysToDelete, linkKey)
	}
	oIt := NewOrphanIteratorForHeight(s.ParentDB, height, s.Prefix)
	defer oIt.Close()
	for ; oIt.Valid(); oIt.Next() {
		orphanKey := oIt.it.Key()
		hashKey := oIt.Key() // key for actual data
		keysToDelete = append(keysToDelete, orphanKey)
		keysToDelete = append(keysToDelete, hashKey)
	}
	for _, k := range keysToDelete {
		if err := s.ParentDB.Delete(k); err != nil {
			panic("an error occurred deleting in cache")
		}
	}
	return
}

// An orphan in this context is a data K/V that has no associated link to it. This logic is only applicable
// during pruning (which is a cache only operation at this point). In order to not prematurely delete the
// datakey before it's needed in a valid (not pruned) historical lookup, it must be tracked under a specific
// key set. This enables the deletion of all orphans at the appropriate height (prune height).
//
// In this implementation, we store orphan keys at DELETED_HEIGHT+1 OR SET_HEIGHT+1 (for overwrites) which is exactly
// when it will be removed during the (delete/prune height).

func (s *Store) trackOrphan(linkStoreKey []byte, op OperationType) error {
	// only track orphan if is cached store because we only prune on cache store
	if s.isCache {
		// attempt to delete any previous orphan key for this height to prevent a sameHeight set-del-set condition
		if op == Set {
			_ = s.ParentDB.Delete(OrphanKey(s.Height, s.Prefix, HashKey(linkStoreKey)))
		}
		dataKey, _ := s.ParentDB.Get(linkStoreKey)
		if dataKey != nil {
			if (op == Set && !bytes.Equal(dataKey, HashKey(linkStoreKey))) || op == Del {
				orphanHeightKey := OrphanKey(s.Height, s.Prefix, dataKey)
				if err := s.ParentDB.Set(orphanHeightKey, nil); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// unused below

func (s *Store) Commit() types.CommitID       { panic("Commit() called in de-dup store ") }
func (s *Store) CacheWrap() types.CacheWrap   { panic("cachewrap not implemented for de-dup store") }
func (s *Store) LastCommitID() types.CommitID { panic("lastCommitID not implemented for de-dup store") }
