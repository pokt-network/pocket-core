package types

import (
	"fmt"
	"sync"
)

// TODO replace redundant code with this general type
type List struct {
	M   map[interface{}]interface{}
	Mux sync.Mutex
}

func NewList() *List {
	return &List{M: make(map[interface{}]interface{})}
}

func (l *List) Add(key interface{}, val interface{}) {
	l.Mux.Lock()
	defer l.Mux.Unlock()
	l.M[key] = val
}

func (l *List) Remove(key interface{}) {
	l.Mux.Lock()
	defer l.Mux.Unlock()
	delete(l.M, key)
}

func (l *List) Contains(key interface{}) bool {
	l.Mux.Lock()
	defer l.Mux.Unlock()
	_, ok := l.M[key]
	return ok
}

func (l *List) Print() {
	fmt.Println(l.M)
}
