package types

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/kv"

	dbm "github.com/tendermint/tm-db"
)

type Store interface { //nolint
	GetStoreType() StoreType
	CacheWrapper
}

// something that can persist to disk
type Committer interface {
	Commit() CommitID
	LastCommitID() CommitID
}

// Stores of MultiStore must implement CommitStore.
type CommitStore interface {
	Committer
	Store
}

//----------------------------------------
// MultiStore

type MultiStore interface { //nolint
	Store

	// Cache wrap MultiStore.
	// NOTE: Caller should probably not call .Write() on each, but
	// call CacheMultiStore.Write().
	CacheMultiStore() CacheMultiStore

	// Convenience for fetching substores.
	// If the store does not exist, panics.
	GetStore(StoreKey) Store
	GetKVStore(StoreKey) KVStore
}

// From MultiStore.CacheMultiStore()....
type CacheMultiStore interface {
	MultiStore
	Write() // Writes operations to underlying KVStore
}

// A non-cache MultiStore.
type CommitMultiStore interface {
	Committer
	MultiStore

	// Mount a store of type using the given db.
	// If db == nil, the new store will use the CommitMultiStore db.
	MountStoreWithDB(key StoreKey, typ StoreType, db dbm.DB)

	// Panics on a nil key.
	GetCommitStore(key StoreKey) CommitStore

	// Panics on a nil key.
	GetCommitKVStore(key StoreKey) CommitKVStore

	// Load the latest persisted version. Called once after all calls to
	// Mount*Store() are complete.
	LoadLatestVersion() error
	// Load a specific persisted version in a memory saving fashion.
	// Don't iterate through and collect all the roots and versions
	LoadHistoricalVersion(ver int64) (*Store, error)
	CopyStore() *Store
}

//---------subsp-------------------------------
// KVStore

// KVStore is a simple interface to get/set data
type KVStore interface {
	Store

	// Get returns nil iff key doesn't exist. Panics on nil key.
	Get(key []byte) ([]byte, error)

	// Has checks if a key exists. Panics on nil key.
	Has(key []byte) (bool, error)

	// Set sets the key. Panics on nil key or value.
	Set(key, value []byte) error

	// Delete deletes the key. Panics on nil key.
	Delete(key []byte) error

	// Iterator over a domain of keys in ascending order. End is exclusive.
	// Start must be less than end, or the Iterator is invalid.
	// Iterator must be closed by caller.
	// To iterate over entire domain, use store.Iterator(nil, nil)
	// CONTRACT: No writes may happen within a domain while an iterator exists over it.
	// Exceptionally allowed for cachekv.Store, safe to write in the modules.
	Iterator(start, end []byte) (Iterator, error)

	// Iterator over a domain of keys in descending order. End is exclusive.
	// Start must be less than end, or the Iterator is invalid.
	// Iterator must be closed by caller.
	// CONTRACT: No writes may happen within a domain while an iterator exists over it.
	// Exceptionally allowed for cachekv.Store, safe to write in the modules.
	ReverseIterator(start, end []byte) (Iterator, error)
}

// Alias iterator to db's Iterator for convenience.
type Iterator = dbm.Iterator

// CacheKVStore cache-wraps a KVStore.  After calling .Write() on
// the CacheKVStore, all previously created CacheKVStores on the
// object expire.
type CacheKVStore interface {
	KVStore

	// Writes operations to underlying KVStore
	Write()
}

// Stores of MultiStore must implement CommitStore.
type CommitKVStore interface {
	Committer
	KVStore
}

//----------------------------------------
// CacheWrap

// CacheWrap makes the most appropriate cache-wrap. For example,
// IAVLStore.CacheWrap() returns a CacheKVStore. CacheWrap should not return
// a Committer, since Commit cache-wraps make no sense. It can return KVStore,
// HeapStore, SpaceStore, etc.
type CacheWrap interface {
	// Write syncs with the underlying store.
	Write()

	// CacheWrap recursively wraps again.
	CacheWrap() CacheWrap
}

type CacheWrapper interface { //nolint
	// CacheWrap cache wraps.
	CacheWrap() CacheWrap
}

//----------------------------------------
// CommitID

// CommitID contains the tree version number and its merkle root.
// type CommitID struct {
// 	Version int64
// 	Hash    []byte
// }

func (cid CommitID) IsZero() bool { //nolint
	return cid.Version == 0 && len(cid.Hash) == 0
}

func (cid CommitID) String() string {
	return fmt.Sprintf("CommitID{%v:%X}", cid.Hash, cid.Version)
}

//----------------------------------------
// Store types

// kind of store
type StoreType int

const (
	//nolint
	StoreTypeMulti StoreType = iota
	StoreTypeDB
	StoreTypeIAVL
)

//----------------------------------------
// PublicKeys for accessing substores

// storeKey is a key used to index stores in a MultiStore.
type StoreKey interface {
	Name() string
	String() string
}

// KVStoreKey is used for accessing substores.
// Only the pointer value should ever be used - it functions as a capabilities key.
type KVStoreKey struct {
	name string
}

// NewKVStoreKey returns a new pointer to a KVStoreKey.
// Use a pointer so keys don't collide.
func NewKVStoreKey(name string) *KVStoreKey {
	return &KVStoreKey{
		name: name,
	}
}

func (key *KVStoreKey) Name() string {
	return key.name
}

func (key *KVStoreKey) String() string {
	return fmt.Sprintf("KVStoreKey{%p, %s}", key, key.name)
}

// key-value result for iterator queries
type KVPair kv.Pair
