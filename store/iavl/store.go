package iavl

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/pokt-network/pocket-core/store/cachekv"
	serrors "github.com/pokt-network/pocket-core/store/errors"
	"github.com/pokt-network/pocket-core/store/tracekv"
	"github.com/pokt-network/pocket-core/store/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	dbm "github.com/tendermint/tm-db"
)

const (
	defaultIAVLCacheSize = 5000000
)

// LoadStore loads the iavl store
func LoadStore(db dbm.DB, id types.CommitID, pruning types.PruningOptions, lazyLoading bool, cache types.SingleStoreCache, iavlCacheSize int64) (types.CommitStore, error) {
	var err error

	if iavlCacheSize <= 0 {
		iavlCacheSize = defaultIAVLCacheSize
	}

	tree, err := NewMutableTree(db, int(iavlCacheSize))
	if err != nil {
		return nil, err
	}

	if lazyLoading {
		_, err = tree.LazyLoadVersion(id.Version)
	} else {
		_, err = tree.LoadVersion(id.Version)
	}

	if err != nil {
		return nil, err
	}

	iavl := UnsafeNewStore(tree, int64(0), int64(0), cache)
	iavl.SetPruning(pruning)

	if iavl.cache.IsValid() {
		cacheDatasetIterator, _ := iavl.Iterator(nil, nil)
		cacheDataset := make(map[string]string, 10)
		for cacheDatasetIterator.Valid() {
			cacheDataset[string(cacheDatasetIterator.Key())] = string(cacheDatasetIterator.Value())
			cacheDatasetIterator.Next()
		}

		iavl.cache.Initialize(cacheDataset, iavl.tree.Version())
		log.Println("Cache warmed.")
	}
	return iavl, nil
}

//----------------------------------------

var _ types.KVStore = (*Store)(nil)
var _ types.CommitStore = (*Store)(nil)
var _ types.Queryable = (*Store)(nil)

// Store Implements types.KVStore and CommitStore.
type Store struct {
	tree Tree

	// How many old versions we hold onto.
	// A value of 0 means keep no recent states.
	numRecent int64

	// This is the distance between state-sync waypoint states to be stored.
	// See https://github.com/tendermint/tendermint/issues/828
	// A value of 1 means store every state.
	// A value of 0 means store no waypoints. (node cannot assist in state-sync)
	// By default this value should be set the same across all nodes,
	// so that nodes can know the waypoints their peers store.
	storeEvery int64

	cache types.SingleStoreCache
}

// CONTRACT: tree should be fully loaded.
// nolint: unparam
func UnsafeNewStore(tree *MutableTree, numRecent int64, storeEvery int64, cache types.SingleStoreCache) *Store {
	st := &Store{
		tree:       tree,
		numRecent:  numRecent,
		storeEvery: storeEvery,
		cache:      cache,
	}
	return st
}

// LoadLazyVersion returns a reference to a new store backed by an immutable IAVL
// tree at a specific version (height) without any pruning options. This should
// be used for querying and iteration only. If the version does not exist or has
// been pruned, an error will be returned. Any mutable operations executed will
// result in a panic.
func (st *Store) LazyLoadStore(version int64, cache types.SingleStoreCache) (*Store, error) {
	a, ok := st.tree.(*MutableTree)
	if !ok {
		return nil, fmt.Errorf("not immutable tree in LazyLoadStore")
	}

	tree, err := a.LazyLoadVersion(version)
	if err != nil {
		return nil, err
	}
	iavl := UnsafeNewStore(tree, int64(0), int64(0), cache)
	return iavl, nil
}

func (st *Store) Rollback(version int64) error {
	r, ok := st.tree.(*MutableTree)
	if !ok {
		return fmt.Errorf("cant turn st.Tree into mutable tree for rollback")
	}
	_, err := r.LoadVersionForOverwriting(version)
	if err != nil {
		return err
	}
	return nil
}

// Implements Committer.
func (st *Store) Commit() types.CommitID {
	// Save a new version.
	hash, version, err := st.tree.SaveVersion()
	st.cache.Commit(version)
	if err != nil {
		// TODO: Do we want to extend Commit to allow returning errors?
		panic(err)
	}

	// Release an old version of history, if not a sync waypoint.
	//previous := version - 1 TODO removed for testing
	//if st.numRecent < previous {
	//	toRelease := previous - st.numRecent
	//	if st.storeEvery == 0 || toRelease%st.storeEvery != 0 {
	//		err := st.tree.DeleteVersion(toRelease)
	//		if errCause := errors.Cause(err); errCause != nil && errCause != iavl.ErrVersionDoesNotExist {
	//			panic(err)
	//		}
	//	}
	//}

	return types.CommitID{
		Version: version,
		Hash:    hash,
	}
}

