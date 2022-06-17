package slim

import (
	"github.com/pokt-network/pocket-core/store/cachemulti"
	"github.com/pokt-network/pocket-core/store/types"
	db "github.com/tendermint/tm-db"
)

var _ types.CommitMultiStore = &MultiStore{}

type MultiStore struct {
	DB         *db.GoLevelDB
	Stores     map[types.StoreKey]types.CommitStore
	LastCommit types.CommitID
}

func (m *MultiStore) LoadLatestVersion() (err error) {
	latestHeight := getLatestVersion(m.DB)
	commitID := types.CommitID{}
	if latestHeight != 0 {
		commitID, err = getCommitID(m.DB, latestHeight)
	}
	for key := range m.Stores {
		m.Stores[key] = NewStoreWithIAVL(m.DB, latestHeight, key.Name(), commitID)
	}
	// reset next height upon launch incase there was 'dirty' data on the next working height
	return m.ResetNextHeight()
}

func (m *MultiStore) LoadVersion(ver int64) (store *types.Store, err error) {
	newStores := make(map[types.StoreKey]types.CommitStore)
	for key := range m.Stores {
		newStores[key] = NewStoreWithoutIAVL(m.DB, ver, key.Name())
	}
	return multiStoreToStore(m.DB, newStores), nil
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
		commitID, batch = s.(*Store).CommitBatch(batch)
		commitInfo.StoreInfos = append(commitInfo.StoreInfos, StoreInfo{
			Name: key.Name(),
			Core: StoreCore{commitID},
		})
	}
	setCommitInfo(batch, nextVersion, commitInfo)
	setLatestVersion(batch, nextVersion)
	_ = batch.Write()
	m.LastCommit = types.CommitID{
		Version: nextVersion,
		Hash:    commitInfo.Hash(),
	}
	return m.LastCommit
}

func (m *MultiStore) CopyStore() *types.Store {
	newStores := make(map[types.StoreKey]types.CommitStore)
	for key, store := range m.Stores {
		newStores[key] = store
	}
	return multiStoreToStore(m.DB, newStores)
}

func (m *MultiStore) ResetNextHeight() (err error) {
	batch := m.DB.NewBatch()
	defer batch.Close()
	for _, store := range m.Stores {
		batch, err = store.(*Store).Dedup.ResetNextHeight(batch)
		if err != nil {
			return err
		}
	}
	return batch.Write()
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
