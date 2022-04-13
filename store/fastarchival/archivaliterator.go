package fastarchival

import "github.com/tendermint/tm-db"

var _ db.Iterator = ArchivalIterator{}

// ArchivalIterator exists to make the iterator of archival capable of knowing its own suffix.
// this is used to answer the Key and Domain methods.
type ArchivalIterator struct {
	db.Iterator
}

func (r ArchivalIterator) Key() (key []byte) {
	return StoreKeySuffix(r.Iterator.Key())
}

func (r ArchivalIterator) Domain() ([]byte, []byte) {
	start, end := r.Iterator.Domain()
	return StoreKeySuffix(start), StoreKeySuffix(end)
}