// Implements Committer.
func (st *Store) LastCommitID() types.CommitID {
	return types.CommitID{
		Version: st.tree.Version(),
		Hash:    st.tree.Hash(),
	}
}

// Implements Committer.
func (st *Store) SetPruning(opt types.PruningOptions) {
	st.numRecent = opt.KeepRecent()
	st.storeEvery = opt.KeepEvery()
}

// VersionExists returns whether or not a given version is stored.
func (st *Store) VersionExists(version int64) bool {
	return st.tree.VersionExists(version)
}

// Implements Store.
func (st *Store) GetStoreType() types.StoreType {
	return types.StoreTypeIAVL
}

// Implements Store.
func (st *Store) CacheWrap() types.CacheWrap {
	return cachekv.NewStore(st)
}

// CacheWrapWithTrace implements the Store interface.
func (st *Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return cachekv.NewStore(tracekv.NewStore(st, w, tc))
}

// Implements types.KVStore.
func (st *Store) Set(key, value []byte) error {
	types.AssertValidValue(value)
	st.tree.Set(key, value)
	st.cache.Set(key, value)
	return nil
}

// Implements types.KVStore.
func (st *Store) Get(key []byte) (value []byte, err error) {
	if val, err := st.cache.Get(st.tree.Version(), key); err == nil {
		return val, nil
	}
	_, v := st.tree.Get(key)
	return v, nil
}

// Implements types.KVStore.
func (st *Store) Has(key []byte) (exists bool, err error) {
	if val, err := st.cache.Has(st.tree.Version(), key); err == nil {
		return val, nil
	}
	return st.tree.Has(key), nil
}

// Implements types.KVStore.
func (st *Store) Delete(key []byte) error {
	st.tree.Remove(key)
	st.cache.Remove(key)
	return nil
}

// Implements types.KVStore.
func (st *Store) Iterator(start, end []byte) (types.Iterator, error) {
	var iTree *ImmutableTree

	switch tree := st.tree.(type) {
	case *immutableTree:
		iTree = tree.ImmutableTree
	case *MutableTree:
		iTree = tree.ImmutableTree
	}
	if val, err := st.cache.Iterator(iTree.version, start, end); err == nil {
		return val, nil
	}
	return newIAVLIterator(iTree, start, end, true), nil
}

// Implements types.KVStore.
func (st *Store) ReverseIterator(start, end []byte) (types.Iterator, error) {
	var iTree *ImmutableTree

	switch tree := st.tree.(type) {
	case *immutableTree:
		iTree = tree.ImmutableTree
	case *MutableTree:
		iTree = tree.ImmutableTree
	}
	if val, err := st.cache.ReverseIterator(iTree.version, start, end); err == nil {
		return val, nil
	}

	return newIAVLIterator(iTree, start, end, false), nil
}

// Handle gatest the latest height, if height is 0
func getHeight(tree Tree, req abci.RequestQuery) int64 {
	height := req.Height
	if height == 0 {
		latest := tree.Version()
		if tree.VersionExists(latest - 1) {
			height = latest - 1
		} else {
			height = latest
		}
	}
	return height
}

// Query implements ABCI interface, allows queries
//
// by default we will return from (latest height -1),
// as we will have merkle proofs immediately (header height = data height + 1)
// If latest-1 is not present, use latest (which must be present)
// if you care to have the latest data to see a tx results, you must
// explicitly set the height you want to see
func (st *Store) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if len(req.Data) == 0 {
		msg := "Query cannot be zero length"
		return serrors.ErrTxDecode(msg).QueryResult()
	}

	tree := st.tree

	// store the height we chose in the response, with 0 being changed to the
	// latest height
	res.Height = getHeight(tree, req)

	switch req.Path {
	case "/key": // get by key
		key := req.Data // data holds the key bytes

		res.Key = key
		if !st.VersionExists(res.Height) {
			res.Log = ErrVersionDoesNotExist.Error()
			break
		}

		if req.Prove {
			value, proof, err := tree.GetVersionedWithProof(key, res.Height)
			if err != nil {
				res.Log = err.Error()
				break
			}
			if proof == nil {
				// Proof == nil implies that the store is empty.
				if value != nil {
					panic("unexpected value for an empty proof")
				}
			}
			if value != nil {
				// value was found
				res.Value = value
				res.Proof = &merkle.Proof{Ops: []merkle.ProofOp{NewValueOp(key, proof).ProofOp()}}
			} else {
				// value wasn't found
				res.Value = nil
				res.Proof = &merkle.Proof{Ops: []merkle.ProofOp{NewAbsenceOp(key, proof).ProofOp()}}
			}
		} else {
			_, res.Value = tree.GetVersioned(key, res.Height)
		}

	case "/subspace":
		var KVs []types.KVPair

		subspace := req.Data
		res.Key = subspace

		iterator, _ := types.KVStorePrefixIterator(st, subspace)
		for ; iterator.Valid(); iterator.Next() {
			KVs = append(KVs, types.KVPair{Key: iterator.Key(), Value: iterator.Value()})
		}

		iterator.Close()
		res.Value, _ = cdc.LegacyMarshalBinaryLengthPrefixed(KVs)

	default:
		msg := fmt.Sprintf("Unexpected Query path: %v", req.Path)
		return serrors.ErrUnknownRequest(msg).QueryResult()
	}

	return
}

