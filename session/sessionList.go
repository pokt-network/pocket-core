package session

import (
	"sync"
	
	"github.com/pokt-network/pocket-core/types"
)

type List types.List

var (
	sesList *List
	sListO  sync.Once
)

func GetSList() *List {
	sListO.Do(func() {
		sesList = (*List)(types.NewList())
	})
	return sesList
}

// "AddSession" adds a session object to the global list
func (sList *List) AddMulti(s ...Session) {
	sList.Mux.Lock()
	defer sList.Mux.Unlock()
	for _, session := range s {
		sList.M[session.DevID] = session
	}
}

// "RemoveSession" removes a session object from the global list
func (sList *List) RemoveSession(s ...Session) {
	sList.Mux.Lock()
	defer sList.Mux.Unlock()
	for _, session := range s {
		delete(sList.M, session.DevID)
	}
}
