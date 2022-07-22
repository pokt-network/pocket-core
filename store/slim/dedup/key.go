package dedup

import (
	"github.com/jordanorelli/lexnum"
	dbm "github.com/tendermint/tm-db"
)

var elenEncoder = lexnum.NewEncoder('=', '-')

func GetLatestHeightKeyRelativeToQuery(prefix string, exclusiveEndHeight, key []byte, parentDB dbm.DB) ([]byte, bool) {
	// use a reverse iterator's 'seek' functionality in order to track down the latest height key
	// relative to the query height. Note: this only works because we use ELEN encoding for int64
	startKey := HeightKey(prefix, SearchStartHeight, key)
	endKey := HeightKey(prefix, exclusiveEndHeight, key)
	it, _ := parentDB.ReverseIterator(startKey, endKey)
	defer it.Close()
	if it.Valid() {
		return it.Key(), true
	}
	return nil, false
}

func KeyForExists(prefix string, k []byte) []byte {
	return append([]byte(ExistsPrefix+prefix), k...)
}

func KeyFromExistsKey(prefix string, k []byte) []byte {
	return k[len(ExistsPrefix+prefix):]
}

func HeightKey(prefix string, latestHeightString []byte, key []byte) []byte {
	return append(append([]byte(prefix), key...), latestHeightString...)
}

// util

var (
	DeletedValue      = []byte("del")
	SearchStartHeight = []byte("0")
)

const (
	ExistsPrefix = "ex"
)
