package fastarchival

import (
	"fmt"
	"github.com/pokt-network/pocket-core/store/cachekv"
	"github.com/pokt-network/pocket-core/store/iavl"
	"github.com/pokt-network/pocket-core/store/rootmulti/heightcache"
	"github.com/pokt-network/pocket-core/store/tracekv"
	"github.com/pokt-network/pocket-core/types"
	dbm "github.com/tendermint/tm-db"
	"io"
)

var _ types.KVStore = (*Store)(nil)
var _ types.CommitStore = (*Store)(nil)
var _ types.CommitKVStore = (*Store)(nil)

type Store struct {
	merkle                *iavl.Store
	liveState             dbm.DB
	archival              dbm.DB
	storeKey              string
	height                int64
	isMutable             bool
	latestCommittedHeight int64
}

// Complying with types.CommitStore

func (s *Store) Commit() types.CommitID {
	commitID := s.merkle.Commit()

	it, err := s.liveState.Iterator(nil, nil)
	if err != nil {
		panic(fmt.Sprintf("unable to create an iterator for height %d storeKey %s in Commit()", s.height, s.storeKey))
	}
	defer it.Close()
	for ; it.Valid(); it.Next() {
		err := s.archival.Set(s.storePrefixForHeight(s.height+1, string(it.Key())), it.Value())
		if err != nil {
			panic("unable to commit to archival, data possibly corrupted.")
		}
	}
	return commitID
}

func (s *Store) LastCommitID() types.CommitID {
	if s.isMutable {
		return s.merkle.LastCommitID()
	}
	panic("'LastCommitID()' for called on an immutable store")
}

// SetPruning exists solely to comply with the CommitStore interface.
// it is a no-op and should be removed.
func (s *Store) SetPruning(options types.PruningOptions) {
}

// - End of Complying with types.CommitStore

// Complying with types.KVStore

// GetStoreType is here to comply with the KVStore interface
// returning the sanest answer in context. Should be adapted to fit reality.
func (s *Store) GetStoreType() types.StoreType {
	return types.StoreTypeIAVL
}

func (s *Store) CacheWrap() types.CacheWrap {
	return cachekv.NewStore(s)
}

func (s *Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return cachekv.NewStore(tracekv.NewStore(s, w, tc))
}

func (s *Store) Get(key []byte) ([]byte, error) {
	if s.isMutable {
		return s.liveState.Get(key)
	}
	return s.archival.Get(s.storePrefix(string(key)))
}

func (s *Store) Has(key []byte) (bool, error) {
	if s.isMutable {
		return s.liveState.Has(key)
	}
	return s.archival.Has(s.storePrefix(string(key)))
}

func (s *Store) Set(key, value []byte) error {
	if !s.isMutable {
		panic("'Set()' called on immutable store")
	}

	if err := s.merkle.Set(key, value); err != nil {
		panic("unable to set to iavl")
	}

	if err := s.liveState.Set(key, value); err != nil {
		panic("unable to set to live state")
	}
	return nil
}

func (s *Store) Delete(key []byte) error {
	if !s.isMutable {
		panic("'Delete()' called on immutable store")
	}
	if err := s.merkle.Delete(key); err != nil {
		panic("unable to delete from merkle tree")
	}
	if err := s.liveState.Delete(key); err != nil {
		panic("unable to delete from live state")
	}
	return nil
}

func (s *Store) Iterator(start, end []byte) (types.Iterator, error) {
	if s.isMutable {
		return s.liveState.Iterator(start, end)
	}
	archivalIterator, err := s.archival.Iterator(s.storePrefix(string(start)), s.storePrefix(string(end)))
	return ArchivalIterator{archivalIterator}, err
}

func (s *Store) ReverseIterator(start, end []byte) (types.Iterator, error) {
	if s.isMutable {
		return s.liveState.ReverseIterator(start, end)
	}
	archivalIterator, err := s.archival.ReverseIterator(s.storePrefix(string(start)), s.storePrefix(string(end)))
	return ArchivalIterator{archivalIterator}, err
}

// - End of Complying with types.KVStore

func NewStore(archival dbm.DB, storeKey string, commitID types.CommitID, height int64, isMutable bool, latestCommittedHeight int64) *Store {
	store := &Store{
		merkle:                nil,
		liveState:             nil,
		archival:              archival,
		storeKey:              storeKey,
		height:                height,
		isMutable:             isMutable,
		latestCommittedHeight: latestCommittedHeight,
	}
	if isMutable {
		prefix := store.storePrefix("")
		it, err := archival.Iterator(prefix, types.PrefixEndBytes(prefix))
		if err != nil {
			panic(fmt.Sprintf("unable to create an iterator for height %d with storeKey %s", height, storeKey))
		}
		defer it.Close()
		for ; it.Valid(); it.Next() {
			err = store.liveState.Set(StoreKeySuffix(it.Key()), it.Value())
			if err != nil {
				panic("unable to set k/v in state: " + err.Error())
			}
		}
		store.merkle, err = iavl.LoadStore(dbm.NewPrefixDB(archival, []byte("s/k:"+storeKey+"/")), commitID, types.PruningOptions{}, false, heightcache.InvalidCache{}, 0)
		if err != nil {
			panic("unable to load iavlStore in rootmultistore: " + err.Error())
		}
	}
	return store
}

func StoreKeySuffix(storeKey []byte) []byte {
	delim := 0
	for i, b := range storeKey {
		if b == byte('/') {
			delim++
		}
		if delim == 2 {
			return storeKey[i+1:]
		}
	}
	panic("attempted to get suffix from store key that doesn't have exactly 2 delims")
}

func (s *Store) storePrefix(datumKey string) []byte {
	return s.storePrefixForHeight(s.height, datumKey)
}

func (s *Store) storePrefixForHeight(height int64, datumKey string) []byte {
	if s.storeKey == "" {
		return []byte(fmt.Sprintf("%d/", height))
	}
	if datumKey == "" {
		return []byte(fmt.Sprintf("%d/%s/", height, s.storeKey))
	}
	return []byte(fmt.Sprintf("%d/%s/%s", height, s.storeKey, datumKey))
}
