package memdb

import (
	"bytes"
	sdk "github.com/pokt-network/pocket-core/types"
	db "github.com/tendermint/tm-db"
	"github.com/tidwall/btree"
)

// PocketMemDBIterator is one of the trickier design points of the Persistence Replacement project. The main reason
// the iterator is such a difficulty - is the functional legacy the design must follow in order to not be consensus
// breaking (which is the requirement for V0 persistence to enable syncing from scratch).
//
// The legacy implementation uses an IAVL tree built with goleveldb iterators. The tricky part is to
// exactly match the behavior pattern when items of the Set are deleted, inserted, or modified during the iteration
// process ( a feature that pocket's utility layer uses in multiple places ).
//
// Alternative Implementations: The memdb implementation in the goleveldb library is the natural choice to satisfy these
// exact requirements however, the memdb is append only (memory is never shrunk even with deletes). This is obviously not
// a satisfactory solution to the problem. Tendermint's memdb iterator has a similar issue with locking and unlocking the
// entire db during iteration, which results in a deadlock. If you remove the lock, the db panics under certain conditions
// as the behavior for deleting during iteration is not properly implemented.
//
// Thus, a custom iterator implementation is chosen and modified to fit the exact functional requirements of the
// goleveldb iterators.

type PocketMemDBIterator struct {
	start     []byte                        // the start key
	end       []byte                        // the end key
	isValid   bool                          // if the iterator is 'valid' within the specified range
	isReverse bool                          // if the iterator is in reverse order or not
	parent    *PocketMemDB                  // the parent db
	it        btree.MapIter[string, []byte] // the base btree iterator
	prevKey   []byte                        // the previous key
	key       []byte                        // current iterator key
	value     []byte                        // current iterator value
}

var _ db.Iterator = &PocketMemDBIterator{}

func NewPocketMemDBIterator(startKey, endKey []byte, parent *PocketMemDB, isReverse bool) *PocketMemDBIterator {
	return (&PocketMemDBIterator{
		start:     startKey,
		end:       endKey,
		isValid:   true,
		isReverse: isReverse,
		parent:    parent,
		it:        parent.BTree.Iter(),
	}).Init()
}

func (p *PocketMemDBIterator) Seek() {
	p.prevKey = sdk.Cpy(p.key)
	if p.isReverse {
		p.seek(p.start, p.key, true)
	} else {
		p.seek(p.key, p.end, true)
	}
}

func (p *PocketMemDBIterator) seek(start, end []byte, initialized bool) {
	// if is reverse  we have two options  if endKey isn't nil, let's seek right to endKey. If endKey is
	// nil, let's seek to the last elem in the tree, but skip the end elem (end, start]. Similarly, if the
	// iterator isn't reverse, we have two options. if the start key isn't nil, seek to the key [start, end).
	// If startKey is nil, let's seek to the first item in the tree.
	if p.isReverse && end == nil {
		p.it.Last()
	} else if p.isReverse && end != nil {
		p.seekReverse(string(end), initialized)
	} else if !p.isReverse && start == nil {
		p.it.First()
	} else if !p.isReverse && start != nil {
		p.it.Seek(string(start))
	}
	p.Set()
}

func (p *PocketMemDBIterator) seekReverse(end string, initialized bool) {
	if !p.it.Seek(end) {
		p.it.Last()
	}
	p.SetKV()
	// if we are not in the initialization phase of the iterator (we've already initialized)
	// then we don't have to do the 'skip end' logic [start, end). Thus the comparator
	// is going to be '1' which is > rather than '0' which is just >=
	comparator := 0
	if initialized {
		comparator = 1
	}
	for bytes.Compare(p.key, []byte(end)) >= comparator { // while key >= end
		if !p.it.Prev() {
			p.isValid = false
			break
		}
		p.SetKV()
	}
}

func (p *PocketMemDBIterator) reSeek() (alreadyNexted bool) {
	// we re-initialize and re-seek the iterator everytime to perfectly match the goleveldb iterator logic
	// this comes with a significant overhead of O(log(n)), but at this point it's the best option without
	// modifying the utilty layer and possibly breaking consensus.
	p.it = p.parent.BTree.Iter()
	p.Seek()
	// seek already did next, this can happen when the Set is modified during iteration
	if p.prevKey != nil && !bytes.Equal(p.prevKey, p.key) {
		p.Set()
		return true
	}
	return false
}

func (p *PocketMemDBIterator) Next() {
	if p.reSeek() {
		return
	}
	if p.isReverse {
		if !p.it.Prev() {
			p.isValid = false
		}
	} else {
		if !p.it.Next() {
			p.isValid = false
		}
	}
	p.Set()
	return
}

func (p *PocketMemDBIterator) Set() {
	if !p.Valid() {
		return
	}
	p.SetKV()
	p.SetValidity()
}

// SetValidity : we always include <start> and exclude <end> regardless of isReverse
func (p *PocketMemDBIterator) SetValidity() {
	if p.KeyIsEmpty() {
		p.isValid = false
		return
	}
	if p.isReverse {
		if p.start == nil {
			return
		}
		if bytes.Compare(p.key, p.start) == -1 { // if key < start
			p.isValid = false
		}
	} else {
		if p.end == nil {
			return
		}
		if bytes.Compare(p.key, p.end) >= 0 { // if key <= end
			p.isValid = false
		}
	}
}

func (p *PocketMemDBIterator) SetKV()                     { p.key = []byte(p.it.Key()); p.value = p.it.Value() }
func (p *PocketMemDBIterator) Init() *PocketMemDBIterator { p.seek(p.start, p.end, false); return p }
func (p *PocketMemDBIterator) KeyIsEmpty() bool           { return p.key == nil || len(p.key) == 0 }
func (p *PocketMemDBIterator) Key() []byte                { return p.key }
func (p *PocketMemDBIterator) Value() []byte              { return p.value }
func (p *PocketMemDBIterator) Valid() bool                { return p.isValid }
func (p *PocketMemDBIterator) Close()                     { return /* no op */ }
func (p *PocketMemDBIterator) Domain() ([]byte, []byte)   { panic("not implemented") }
func (p *PocketMemDBIterator) Error() error               { panic("not implemented") }
