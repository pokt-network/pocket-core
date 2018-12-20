// This package is for all 'session' related code.
package sessio
// TODO thread safety
import (
	"fmt"
	"sync"
)

var (
	globalSessionPool *sessionPool // global session pool instance
	o                 sync.Once    // only occurs o throughout program
	lock              sync.Mutex   // for thread safety // TODO consider making a member of the sessionPool
)

/*
 "GetSessionPoolInstance() returns the singleton instance of the global session pool
 */
func GetSessionPoolInstance() *sessionPool {
	o.Do(func() { // only do o
		if 	globalSessionPool == nil { 					  		// if no existing globalSessionPool
			globalSessionPool = &sessionPool{}                	// create a new session pool
			globalSessionPool.list = make(map[string]Session) 	// create a map of sessions
		}
	})
	return globalSessionPool // return the session pool
}

/*
"createNewSession" creates a new session for the specific devID and adds to global sessionPool (map)
 */
func CreateAndRegisterSession(dID string) {
	RegisterSession(NewSession(dID))
}

func RegisterSession(session Session){
	lock.Lock()
	defer lock.Unlock()
	if globalSessionPool==nil{
		GetSessionPoolInstance()
	}
	if !sessionListContains(session.devID) {
		sList := GetSessionPoolInstance().list           							// pulls the global list from the singleton
		sList[session.devID] = session												// adds a new session to the sessionlist (map)
	}
}

func NewSession(dID string) Session{
	return Session{dID,make(map[string]Connection)}
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
	if globalSessionPool==nil{				// TODO check if this is necessary in sessiogo and peers.go
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

func (session *Session) GetSessionConnList() map[string]Connection {
	once.Do(func() {
		session.connectionList = make(map[string]Connection)
	})
	return session.connectionList
}

func (session *Session) RegisterSessionConn(connection Connection) {
	fmt.Println("REGISTERING CONN "+ connection.Conn.RemoteAddr().String())
	session.GetSessionConnList()[connection.Conn.RemoteAddr().String()] = connection // added by remote addr
}

func (session *Session) ClearSessionConnList(){
	session.GetSessionConnList() = make(map[string]Connection)
}
