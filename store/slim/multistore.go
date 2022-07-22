package slim

import (
	"fmt"
	"github.com/pokt-network/pocket-core/store/cachemulti"
	"github.com/pokt-network/pocket-core/store/slim/cache"
	"github.com/pokt-network/pocket-core/store/types"
	db "github.com/tendermint/tm-db"
	"time"
)

var _ types.CommitMultiStore = &MultiStore{}

// Multi stores are abstractions over the db; breaking it up into a few 'prefix' stores
// this is built in Cosmos SDK architecture that can't be removed without modifying the
// entire app structure.

type MultiStore struct {
	DB         *db.GoLevelDB
	CacheDB    *cache.CacheDB
	Stores     map[types.StoreKey]types.CommitStore
	LastCommit types.CommitID
}

func NewStore(d db.DB) *MultiStore {
	return &MultiStore{
		DB:         d.(*db.GoLevelDB),
		CacheDB:    cache.NewCacheDB(d, 0, maxCacheKeepHeights),
		Stores:     make(map[types.StoreKey]types.CommitStore),
		LastCommit: types.CommitID{},
	}
}

func (m *MultiStore) LoadLatestVersion() (err error) {
	latestHeight := getLatestVersion(m.DB)
	commitID := types.CommitID{}
	if latestHeight != 0 {
		commitID, err = getCommitID(m.DB, latestHeight)
		if err != nil {
			return err
		}
		m.LastCommit = commitID
	}
	for key := range m.Stores {
		m.Stores[key] = NewStoreWithIAVL(m.DB, m.CacheDB, latestHeight, key.Name(), commitID)
	}
	m.Preload(latestHeight)
	return nil
}

func (m *MultiStore) LoadVersion(ver int64) (store *types.Store, err error) {
	newStores := make(map[types.StoreKey]types.CommitStore)
	for key := range m.Stores {
		newStores[key] = NewStoreWithoutIAVL(m.DB, m.CacheDB, m.LastCommit.Version, ver-1, key.Name())
	}
	return multiStoreToStore(m.DB, m.CacheDB, m.LastCommit, newStores), nil
}

func (m *MultiStore) Commit() (commitID types.CommitID) {
	batch := m.DB.NewBatch()
	defer batch.Close()
	nextVersion := m.LastCommit.Version + 1
	commitInfo := CommitInfo{
		Version:    nextVersion,
		StoreInfos: make([]StoreInfo, 0),
	}
	for key, s := range m.Stores {
		commitID = s.(*Store).Commit()
		commitInfo.StoreInfos = append(commitInfo.StoreInfos, StoreInfo{
			Name: key.Name(),
			Core: StoreCore{commitID},
		})
	}
	m.CommitCache() // write preload to disk & prune cache
	setCommitInfo(batch, nextVersion, commitInfo)
	setLatestVersion(batch, nextVersion)
	_ = batch.Write()
	m.LastCommit = types.CommitID{
		Version: nextVersion,
		Hash:    commitInfo.Hash(),
	}
	return m.LastCommit
}

func (m *MultiStore) Preload(latestVersion int64) {
	m.CacheDB.LatestHeight = latestVersion
	m.CacheDB.Preload()
}

func (m *MultiStore) CommitCache() {
	// commit cache
	t := time.Now()
	m.CacheDB.Commit()
	fmt.Println("Cache commit took: " + time.Since(t).String())
	// prune cache
	oldestHeight := m.LastCommit.Version - maxCacheKeepHeights
	if oldestHeight <= 1 {
		return
	}
	m.CacheDB.DeleteHeight(oldestHeight)
}

func (m *MultiStore) CopyStore() *types.Store {
	newStores := make(map[types.StoreKey]types.CommitStore)
	for key, store := range m.Stores {
		newStores[key] = store
	}
	return multiStoreToStore(m.DB, m.CacheDB, m.LastCommit, newStores)
}

func (m *MultiStore) LastCommitID() types.CommitID { return m.LastCommit }
func (m *MultiStore) CacheWrap() types.CacheWrap   { return m.CacheMultiStore() }
func (m *MultiStore) CacheMultiStore() types.CacheMultiStore {
	return cachemulti.NewCacheMulti(m.Stores)
}
func (m *MultiStore) GetKVStore(key types.StoreKey) types.KVStore {
	return m.Stores[key].(types.KVStore)
}
func (m *MultiStore) MountStoreWithDB(key types.StoreKey, _ types.StoreType, _ db.DB) {
	m.Stores[key] = nil
}
