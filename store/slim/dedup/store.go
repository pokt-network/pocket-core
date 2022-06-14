package dedup

import (
	"bytes"
	"github.com/pokt-network/pocket-core/store/types"
	db "github.com/tendermint/tm-db"
)

type Store struct {
	Height                int64
	HeightString          []byte
	ExclusiveHeightString []byte
	Prefix                string
	ParentDB              db.DB
}

func NewStore(height int64, prefix string, parent db.DB) *Store {
	heightString := []byte(elenEncoder.EncodeInt(int(height)))
	return &Store{
		Height:                height,
		HeightString:          heightString,
		ExclusiveHeightString: types.PrefixEndBytes(heightString),
		Prefix:                prefix,
		ParentDB:              parent,
	}
}

// reads

func (s *Store) Get(k []byte) ([]byte, error) {
	latestKey, found := s.GetLatestHeightKeyRelativeToQuery(k)
	if !found {
		return nil, nil
	}
	value, err := s.ParentDB.Get(latestKey)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(value, DeletedValue) {
		return nil, nil
	}
	return value, nil
}

func (s *Store) Has(k []byte) (bool, error) {
	k, err := s.Get(k)
	return k != nil, err
}

func (s *Store) GetLatestHeightKeyRelativeToQuery(k []byte) ([]byte, bool) {
	return GetLatestHeightKeyRelativeToQuery(s.Prefix, s.ExclusiveHeightString, k, s.ParentDB)
}

func (s *Store) Iterator(start, end []byte) (types.Iterator, error) {
	return NewCascadeIterator(s.ParentDB, s.Height, s.Prefix, s.ExclusiveHeightString, start, end, false)
}

func (s *Store) ReverseIterator(start, end []byte) (types.Iterator, error) {
	return NewCascadeIterator(s.ParentDB, s.Height, s.Prefix, s.ExclusiveHeightString, start, end, true)
}

// writes

func (s *Store) Set(k, value []byte) error {
	heightKey := HeightKey(s.Prefix, s.HeightString, k)
	existsKey := KeyForExists(s.Prefix, k)
	if err := s.ParentDB.Set(existsKey, nil); err != nil {
		return err
	}
	if err := s.ParentDB.Set(heightKey, value); err != nil {
		return err
	}
	return nil
}

func (s *Store) Delete(k []byte) error {
	heightKey := HeightKey(s.Prefix, s.HeightString, k)
	if err := s.ParentDB.Set(heightKey, DeletedValue); err != nil {
		return err
	}
	return nil
}

func (s *Store) Commit() {
	s.Height++
	s.HeightString = []byte(elenEncoder.EncodeInt(int(s.Height)))
	s.ExclusiveHeightString = types.PrefixEndBytes(s.HeightString)
	return
}

// unused below

func (s *Store) CacheWrap() types.CacheWrap   { panic("cachewrap not implemented for de-dup store") }
func (s *Store) LastCommitID() types.CommitID { panic("lastCommitID not implemented for de-dup store") }
