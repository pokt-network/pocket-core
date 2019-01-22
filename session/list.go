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

func GetSessionList() *List {
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

// "Add" adds a session to the global list
func (sList *List) Add(s Session) {
	(*types.List)(sList).Add(s.DevID, s)
}

// "Remove" removes a session from the global list
func (sList *List) Remove(s Session) {
	(*types.List)(sList).Remove(s.DevID)
}

// "Contains" returns whether or not it is within the list.
func (sList *List) Contains(devID string) bool {
	return (*types.List)(sList).Contains(devID)
}

// "Count" returns number of sessions
func (sList *List) Count() int {
	return (*types.List)(sList).Count()
}

// "Get" returns the number of sessions
func (sList *List) Get(devID string) Session {
	return (*types.List)(sList).Get(devID).(Session)
}

// "Print" prints the global sessionList
func (sList *List) Print() {
	(*types.List)(sList).Print()
}
