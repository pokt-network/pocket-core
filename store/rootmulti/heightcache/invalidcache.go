package heightcache

import (
	"errors"
	"github.com/pokt-network/pocket-core/store/types"
)

var _ types.SingleStoreCache = &InvalidCache{}

type InvalidCache struct {
}

func (i InvalidCache) Get(height int64, key []byte) ([]byte, error) {
	return nil, errors.New("invalid cache cannot get")
}

func (i InvalidCache) Has(height int64, key []byte) (bool, error) {
	return false, errors.New("invalid cache cannot get")
}

func (i InvalidCache) Set(key []byte, value []byte) {
}

func (i InvalidCache) Remove(key []byte) error {
	return errors.New("invalid cache cannot delete")
}

func (i InvalidCache) Iterator(height int64, start, end []byte) (types.Iterator, error) {
	return nil, errors.New("invalid cache has no iterators")
}

func (i InvalidCache) ReverseIterator(height int64, start, end []byte) (types.Iterator, error) {
	return nil, errors.New("invalid cache has no iterators")
}

func (i InvalidCache) Commit(height int64) {
}

func (i InvalidCache) Initialize(currentData map[string]string, version int64) {
}

func (i InvalidCache) IsValid() bool {
	return false
}
