// This package is for all 'session' related code.
package sessio
import (
	"fmt"
	"github.com/pokt-network/pocket-core/node"
	"sync"
)

/***********************************************************************************************************************
Session Pool Code // TODO naming convention ' sessionList' or 'sessionpool'
 */

var (
	globalSessionPool *sessionPool // global session pool instance
	sPoolLock         sync.Mutex   // for thread safety
)

/*
 "GetSessionPoolInstance() returns the singleton instance of the global session pool
 */
func GetSessionPoolInstance() *sessionPool {
	sync.Once{}.Do(func() { // only do o
		if 	globalSessionPool == nil { 					  		// if no existing globalSessionPool
			globalSessionPool = &sessionPool{}                	// create a new session pool
			globalSessionPool.list = make(map[string]Session) 	// create a map of sessions
		}
	})
	return globalSessionPool // return the session pool
}

func CreateAndRegisterSession(dID string) {
	RegisterSession(NewSession(dID))
}

func RegisterSession(session Session){ // this is a function because only 1 global sessionPool instance
	sPoolLock.Lock()
	defer sPoolLock.Unlock()
	if globalSessionPool==nil{
		GetSessionPoolInstance()
	}
	if !sessionListContains(session.DevID) {
		sList := GetSessionPoolInstance().list // pulls the global list from the singleton
		sList[session.DevID] = session         // adds a new session to the sessionlist (map)
	}
}

/*
"sessionListContains" searches the session list for the specific devID and returns whether or not it is held
 */
func sessionListContains(dID string) bool{
	_,ok := GetSessionPoolInstance().list[dID]
	return ok
}

/*
"SessionListContains" searches the session list for the specific DevID and returns whether or not it is held
Thread safe
 */
func SessionListContains(dID string) bool{
	sPoolLock.Lock()
	defer sPoolLock.Unlock()
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

/***********************************************************************************************************************
Session Code
 */

func NewSession(dID string) Session{ // TODO will derive peers for session using blockchain
	connList:=make(map[string]Connection)
	// create dummy nodes
	sNode1 := node.Node{GID: "snode1", RemoteIP: "localhost", LocalIP: "localhost"}
	sNode2 := node.Node{GID: "snode2", RemoteIP: "localhost", LocalIP: "localhost"}
	vNode := node.Node{GID:"vnode", RemoteIP:"localhost", LocalIP:"localhost"}
	// create dummy connections
	sNodeConn1:= Connection{Mutex: sync.Mutex{}, Node: sNode1, Role: SERVICER}
	sNodeConn2:= Connection{Mutex: sync.Mutex{}, Node: sNode2, Role: SERVICER}
	vNodeConn:= Connection{Mutex:sync.Mutex{}, Node: vNode, Role: VALIDATOR}
	// add to list
	connList[sNode1.GID]=sNodeConn1
	connList[sNode2.GID]=sNodeConn2
	connList[vNode.GID]=vNodeConn
	// return resulting session
	return Session{dID,connList, sync.Mutex{}}
}

func (session *Session) GetSessionConnList() map[string]Connection {
	sync.Once{}.Do(func() {
		session.ConnList = make(map[string]Connection)
	})
	return session.ConnList
}

func (session *Session) RegisterSessionConn(connection Connection) {
	fmt.Println("REGISTERING CONN "+ connection.Conn.RemoteAddr().String())
	session.Lock()
	defer session.Unlock()
	session.GetSessionConnList()[connection.Conn.RemoteAddr().String()] = connection // added by remote addr
}

func (session *Session) GetConnectionFromList(gid string) Connection {
	session.Lock()
	defer session.Unlock()
	sessConnList:= session.GetSessionConnList()
	return sessConnList[gid]
}

func (session *Session) ClearSessionConnList(){
	session.Lock()
	defer session.Unlock()
	session.ConnList = make(map[string]Connection)
}
