// This package is network code relating to pocket 'sessions'
package session

import (
	"fmt"
	"github.com/pokt-network/pocket-core/logs"
	"sync"
)

// "sessionList.go" holds the sessionList structure, methods and functions.

/*
This holds a List of List that are active (needs to confirm using liveness check).
*/
type sessionList struct {
	List       map[string]Session 						// "List" is the local List of ongoing List.
	sync.Mutex                    						// for thread safety
}

var (
	sList *sessionList	 								// global session pool instance
	once sync.Once
)
// TODO consider abstracting list type
/***********************************************************************************************************************
Singleton getter
*/

/*
"GetSessionList" returns the global sessionList object
 */
func GetSessionList() *sessionList {
	once.Do(func() { 									// only do once
		sList = &sessionList{}                		// create a new session pool
		sList.List = make(map[string]Session) 		// create a map of sessions
	})
	return sList 										// return the session pool
}

/***********************************************************************************************************************
sessionList Methods
*/

/*
"AddSession" adds a session object to the global list
 */
func (sList *sessionList) AddSession(s... Session) { 	// this is a function because only 1 global sessionList instance
	sList.Lock()										// lock the list
	defer sList.Unlock()								// unlock after complete
	for _, session := range s {							// for each session
		logs.NewLog("New session added to list: "+session.DevID, logs.InfoLevel, logs.JSONLogFormat)
		sList.List[session.DevID] = session 			// adds a new session to the sessionlist (map)
	}
}

/*
"RemoveSession" removes a session object from the global list
 */
func (sList *sessionList) RemoveSession(s... Session) {
	sList.Lock()										// locks the list
	defer sList.Unlock()								// unlock after complete
	for _,session := range s {
		logs.NewLog("Session "+session.DevID+" removed from list", logs.InfoLevel, logs.JSONLogFormat)
		delete(sList.List, session.DevID) 				// delete from list
	}
}

/*
"Contains" checks to see if a session is within the global list
 */
func (sList *sessionList) Contains(dID string) bool {
	sList.Lock()										// lock the list
	defer sList.Unlock()								// unlock the list
	_, ok := sList.List[dID]							// check if contains
	return ok											// return the bool
}

/*
"Count" returns the number of sessions within the global list
 */
func (sList *sessionList) Count() int {
	sList.Lock()										// lock the list
	defer sList.Unlock()								// unlock the list
	return len(sList.List)								// return the length of the global list
}

/*
"Print" prints the global session list to the CLI
 */
func (sList *sessionList) Print() {
	fmt.Println(sList.List)								// print to the CLI
}

/*
"Get" returns a session from the list based on the developer ID
 */
func (sList *sessionList) Get(dID string) Session{
	sList.Lock()
	defer sList.Unlock()
	return sList.List[dID]
}

/*
"Set" updates a session within the list based on the developer ID
 */
 func (sList *sessionList) Set(dID string, s Session){
 	sList.Lock()
 	defer sList.Unlock()
 	sList.List[dID]= s
 }

/***********************************************************************************************************************
sessionList Functions
*/

/*
"GetSessionCount" returns the number of sessions without the global sessionList object
 */
func GetSessionCount() int {
	sList := GetSessionList()							// get the list
	sList.Lock()										// lock the list
	defer sList.Unlock()								// after unlock the list
	return len(sList.List)								// return the length of the global list
}
