package memdb

import (
	"encoding/hex"
	dbm "github.com/tendermint/tm-db"
	"sort"
	"sync"
)

var _ dbm.DB = &PocketMemDB{}

type PocketMemDB struct {
	L       sync.RWMutex
	M       map[string][]byte
	Ordered []string
}

func NewPocketMemDB() *PocketMemDB {
	return &PocketMemDB{
		L:       sync.RWMutex{},
		M:       make(map[string][]byte),
		Ordered: make([]string, 0),
	}
}

func (p *PocketMemDB) Get(key []byte) ([]byte, error) {
	p.L.RLock()
	defer p.L.RUnlock()
	return p.M[hex.EncodeToString(key)], nil
}

func (p *PocketMemDB) Has(key []byte) (bool, error) {
	p.L.RLock()
	defer p.L.RUnlock()
	_, ok := p.M[hex.EncodeToString(key)]
	return ok, nil
}

func (p *PocketMemDB) Set(key []byte, value []byte) error {
	k := hex.EncodeToString(key)
	p.L.Lock()
	defer p.L.Unlock()
	p.M[k] = value
	p.InsertOrdered(k)
	return nil
}

func (p *PocketMemDB) Delete(key []byte) error {
	if ok, _ := p.Has(key); !ok {
		return nil
	}
	k := hex.EncodeToString(key)
	p.L.Lock()
	defer p.L.Unlock()
	delete(p.M, k)
	p.DeleteOrdered(k)
	return nil
}

func (p *PocketMemDB) Iterator(start, end []byte) (dbm.Iterator, error) {
	return NewPocketMemDBIterator(start, end, p, false), nil
}

func (p *PocketMemDB) ReverseIterator(start, end []byte) (dbm.Iterator, error) {
	return NewPocketMemDBIterator(start, end, p, true), nil
}

func (p *PocketMemDB) NewBatch() dbm.Batch {
	return NewPocketMemDBBatch(p)
}

// InsertOrdered : Contract, the structure is WLocked
func (p *PocketMemDB) InsertOrdered(key string) {
	i := sort.SearchStrings(p.Ordered, key)
	p.Ordered = append(p.Ordered, "")
	copy(p.Ordered[i+1:], p.Ordered[i:])
	p.Ordered[i] = key
	return
}

// DeleteOrdered : Contract, the structure is WLocked & the element is present in the slice
func (p *PocketMemDB) DeleteOrdered(key string) {
	x := sort.SearchStrings(p.Ordered, key)
	p.Ordered = append(p.Ordered[:x], p.Ordered[x+1:]...)
	return
}

func (p *PocketMemDB) SetSync([]byte, []byte) error { panic("not implemented in PocketMemDB") }
func (p *PocketMemDB) DeleteSync([]byte) error      { panic("not implemented in PocketMemDB") }
func (p *PocketMemDB) Close() error                 { panic("not implemented in PocketMemDB") }
func (p *PocketMemDB) Print() error                 { panic("not implemented in PocketMemDB") }
func (p *PocketMemDB) Stats() map[string]string     { panic("not implemented in PocketMemDB") }
