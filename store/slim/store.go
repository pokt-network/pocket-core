package slim

import (
	"github.com/jackc/pgx/v4"
	"github.com/pokt-network/pocket-core/store/iavl"
	"github.com/pokt-network/pocket-core/store/slim/dedup"
	"github.com/pokt-network/pocket-core/store/slim/historical"
	"github.com/pokt-network/pocket-core/store/types"
	db "github.com/tendermint/tm-db"
)

var _ types.KVStore = &Store{}
var _ types.CommitStore = &Store{}

type Store struct {
	RecentStore     dedup.Store
	HistoricalStore historical.Store
	IAVLStore       iavl.Store
	OnlyHistorical  bool
}

func NewStore(height int64, prefix string, onlyHistorical bool, commitID types.CommitID, parentForIAVL db.GoLevelDB, parentForRecentStore db.GoLevelDB, parentForHistoricalStore *pgx.Conn) *Store {
	if onlyHistorical {
		return &Store{
			HistoricalStore: historical.NewStore(height, prefix, parentForHistoricalStore),
			OnlyHistorical:  onlyHistorical,
		}
	}
	iavlStore, err := iavl.NewStore(&parentForIAVL, commitID)
	if err != nil {
		panic("iavl store failed to load for height: %s prefix: %s")
	}
	return &Store{
		RecentStore:     dedup.NewStore(height, prefix, parentForRecentStore),
		HistoricalStore: historical.NewStore(height, prefix, parentForHistoricalStore),
		IAVLStore:       *iavlStore,
		OnlyHistorical:  false,
	}
}

func (s *Store) Get(key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Has(key []byte) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Set(key, value []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Delete(key []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Iterator(start, end []byte) (types.Iterator, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) ReverseIterator(start, end []byte) (types.Iterator, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) CacheWrap() types.CacheWrap {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Commit() types.CommitID {
	//TODO implement me
	panic("implement me")
}

func (s *Store) LastCommitID() types.CommitID {
	//TODO implement me
	panic("implement me")
}
