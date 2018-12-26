// This package is for all 'session' related code.
package sessio

import (
	"fmt"
	"sync"
)

/*
This is the session structure.
 */
type Session struct {
	DevID    string                	`json:"devid"` 			// "DevID" is the developer's ID that identifies the sessio
	ConnList map[string]Connection 	`json:"connList"`		// "ConnList" is the List of peer connections [GID] Connection
	sync.Mutex						`json:"mutex"`
}


/***********************************************************************************************************************
Session Pool Code // TODO naming convention ' sessionPool' or 'sessionpool'
 */

var (
	globalSessionPool *sessionPool // global session pool instance
	sPoolLock         sync.Mutex   // for thread safety
	once			  sync.Once
)

/*
 "GetSessionPoolInstance() returns the singleton instance of the global session pool
 */
func GetSessionPoolInstance() *sessionPool {
	once.Do(func() { // only do o
		if 	globalSessionPool == nil { 					  		// if no existing globalSessionPool
			globalSessionPool = &sessionPool{}                	// create a new session pool
			globalSessionPool.List = make(map[string]Session) 	// create a map of sessions
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
		sList := GetSessionPoolInstance().List // pulls the global list from the singleton
		sList[session.DevID] = session         // adds a new session to the sessionlist (map)
	}
}

/*
"sessionListContains" searches the session list for the specific devID and returns whether or not it is held
 */
func sessionListContains(dID string) bool{
	_,ok := GetSessionPoolInstance().List[dID]
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
	_,ok := GetSessionPoolInstance().List[dID]
	return ok
}

func NewSession(dID string) Session{
	return Session {DevID: dID}
}

/***********************************************************************************************************************
Session Methods
 */
func (session *Session) GetConnections() map[string]Connection {
	var once sync.Once
	once.Do(func() {
		if session.ConnList == nil {
			session.ConnList = make(map[string]Connection)
		}
	})
	return session.ConnList
}

func (session *Session) AddConnection(connection Connection) {
	fmt.Println("Adding Connection: " + connection.Peer.(SessionPeer).GID + " to Session: " + session.DevID + "")
	session.Lock()
	defer session.Unlock()
	session.GetConnections()[connection.Peer.(SessionPeer).GID] = connection
}

func (session *Session) GetConnection(gid string) Connection {
	session.Lock()
	defer session.Unlock()
	sessConnList := session.GetConnections()
	return sessConnList[gid]
}

func (session *Session) GetConnectionByIP(ip string) Connection {
	session.Lock()
	defer session.Unlock()
	fmt.Println(session.GetConnections())
	fmt.Println(ip)
	connection := session.GetConnections()[ip]
	if new(Connection) == &connection {
		panic("Unable to locate connection from IP")
	}
	return connection
}

func (session *Session) NewConnections(sp []SessionPeer) {
	for _, sessionPeer := range sp {	// TODO dial each individual sessionPeer -> add connection to connList
		connection := Connection{Peer: sessionPeer}
		//connection.Dial("3333", sessionPeer.LocalIP) // TODO allow flexible ports: startup flag -> in Node structure -> in Session Peer Structure -> called here
		session.AddConnection(connection)
	}
}

func (session *Session) ClearConnections() {
	session.Lock()
	defer session.Unlock()
	session.ConnList = make(map[string]Connection)
}
