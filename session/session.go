// This package is for all 'session' related code.
package session

import (
	"fmt"
	"github.com/pocket_network/pocket-core/node"
	"sync"
)

var (
	globalSessionList *sessionList
)
/*
This is the session structure.
 */
type Session struct {
	// "devID" is the developer's ID that identifies the session.
	devID string
	// "validators" is an array of validator nodes.
	validators []node.Validator
	// "servicers" is an array of service nodes.
	servicers []node.Service
}

/*
This holds a list of sessions that are active (needs to confirm using liveness check).
 */
type sessionList struct {
	// "sessions" is the local list of ongoing sessions.
	sessions map[string]Session
}

func GetSessionListInstance() *sessionList {
	sync.Once{}.Do(func() {
		if (globalSessionList == nil) {
			globalSessionList = &sessionList{}
			globalSessionList.sessions = make(map[string]Session)
		}
	})
	return globalSessionList
}

/*
"createNewSession" creates a new session for the specific devID and adds to global sessionList (map)
 */
func CreateNewSession(dID string) {
	if(SearchSessionList(dID)==nil){
		// pulls the global list from the singleton
		sList :=GetSessionListInstance().sessions
		// simulated List of Validators
		// TODO turn into real list of validators
		validators :=[]node.Validator{}
		// simulated List of Servicers
		// TODO turn into real list of servicers
		servicers :=[]node.Service{}
		// adds a new session to the sessionlist (map)
		sList[dID]=Session{dID,validators, servicers}
	}
}

/*
"searchSessionList" searches the session list for the specific devID
 */
func SearchSessionList(dID string) *Session{
	// gets global session list from singleton
	list := GetSessionListInstance()
	// pulls the session with the developer ID
	session:=list.sessions[dID]
	// if the session is found
	// TODO may be an error here
	if &session!=nil{
		fmt.Println("Session Found!")
		return &session
	}
	// else return nil
	fmt.Println("Session Not Found!")
	return nil
}
