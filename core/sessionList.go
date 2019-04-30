package core

import (
	"github.com/pokt-network/pocket-core/types"
	"sync"
)

type SessionList types.List // [developer ID] -> session object

var (
	globalSessionList *SessionList
	sessionListOnce   sync.Once
)

func GetSessionList() *SessionList {
	sessionListOnce.Do(func() {
		globalSessionList = (*SessionList)(types.NewList())
	})
	return globalSessionList
}

func (sl *SessionList) AddSession(s Session) {
	(*types.List)(sl).Add(s.DevID, s)
}

func (sl *SessionList) RemoveSession(s Session) {
	(*types.List)(sl).Remove(s.DevID)
}

func (sl *SessionList) Clear() {
	(*types.List)(sl).Clear()
}

func (sl *SessionList) GetSession(devID string) Session {
	return (*types.List)(sl).Get(devID).(Session)
}

func (sl *SessionList) Contains(devID string) bool {
	return (*types.List)(sl).Contains(devId)
}

func (sl *SessionList) Print() {
	(*types.List)(sl).Print()
}

func (sl *SessionList) Len() int {
	return (*types.List)(sl).Count()
}
