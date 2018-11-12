// This package is for all 'session' related code.
package session

import (
	"fmt"
	"github.com/pocket_network/pocket-core/node"
)

var (
	globalSessionList *sessionPool
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
This holds a list of list that are active (needs to confirm using liveness check).
 */
type sessionPool struct {
	// "list" is the local list of ongoing list.
	list map[string]Session
}

/*
 "GetSessionPoolInstance() returns the singleton instance of the global session pool
  TODO make thread safe
 */
func GetSessionPoolInstance() *sessionPool {
		if (globalSessionList == nil) {
			globalSessionList = &sessionPool{}
			globalSessionList.list = make(map[string]Session)
		}
	return globalSessionList
}

/*
"createNewSession" creates a new session for the specific devID and adds to global sessionPool (map)
 */
func CreateNewSession(dID string) {
	if(SearchSessionList(dID)==nil){
		// pulls the global list from the singleton
		sList := GetSessionPoolInstance().list
		// simulated List of Validators
		// TODO turn into real list of validators
		validators := make(map[string]node.Validator)
		// simulated List of Servicers
		// TODO turn into real list of servicers
		servicers := make(map[string]node.Service)
		// adds a new session to the sessionlist (map)
		sList[dID]=Session{dID,validators, servicers}
	}
}

/*
"SearchSessionList" searches the session list for the specific devID
 */
func SearchSessionList(dID string) *Session{
	// gets global session pool from singleton
	list := GetSessionPoolInstance()
	// pulls the session with the developer ID
	session:=list.list[dID]
	// if the session is found
	if session.devID!=""{
		return &session
	}
	return nil
}

/*
"PrintSessionList" prints the list from the session pool map"
 */
func PrintSessionList(){
	fmt.Println(GetSessionPoolInstance().list)
}
