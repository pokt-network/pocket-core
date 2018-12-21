package sessio

import (
	"fmt"
	"sync"
)

var (
	sList *sessionPool // global session pool instance
)

/***********************************************************************************************************************
Singleton getter
 */
func GetSessionList() *sessionPool {
	var once sync.Once
	once.Do(func() { // only do once
		if 	sList == nil { // if no existing sList
			sList = &sessionPool{}                // create a new session pool
			sList.List = make(map[string]Session) // create a map of sessions
		}
	})
	return sList // return the session pool
}

/***********************************************************************************************************************
sessionPool Methods
 */
func (sList *sessionPool) AddSession(session Session){ // this is a function because only 1 global sessionPool instance
	if !sList.Contains(session.DevID) {
		sList.Lock()
		defer sList.Unlock()
		sList.List[session.DevID] = session // adds a new session to the sessionlist (map)
	}
}

func (sList *sessionPool) RemoveSession(session Session){
	sList.Lock()
	defer sList.Unlock()
	delete(sList.List, session.DevID)
}

func (sList *sessionPool) Contains(dID string) bool{
	sList.Lock()
	defer sList.Unlock()
	_,ok := sList.List[dID]
	return ok
}

func (pList *sessionPool) Count() int{
	pList.Lock()
	defer pList.Unlock()
	return len(pList.List)
}

func (sList *sessionPool) Print() {
	fmt.Println(sList.List)
}

/***********************************************************************************************************************
sessionPool Functions
 */

func GetSessionCountt() int{
	sList := GetSessionList()
	sList.Lock()
	defer sList.Unlock()
	return len(sList.List)
}
