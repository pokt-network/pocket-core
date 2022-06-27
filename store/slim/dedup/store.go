package dedup

import (
	"bytes"
	"fmt"
	"github.com/pokt-network/pocket-core/store/types"
	sdk "github.com/pokt-network/pocket-core/x/pocketcore/types"
	db "github.com/tendermint/tm-db"
	"strconv"
	"strings"
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
	val, err := s.ParentDB.Get(dataStoreKey)
	return val, err
}

func (s *Store) Has(k []byte) (bool, error) {
	linkStoreKey := HeightKey(s.Height, s.Prefix, k)
	return s.ParentDB.Has(linkStoreKey)
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
	if err := s.TrackOrphan(linkStoreKey); err != nil {
		return err
	}
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
	if err := s.TrackOrphan(linkStoreKey); err != nil {
		return err
	}
	return s.ParentDB.Delete(linkStoreKey)
}

// lifecycle ops

func (s *Store) CommitCache(b db.Batch, dedupStore *Store) {
	s.CopyHeight(s.Height, dedupStore, b)
	s.Height++
	dedupStore.Height++
	s.DeleteCacheHeight(s.Height - DefaultCacheKeepHeights)
	s.PrepareNextHeight(nil)
	return
}

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

func (s *Store) TrackOrphan(linkStoreKey []byte) error {
	// only track orphan if is cached store because we only prune on cache store
	if s.isCache {
		oldDataKey, _ := s.ParentDB.Get(linkStoreKey)
		if oldDataKey != nil && !bytes.Equal(oldDataKey, HashKey(linkStoreKey)) {
			orphanHeightKey := OrphanKey(HeightKey(s.Height+1, s.Prefix, oldDataKey))
			if err := s.ParentDB.Set(orphanHeightKey, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Store) DeleteCacheHeight(height int64) {
	if height < 0 {
		height = 0
	}
	keysToDelete := make([][]byte, 0)
	it := NewDedupIteratorForHeight(s.ParentDB, height, s.Prefix)
	for ; it.Valid(); it.Next() {
		linkKey := it.it.Key()
		keysToDelete = append(keysToDelete, linkKey)
	}
	it.Close()
	oIt := NewOrphanIteratorForHeight(s.ParentDB, height, s.Prefix)
	for ; oIt.Valid(); oIt.Next() {
		orphanKey := oIt.it.Key()
		linkKey := oIt.Key()
		keysToDelete = append(keysToDelete, orphanKey)
		keysToDelete = append(keysToDelete, linkKey)
	}
	oIt.Close()
	for _, k := range keysToDelete {
		if ok, _ := s.ParentDB.Has(k); !ok {
			fmt.Println("deleting a key that the cache store doesn't have")
		}
		if err := s.ParentDB.Delete(k); err != nil {
			panic("an error occurred deleting in cache")
		}
	}
	return
}

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

func (s *Store) CopyHeight(height int64, dedupStore *Store, b db.Batch) {
	it := NewDedupIteratorForHeight(s.ParentDB, height, s.Prefix)
	for ; it.Valid(); it.Next() {
		// linkstore key
		linkKey := it.it.Key()
		// datastore key
		dataKey := it.it.Value()
		b.Set(linkKey, dataKey)
		// data value
		if found, err := dedupStore.ParentDB.Has(dataKey); !found {
			if err != nil {
				panic("an error occurred in CopyHeight() .Has() func: " + err.Error())
			}
			dataValue, err := s.ParentDB.Get(dataKey)
			if err != nil {
				panic("an error occurred in CopyHeight() commit getting datavalue: " + err.Error())
			}
			b.Set(dataKey, dataValue)
		}
	}
	it.Close()
}

// key ops

func HeightKey(height int64, prefix string, k []byte) []byte {
	return []byte(fmt.Sprintf("%d/%s/%s", height, prefix, string(k)))
}

func FromHeightKey(heightKey string) (height int64, prefix string, k []byte) {
	var delim = "/"
	arr := strings.Split(heightKey, delim)
	// get height
	height, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		panic("unable to parse height from height key: " + heightKey)
	}
	prefix = arr[1]
	k = []byte(strings.Join(arr[2:], delim))
	return
}

func KeyFromHeightKey(heightKey []byte) (k []byte) {
	_, _, k = FromHeightKey(string(heightKey))
	return
}

func HashKey(key []byte) []byte {
	return sdk.Hash(key)
}

const orphanPrefix = "orphan/"

func OrphanKey(key []byte) []byte {
	return append([]byte(orphanPrefix), key...)
}

func KeyFromOrphanKey(orphanKey []byte) []byte {
	return orphanKey[len(orphanPrefix):]
}

// util

const (
	DefaultCacheKeepHeights = 25
)

func getPreloadStartHeight(latestHeight int64) int64 {
	startHeight := latestHeight - DefaultCacheKeepHeights
	if startHeight < 0 {
		startHeight = 0
	}
	return startHeight
}

// unused below

func (s *Store) Commit() types.CommitID {
	panic("Commit() called in de-dup store, when commitBatch should be used")
}
func (s *Store) CacheWrap() types.CacheWrap   { panic("cachewrap not implemented for de-dup store") }
func (s *Store) LastCommitID() types.CommitID { panic("lastCommitID not implemented for de-dup store") }
