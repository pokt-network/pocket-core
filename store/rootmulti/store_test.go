package rootmulti

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/merkle"
	dbm "github.com/tendermint/tm-db"

	"github.com/pokt-network/pocket-core/store/types"
)

func TestStoreType(t *testing.T) {
	db := dbm.NewMemDB()
	store := NewStore(db)
	store.MountStoreWithDB(
		types.NewKVStoreKey("store1"), types.StoreTypeIAVL, db)

}

func TestStoreMount(t *testing.T) {
	db := dbm.NewMemDB()
	store := NewStore(db)

	key1 := types.NewKVStoreKey("store1")
	key2 := types.NewKVStoreKey("store2")
	dup1 := types.NewKVStoreKey("store1")

	require.NotPanics(t, func() { store.MountStoreWithDB(key1, types.StoreTypeIAVL, db) })
	require.NotPanics(t, func() { store.MountStoreWithDB(key2, types.StoreTypeIAVL, db) })

	require.Panics(t, func() { store.MountStoreWithDB(key1, types.StoreTypeIAVL, db) })
	require.Panics(t, func() { store.MountStoreWithDB(dup1, types.StoreTypeIAVL, db) })
}

func TestHashStableWithEmptyCommit(t *testing.T) {
	var db dbm.DB = dbm.NewMemDB()
	ms := newMultiStoreWithMounts(db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	commitID := types.CommitID{}
	checkStore(t, ms, commitID, commitID)

	k, v := []byte("wind"), []byte("blows")

	store1 := ms.getStoreByName("store1").(types.KVStore)
	_ = store1.Set(k, v)

	cID := ms.Commit()
	require.Equal(t, int64(1), cID.Version)
	hash := cID.Hash

	// make an empty commit, it should update version, but not affect hash
	cID = ms.Commit()
	require.Equal(t, int64(2), cID.Version)
	require.Equal(t, hash, cID.Hash)
}

func TestMultistoreCommitLoad(t *testing.T) {
	var db dbm.DB = dbm.NewMemDB()
	store := newMultiStoreWithMounts(db)
	err := store.LoadLatestVersion()
	require.Nil(t, err)

	// New store has empty last commit.
	commitID := types.CommitID{}
	checkStore(t, store, commitID, commitID)

	// Make sure we can get stores by name.
	s1 := store.getStoreByName("store1")
	require.NotNil(t, s1)
	s3 := store.getStoreByName("store3")
	require.NotNil(t, s3)
	s77 := store.getStoreByName("store77")
	require.Nil(t, s77)

	// Make a few commits and check them.
	nCommits := int64(3)
	for i := int64(0); i < nCommits; i++ {
		commitID = store.Commit()
		expectedCommitID := getExpectedCommitID(store, i+1)
		checkStore(t, store, expectedCommitID, commitID)
	}

	// Load the latest multistore again and check version.
	store = newMultiStoreWithMounts(db)
	err = store.LoadLatestVersion()
	require.Nil(t, err)
	commitID = getExpectedCommitID(store, nCommits)
	checkStore(t, store, commitID, commitID)

	// Commit and check version.
	commitID = store.Commit()
	expectedCommitID := getExpectedCommitID(store, nCommits+1)
	checkStore(t, store, expectedCommitID, commitID)

	// Load an older multistore and check version.
	ver := nCommits - 1
	store = newMultiStoreWithMounts(db)
	err = store.LoadVersion(ver)
	require.Nil(t, err)
	commitID = getExpectedCommitID(store, ver)
	checkStore(t, store, commitID, commitID)

	// XXX: commit this older version
	commitID = store.Commit()
	expectedCommitID = getExpectedCommitID(store, ver+1)
	checkStore(t, store, expectedCommitID, commitID)

	// XXX: confirm old commit is overwritten and we have rolled back
	// LatestVersion
	store = newMultiStoreWithMounts(db)
	err = store.LoadLatestVersion()
	require.Nil(t, err)
	commitID = getExpectedCommitID(store, ver+1)
	checkStore(t, store, commitID, commitID)
}

func TestParsePath(t *testing.T) {
	_, _, err := parsePath("foo")
	require.Error(t, err)

	store, subpath, err := parsePath("/foo")
	require.NoError(t, err)
	require.Equal(t, store, "foo")
	require.Equal(t, subpath, "")

	store, subpath, err = parsePath("/fizz/bang/baz")
	require.NoError(t, err)
	require.Equal(t, store, "fizz")
	require.Equal(t, subpath, "/bang/baz")

	substore, subsubpath, err := parsePath(subpath)
	require.NoError(t, err)
	require.Equal(t, substore, "bang")
	require.Equal(t, subsubpath, "/baz")

}

//-----------------------------------------------------------------------
// utils

func newMultiStoreWithMounts(db dbm.DB) *Store {
	store := NewStore(db)
	store.MountStoreWithDB(
		types.NewKVStoreKey("store1"), types.StoreTypeIAVL, nil)
	store.MountStoreWithDB(
		types.NewKVStoreKey("store2"), types.StoreTypeIAVL, nil)
	store.MountStoreWithDB(
		types.NewKVStoreKey("store3"), types.StoreTypeIAVL, nil)
	return store
}

func checkStore(t *testing.T, store *Store, expect, got types.CommitID) {
	require.Equal(t, expect, got)
	require.Equal(t, expect, store.LastCommitID())

}

func getExpectedCommitID(store *Store, ver int64) types.CommitID {
	return types.CommitID{
		Version: ver,
		Hash:    hashStores(store.stores),
	}
}

func hashStores(stores map[types.StoreKey]types.CommitStore) []byte {
	m := make(map[string][]byte, len(stores))
	for key, store := range stores {
		name := key.Name()
		m[name] = StoreInfo{
			Name: name,
			Core: StoreCore{
				CommitID: store.LastCommitID(),
				// StoreType: store.GetStoreType(),
			},
		}.Hash()
	}
	return merkle.SimpleHashFromMap(m)
}
