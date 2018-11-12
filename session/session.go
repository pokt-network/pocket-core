// This package is for all 'session' related code.
package session

import (
	"fmt"
	"github.com/pocket_network/pocket-core/node"
)

var (
	globalSessionList *sessionList
)
/*
This is the session structure.
 */
 // TODO consider converting validators and servicers to map
type Session struct {
	// "devID" is the developer's ID that identifies the session.
	devID string
	// "validators" is a map [devid]Node validator nodes.
	validators map[string]node.Validator

	// "validators" is a map [devid]Node servicer nodes.
	servicers  map[string]node.Service
}

/*
This holds a list of sessions that are active (needs to confirm using liveness check).
 */
type sessionList struct {
	// "sessions" is the local list of ongoing sessions.
	sessions map[string]Session
}

/*
 "GetSessionListInstance() returns the singleton instance of the global session list
  TODO make thread safe
 */
func GetSessionListInstance() *sessionList {
		if (globalSessionList == nil) {
			globalSessionList = &sessionList{}
			globalSessionList.sessions = make(map[string]Session)
		}
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
"SearchSessionList" searches the session list for the specific devID
 */
func SearchSessionList(dID string) *Session{
	// gets global session list from singleton
	list := GetSessionListInstance()
	// pulls the session with the developer ID
	session:=list.sessions[dID]
	// if the session is found
	if session.devID!=""{
		return &session
	}
	return nil
}

/*
"PrintSessionList" prints the session list map"
 */
func PrintSessionList(){
	fmt.Println(GetSessionListInstance().sessions)
}
