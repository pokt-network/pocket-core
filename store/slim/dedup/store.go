package dedup

import (
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
	ParentDB db.GoLevelDB
}

func NewStore(height int64, prefix string, parent db.GoLevelDB) Store {
	return Store{
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

func (s *Store) CommitBatch(b db.Batch) (types.CommitID, db.Batch) {
	s.Height = s.Height + 1
	return types.CommitID{}, s.PrepareNextHeight(b)
}

func (s *Store) PrepareNextHeight(b db.Batch) db.Batch {
	var err error
	start := HeightKey(s.Height-1, s.Prefix, nil)
	end := types.PrefixEndBytes(start)
	it := &DedupIterator{parent: s.ParentDB}
	it.it, err = s.ParentDB.Iterator(start, end)
	if err != nil {
		panic(fmt.Sprintf("unable to create an iterator for height/prefix %d/%s in Commit()", s.Height, s.Prefix))
	}
	nextHeight := s.Height
	defer it.Close()
	for ; it.Valid(); it.Next() {
		k := it.Key()
		nextHeightKey := HeightKey(nextHeight, s.Prefix, k)
		linkValue := it.it.Value()
		b.Set(nextHeightKey, linkValue)
	}
	return b
}

// The two functions below are needed because of the 'write immediately' db design where the
// data is written each op to the db instead of saving it for commit() phase at the end.
// While this is conceptually simpler, it comes with the tradeoff that the next working height
// must be reset upon restart.
// NOTE: The assumption is that the entire block is rolled back upon replay. If this isn't the case
// this actually will be both unnecessary (see cache-wrapping) and not work.

func (s *Store) ResetNextHeight(b db.Batch) (batch db.Batch, err error) {
	b, err = s.ClearNextHeight(b)
	if err != nil {
		return b, err
	}
	return s.PrepareNextHeight(b), nil
}

func (s *Store) ClearNextHeight(b db.Batch) (db.Batch, error) {
	// iterate through the LINKSTORE to clear the 'next height'
	// we need to do this in case the db was shut down at an
	// unsafe point
	nextHeight := s.Height
	startKey := HeightKey(nextHeight, s.Prefix, nil)
	endKey := types.PrefixEndBytes(startKey)
	it, err := s.Iterator(startKey, endKey)
	if err != nil {
		panic(fmt.Sprintf("unable to create an iterator for height/prefix %d/%s in ClearNextHeight()", s.Height, s.Prefix))
	}
	keysToDelete := make([][]byte, 0)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		linkKey := it.Key()
		keysToDelete = append(keysToDelete, linkKey)
		// delete any data that was created for next height as well
		// or the result will be orphaned data that has no link
		dataKey := HashKey(linkKey)
		keysToDelete = append(keysToDelete, dataKey)
	}
	for _, k := range keysToDelete {
		b.Delete(k)
	}
	return b, nil
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

// unused below

func (s *Store) Commit() types.CommitID {
	panic("Commit() called in de-dup store, when commitBatch should be used")
}
func (s *Store) CacheWrap() types.CacheWrap   { panic("cachewrap not implemented for de-dup store") }
func (s *Store) LastCommitID() types.CommitID { panic("lastCommitID not implemented for de-dup store") }
