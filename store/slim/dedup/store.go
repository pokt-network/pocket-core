package dedup

import (
	"fmt"
	"github.com/pokt-network/pocket-core/store/types"
	db "github.com/tendermint/tm-db"
	"strconv"
	"strings"
)

var _ types.KVStore = &Store{}
var _ types.CommitStore = &Store{}

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

func HeightKey(height int64, prefix string, key []byte) []byte {
	return []byte(fmt.Sprintf("%d/%s/%s/", height, prefix, string(key)))
}

func FromHeightKey(heightKey string) (height int64, prefix string, key []byte) {
	var delim = "/"
	arr := strings.Split(heightKey, delim)
	// get height
	height, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		panic("unable to parse height from height key: " + heightKey)
	}
	prefix = arr[1]
	key = []byte(strings.Join(arr[2:], delim))
	return
}
