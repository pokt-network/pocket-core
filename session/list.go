// This package is network code relating to pocket 'sessions'
package session

import (
	"fmt"
	"sync"

	"github.com/pokt-network/pocket-core/logs"
)

type sessionList struct {
	List map[string]Session // [DID]
	sync.Mutex
}

var (
	sList *sessionList
	once  sync.Once
)

// "GetSessionList" returns the global sessionList object
func GetSessionList() *sessionList {
	once.Do(func() {
		sList = &sessionList{}
		sList.List = make(map[string]Session)
	})
	return sList
}

// "AddSession" adds a session object to the global list
func (sList *sessionList) AddSession(s ...Session) {
	sList.Lock()
	defer sList.Unlock()
	for _, session := range s {
		logs.NewLog("New session added to list: "+session.DevID, logs.InfoLevel, logs.JSONLogFormat)
		sList.List[session.DevID] = session
	}
}

// "RemoveSession" removes a session object from the global list
func (sList *sessionList) RemoveSession(s ...Session) {
	sList.Lock()
	defer sList.Unlock()
	for _, session := range s {
		logs.NewLog("Session "+session.DevID+" removed from list", logs.InfoLevel, logs.JSONLogFormat)
		delete(sList.List, session.DevID)
	}
}

// "Contains" checks to see if a session is within the global list
func (sList *sessionList) Contains(dID string) bool {
	sList.Lock()
	defer sList.Unlock()
	_, ok := sList.List[dID]
	return ok
}

// "Count" returns the number of sessions within the global list
func (sList *sessionList) Count() int {
	sList.Lock()
	defer sList.Unlock()
	return len(sList.List)
}

// "Print" prints the global session list to the CLI
func (sList *sessionList) Print() {
	fmt.Println(sList.List) // print to the CLI
}

// "Get" returns a session from the list based on the developer ID
func (sList *sessionList) Get(dID string) Session {
	sList.Lock()
	defer sList.Unlock()
	return sList.List[dID]
}

/*
"Set" updates a session within the list based on the developer ID
*/
func (sList *sessionList) Set(dID string, s Session) {
	sList.Lock()
	defer sList.Unlock()
	sList.List[dID] = s
}
