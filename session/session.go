// This package is for all 'session' related code.
package session
// TODO thread safety
import (
	"fmt"
	"github.com/pokt-network/pocket-core/node"
	"sync"
)

var (
	globalSessionPool *sessionPool							// global session pool instance
	once              sync.Once								// only occurs once throughout program
	lock			  sync.Mutex							// for thread safety
)
/*
This is the session structure.
 */
type Session struct {
	devID string											// "devID" is the developer's ID that identifies the session.
	validators map[string]node.Validator					// "validators" is a map [devid]Node validator nodes.
	servicers map[string]node.Service						// "validators" is a map [devid]Node servicer nodes.
}

/*
This holds a list of list that are active (needs to confirm using liveness check).
 */
type sessionPool struct {
	list map[string]Session // "list" is the local list of ongoing list.
}

/*
 "GetSessionPoolInstance() returns the singleton instance of the global session pool
 */
func GetSessionPoolInstance() *sessionPool {
	once.Do(func() { 										  	// only do once
		if (globalSessionPool == nil) { 					  	// if no existing globalSessionPool
			globalSessionPool = &sessionPool{}                	// create a new session pool
			globalSessionPool.list = make(map[string]Session) 	// create a map of sessions
		}
	})
	return globalSessionPool // return the session pool
}

/*
"createNewSession" creates a new session for the specific devID and adds to global sessionPool (map)
 */
func CreateNewSession(dID string) {
	lock.Lock()
	defer lock.Unlock()
	if globalSessionPool==nil{
		GetSessionPoolInstance()
	}
	if !sessionListContains(dID) {
		sList := GetSessionPoolInstance().list           	// pulls the global list from the singleton
		validators := make(map[string]node.Validator)    	// simulated List of Validators
		servicers := make(map[string]node.Service)       	// simulated List of Servicers
		sList[dID] = Session{dID, validators, servicers} // adds a new session to the sessionlist (map)
	}
}

/*
"sessionListContains" searches the session list for the specific devID and returns whether or not it is held
<<<<<<< 4b882753294f034021971c577be6c5ff147314a1
 */
func sessionListContains(dID string) bool{
	_,ok := GetSessionPoolInstance().list[dID]
	return ok
}

/*
"SessionListContains" searches the session list for the specific devID and returns whether or not it is held
Thread safe
 */
func SessionListContains(dID string) bool{
	lock.Lock()
	defer lock.Unlock()
	if globalSessionPool==nil{				// TODO check if this is necessary in session.go and peers.go
		GetSessionPoolInstance()
	}
	_,ok := GetSessionPoolInstance().list[dID]
	return ok
}

/*
"PrintSessionList" prints the list from the session pool map"
 */
func PrintSessionList() {
	fmt.Println(GetSessionPoolInstance().list)
}
