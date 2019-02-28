package types

import (
	"fmt"
	"sync"
)

type Set struct {
	M   map[interface{}]struct{}
	Mux sync.Mutex
}

func NewSet() *Set {
	return &Set{M: make(map[interface{}]struct{})}
}

func (s *Set) Add(key interface{}) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	s.M[key] = struct{}{}
}

func (s *Set) Remove(key interface{}) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	delete(s.M, key)
}

func (s *Set) Count() int {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	return len(s.M)
}

func (s *Set) Contains(key interface{}) bool {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	_, ok := s.M[key]
	return ok
}

func (s *Set) Print() {
	fmt.Println(s.M)
}

func (s *Set) Clear() {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	s.M = map[interface{}]struct{}{}
}