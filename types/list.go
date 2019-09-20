package types

import (
	"sync"
)

type List struct {
	M   map[interface{}]interface{}
	Mux sync.Mutex
}

func NewList() *List {
	return &List{M: make(map[interface{}]interface{})}
}

func (l *List) Add(key, val interface{}) {
	l.Mux.Lock()
	defer l.Mux.Unlock()
	l.M[key] = val
}

func (l *List) Remove(key interface{}) {
	l.Mux.Lock()
	defer l.Mux.Unlock()
	delete(l.M, key)
}

func (l *List) Get(key interface{}) interface{} {
	l.Mux.Lock()
	defer l.Mux.Unlock()
	return l.M[key]
}

func (l *List) Count() int {
	l.Mux.Lock()
	defer l.Mux.Unlock()
	return len(l.M)
}

func (l *List) Contains(key interface{}) bool {
	l.Mux.Lock()
	defer l.Mux.Unlock()
	_, ok := l.M[key]
	return ok
}

func (l *List) Clear() {
	l.Mux.Lock()
	defer l.Mux.Unlock()
	l.M = map[interface{}]interface{}{}
}
