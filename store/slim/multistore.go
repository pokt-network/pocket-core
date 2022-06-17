package slim

import (
	"github.com/pokt-network/pocket-core/store/types"
	db "github.com/tendermint/tm-db"
)

var _ types.CommitMultiStore = &MultiStore{}

type MultiStore struct {
	DB                   db.GoLevelDB
	LatestHeight         int64
	LoadedHeight         int64
	OldestHeightInRecent int64
}

func (m *MultiStore) Commit() types.CommitID {
	//TODO implement me
	panic("implement me")
}

func (m *MultiStore) LastCommitID() types.CommitID {
	//TODO implement me
	panic("implement me")
}

func (m *MultiStore) CacheWrap() types.CacheWrap {
	//TODO implement me
	panic("implement me")
}

func (m *MultiStore) CacheMultiStore() types.CacheMultiStore {
	//TODO implement me
	panic("implement me")
}

func (m *MultiStore) GetKVStore(key types.StoreKey) types.KVStore {
	// TODO implement me
	panic("implement me")
}

func (m *MultiStore) MountStoreWithDB(key types.StoreKey, typ types.StoreType, db db.DB) {
	//TODO implement me
	panic("implement me")
}

func (m *MultiStore) LoadLatestVersion() error {
	//TODO implement me
	panic("implement me")
}

func (m *MultiStore) LoadVersion(ver int64) (*types.Store, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MultiStore) CopyStore() *types.Store {
	//TODO implement me
	panic("implement me")
}
