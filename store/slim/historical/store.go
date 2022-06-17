package historical

import (
	"github.com/jackc/pgx/v4"
	"github.com/pokt-network/pocket-core/store/types"
)

var _ types.KVStore = &Store{}
var _ types.CommitStore = &Store{}

type Store struct {
	Height int64
	Prefix string
	DB     *pgx.Conn
}

func NewStore(height int64, prefix string, postgresConnection *pgx.Conn) Store {
	return Store{
		Height: height,
		Prefix: prefix,
		DB:     postgresConnection,
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
