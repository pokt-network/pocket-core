package memdb

import (
	dbm "github.com/tendermint/tm-db"
	"github.com/tidwall/btree"
	"sync"
)

var _ dbm.DB = &PocketMemDB{}

type PocketMemDB struct {
	L     sync.RWMutex
	BTree *btree.Map[string, []byte]
}

func NewPocketMemDB() *PocketMemDB {
	return &PocketMemDB{
		L:     sync.RWMutex{},
		BTree: &btree.Map[string, []byte]{},
	}
}

func (p *PocketMemDB) Get(key []byte) ([]byte, error) {
	p.L.RLock()
	defer p.L.RUnlock()
	v, _ := p.BTree.Get(string(key))
	return v, nil
}

func (p *PocketMemDB) Has(key []byte) (bool, error) {
	p.L.RLock()
	defer p.L.RUnlock()
	_, found := p.BTree.Get(string(key))
	return found, nil
}

func (p *PocketMemDB) Set(key []byte, value []byte) error {
	p.L.Lock()
	defer p.L.Unlock()
	p.BTree.Set(string(key), value)
	return nil
}

func (p *PocketMemDB) Delete(key []byte) error {
	p.L.Lock()
	defer p.L.Unlock()
	p.BTree.Delete(string(key))
	return nil
}

func (p *PocketMemDB) Iterator(start, end []byte) (dbm.Iterator, error) {
	return NewPocketMemDBIterator(start, end, p, false), nil
}

func (p *PocketMemDB) ReverseIterator(start, end []byte) (dbm.Iterator, error) {
	return NewPocketMemDBIterator(start, end, p, true), nil
}

func (p *PocketMemDB) NewBatch() dbm.Batch          { return NewPocketMemDBBatch(p) }
func (p *PocketMemDB) SetSync([]byte, []byte) error { panic("not implemented in PocketMemDB") }
func (p *PocketMemDB) DeleteSync([]byte) error      { panic("not implemented in PocketMemDB") }
func (p *PocketMemDB) Close() error                 { panic("not implemented in PocketMemDB") }
func (p *PocketMemDB) Print() error                 { panic("not implemented in PocketMemDB") }
func (p *PocketMemDB) Stats() map[string]string     { panic("not implemented in PocketMemDB") }
