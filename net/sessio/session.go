// This package is for all 'session' related code.
package sessio
import (
	"fmt"
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/node"
	"net"
	"sync"
)
/***********************************************************************************************************************
Session Pool Code // TODO naming convention ' sessionList' or 'sessionpool'
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

func NewDummySession(dID string) Session{
	connList:=make(map[string]Connection)
	// create dummy nodes // TODO will derive peers for session using blockchain
	sNode1 := node.Node{GID: "snode1", RemoteIP: "localhost", LocalIP: "localhost"}
	sNode2 := node.Node{GID: "snode2", RemoteIP: "localhost", LocalIP: "localhost"}
	vNode := node.Node{GID:"vnode", RemoteIP:"localhost", LocalIP:"localhost"}
	// create dummy connections
	sNodeConn1:= Connection{Mutex: sync.Mutex{}, Node: sNode1, Role: SERVICER}
	sNodeConn2:= Connection{Mutex: sync.Mutex{}, Node: sNode2, Role: SERVICER}
	vNodeConn:= Connection{Mutex:sync.Mutex{}, Node: vNode, Role: VALIDATOR}
	// add to List
	connList[sNode1.GID]=sNodeConn1
	connList[sNode2.GID]=sNodeConn2
	connList[vNode.GID]=vNodeConn
	// return resulting session
	return Session{dID,connList, sync.Mutex{}}
}

/***********************************************************************************************************************
Session Methods
 */
func (session *Session) GetConnections() map[string]Connection {
	var once sync.Once
	once.Do(func() {
		session.ConnList = make(map[string]Connection)
	})
	return session.ConnList
}

func (session *Session) AddConnection(connection Connection) {
	fmt.Println("REGISTERING CONN "+ connection.Conn.RemoteAddr().String())
	session.Lock()
	defer session.Unlock()
	session.GetConnections()[connection.Conn.RemoteAddr().String()] = connection // added by remote addr
}

func (session *Session) GetConnection(gid string) Connection {
	session.Lock()
	defer session.Unlock()
	sessConnList:= session.GetConnections()
	return sessConnList[gid]
}

func (session *Session) ClearConnections(){
	session.Lock()
	defer session.Unlock()
	session.ConnList = make(map[string]Connection)
}

func (session *Session) Listen(port string, host string) {
	l, err := net.Listen(_const.SESSION_CONN_TYPE, host+":"+port)				// listen on port & host
	if err != nil {																// handle server creation error
		logs.NewLog("Unable to create a new "+_const.SESSION_CONN_TYPE+" server on port:"+port, logs.PanicLevel, logs.JSONLogFormat)
		logs.NewLog("ERROR: "+err.Error(), logs.PanicLevel, logs.JSONLogFormat)
	}
	defer l.Close()																// close the server after serve and listen finishes
	logs.NewLog("Listening on port :"+port, logs.InfoLevel, logs.JSONLogFormat) // log the new connection
	for {																		// for the duration of incoming requests
		conn, err := l.Accept()													// accept the connection
		if err != nil {															// handle request accept err
			logs.NewLog("Unable to accept the "+_const.SESSION_CONN_TYPE+" Conn on port:"+port, logs.PanicLevel, logs.JSONLogFormat)
			logs.NewLog("ERROR: "+err.Error(), logs.PanicLevel, logs.JSONLogFormat)
		}
		connection:= NewConnection(conn)   // create a new connection from connection
		session.AddConnection(*connection) // TODO consider returning the connection and register it from caller
	}
}

func (session *Session) Dial(port string, host string, connection Connection) {
	conn, err := net.Dial(_const.SESSION_CONN_TYPE, host+":"+port) 				// establish a connection
	if err != nil { 															// handle connection error
		logs.NewLog("Unable to establish "+_const.SESSION_CONN_TYPE+" connection on port "+host+":"+port,
			logs.PanicLevel, logs.JSONLogFormat)
	}
	connection.Conn = conn            // save the connection to this connection instance
	go connection.Receive()           // run receive to listen for incoming messages
	session.AddConnection(connection) // TODO consider returning the connection and register it from caller
}
