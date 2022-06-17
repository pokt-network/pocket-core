package dedup

import (
	"github.com/pokt-network/pocket-core/types"
	dbm "github.com/tendermint/tm-db"
)

var _ types.Iterator = &DedupIterator{}
var _ types.KVStore = &Store{}
var _ types.CommitStore = &Store{}

type DedupIterator struct {
	parent dbm.GoLevelDB
	it     dbm.Iterator
}

func NewDedupIterator(parent dbm.GoLevelDB, height int64, prefix string, startKey, endKey []byte, isReverse bool) (dedupIterator *DedupIterator, err error) {
	start := HeightKey(height, prefix, startKey)
	end := HeightKey(height, prefix, endKey)
	dedupIterator = &DedupIterator{parent: parent}
	if isReverse {
		dedupIterator.it, err = parent.ReverseIterator(start, end)
	} else {
		dedupIterator.it, err = parent.Iterator(start, end)
	}
	return
}

func (d *DedupIterator) Next() {
	d.it.Next()
}

func (d *DedupIterator) Key() (key []byte) {
	return KeyFromHeightKey(d.it.Key())
}

func (d *DedupIterator) Value() (value []byte) {
	dataStoreKey := d.it.Value()
	if dataStoreKey == nil {
		return
	}
	value, err := d.parent.Get(dataStoreKey)
	if err != nil {
		panic("an error occurred in dedup iterator value call: " + err.Error())
	}
	return
}

func (d *DedupIterator) Error() error { return d.Error() }

func (d *DedupIterator) Close() { d.Close() }

func (d *DedupIterator) Valid() bool { return d.Valid() }

func (d *DedupIterator) Domain() (start []byte, end []byte) {
	st, end := d.it.Domain()
	return KeyFromHeightKey(st), KeyFromHeightKey(end)
}