//----------------------------------------

// Implements types.Iterator.
type iavlIterator struct {
	// Underlying store
	tree *ImmutableTree

	// Domain
	start, end []byte

	// Iteration order
	ascending bool

	// Channel to push iteration values.
	iterCh chan kv.Pair

	// Close this to release goroutine.
	quitCh chan struct{}

	// Close this to signal that state is initialized.
	initCh chan struct{}

	//----------------------------------------
	// What follows are mutable state.
	mtx sync.Mutex

	invalid bool   // True once, true forever
	key     []byte // The current key
	value   []byte // The current value
}

func (iter *iavlIterator) Error() error {
	panic("implement me")
}

var _ types.Iterator = (*iavlIterator)(nil)

// newIAVLIterator will create a new iavlIterator.
// CONTRACT: Caller must release the iavlIterator, as each one creates a new
// goroutine.
func newIAVLIterator(tree *ImmutableTree, start, end []byte, ascending bool) *iavlIterator {
	iter := &iavlIterator{
		tree:      tree,
		start:     types.Cp(start),
		end:       types.Cp(end),
		ascending: ascending,
		iterCh:    make(chan kv.Pair), // Set capacity > 0?
		quitCh:    make(chan struct{}),
		initCh:    make(chan struct{}),
	}
	go iter.iterateRoutine()
	go iter.initRoutine()
	return iter
}

// Run this to funnel items from the tree to iterCh.
func (iter *iavlIterator) iterateRoutine() {
	iter.tree.IterateRange(
		iter.start, iter.end, iter.ascending,
		func(key, value []byte) bool {
			select {
			case <-iter.quitCh:
				return true // done with iteration.
			case iter.iterCh <- kv.Pair{Key: key, Value: value}:
				return false // yay.
			}
		},
	)
	close(iter.iterCh) // done.
}

// Run this to fetch the first item.
func (iter *iavlIterator) initRoutine() {
	iter.receiveNext()
	close(iter.initCh)
}

// Implements types.Iterator.
func (iter *iavlIterator) Domain() (start, end []byte) {
	return iter.start, iter.end
}

// Implements types.Iterator.
func (iter *iavlIterator) Valid() bool {
	iter.waitInit()
	iter.mtx.Lock()

	validity := !iter.invalid
	iter.mtx.Unlock()
	return validity
}

// Implements types.Iterator.
func (iter *iavlIterator) Next() {
	iter.waitInit()
	iter.mtx.Lock()
	iter.assertIsValid(true)

	iter.receiveNext()
	iter.mtx.Unlock()
}

// Implements types.Iterator.
func (iter *iavlIterator) Key() []byte {
	iter.waitInit()
	iter.mtx.Lock()
	iter.assertIsValid(true)

	key := iter.key
	iter.mtx.Unlock()
	return key
}

// Implements types.Iterator.
func (iter *iavlIterator) Value() []byte {
	iter.waitInit()
	iter.mtx.Lock()
	iter.assertIsValid(true)

	val := iter.value
	iter.mtx.Unlock()
	return val
}

// Implements types.Iterator.
func (iter *iavlIterator) Close() {
	close(iter.quitCh)
}

//----------------------------------------

func (iter *iavlIterator) setNext(key, value []byte) {
	iter.assertIsValid(false)

	iter.key = key
	iter.value = value
}

func (iter *iavlIterator) setInvalid() {
	iter.assertIsValid(false)

	iter.invalid = true
}

func (iter *iavlIterator) waitInit() {
	<-iter.initCh
}

func (iter *iavlIterator) receiveNext() {
	kvPair, ok := <-iter.iterCh
	if ok {
		iter.setNext(kvPair.Key, kvPair.Value)
	} else {
		iter.setInvalid()
	}
}

// assertIsValid panics if the iterator is invalid. If unlockMutex is true,
// it also unlocks the mutex before panicing, to prevent deadlocks in code that
// recovers from panics
func (iter *iavlIterator) assertIsValid(unlockMutex bool) {
	if iter.invalid {
		if unlockMutex {
			iter.mtx.Unlock()
		}
		panic("invalid iterator")
	}
}
