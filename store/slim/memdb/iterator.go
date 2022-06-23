package memdb

import (
	"encoding/hex"
	"fmt"
	db "github.com/tendermint/tm-db"
	"sort"
	"time"
)

type PocketMemDBIterator struct {
	start     string
	end       string
	index     int
	isValid   bool
	isReverse bool
	size      int
	Ordered   []string
	M         map[string][]byte
}

var _ db.Iterator = &PocketMemDBIterator{}

func NewPocketMemDBIterator(startKey, endKey []byte, parent *PocketMemDB, isReverse bool) *PocketMemDBIterator {
	t := time.Now()
	startK := hex.EncodeToString(startKey)
	endK := hex.EncodeToString(endKey)
	it := &PocketMemDBIterator{
		start:     startK,
		end:       endK,
		isReverse: isReverse,
		isValid:   true,
		M:         make(map[string][]byte),
	}
	startIndex, endIndex := 0, 0
	parent.L.RLock()
	defer parent.L.RUnlock()
	if startKey != nil {
		startIndex = sort.SearchStrings(parent.Ordered, startK)
	}
	if endKey != nil {
		endIndex = sort.SearchStrings(parent.Ordered, endK)
	} else {
		endIndex = len(parent.Ordered)
	}
	window := parent.Ordered[startIndex:endIndex]
	windowSize := len(window)
	it.Ordered = make([]string, windowSize)
	if windowSize == 0 {
		it.isValid = false
		return it
	}
	it.size = windowSize
	if isReverse {
		it.index = it.size - 1
	}
	copy(it.Ordered, window)
	for _, k := range it.Ordered {
		it.M[k] = parent.M[k]
		//bs := parent.M[k] TODO may need to copy value slice for safety?
		//it.M[k] = make([]byte, len(bs))
		//copy(it.M[k], bs)
	}
	fmt.Println("iterator creation took: " + time.Since(t).String())
	return it
}

func (p *PocketMemDBIterator) Next() {
	if p.isReverse {
		p.index--
	} else {
		p.index++
	}
	if p.index >= p.size || p.index < 0 {
		p.isValid = false
	}
	return
}

func (p *PocketMemDBIterator) Key() (key []byte) {
	if !p.isValid {
		return nil
	}
	k := p.Ordered[p.index]
	b, _ := hex.DecodeString(k)
	return b
}

func (p *PocketMemDBIterator) Value() (value []byte) {
	if !p.isValid {
		return nil
	}
	k := p.Ordered[p.index]
	return p.M[k]
}

func (p *PocketMemDBIterator) Valid() bool { return p.isValid }
func (p *PocketMemDBIterator) Close()      { return }
func (p *PocketMemDBIterator) Domain() ([]byte, []byte) {
	panic("not implemented in PocketMemDBIterator")
}
func (p *PocketMemDBIterator) Error() error { panic("not implemented in PocketMemDBIterator") }
