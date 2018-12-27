// This package is network code relating to pocket 'sessions'
package sessio

import (
	"github.com/pokt-network/pocket-core/logs"
	"sync"
)

// "session.go" specifies the session structure, methods, and functions

/*
This is the session structure.
*/
type Session struct {
	DevID      string                `json:"devid"`    	// "DevID" is the developer's ID that identifies the sessio
	ConnList   map[string]Connection `json:"connList"` 	// "ConnList" is the List of peer connections [GID] Connection
	sync.Mutex `json:"mutex"`
}

/***********************************************************************************************************************
Session Constructor
*/

/*
"NewEmptySession" returns an empty session object with the devID prefilled
 */
func NewEmptySession(dID string) Session {
	return Session{DevID: dID}													// prefill the devID and return
}

/***********************************************************************************************************************
Session Methods
*/

/*
"GetConnections" returns a map of Connection objects [GID]Connection
 */
func (session *Session) GetConnections() map[string]Connection {
	var once sync.Once
	once.Do(func() {															// only do once
		if session.ConnList == nil {											// if nil connectionList
			session.ConnList = make(map[string]Connection)						// make a new map
		}
	})
	return session.ConnList														// return the connectionlist
}
/*
"AddConnection" adds a connection object to the session
 */
func (session *Session) AddConnection(connection Connection) {
	logs.NewLog("Adding Connection: " + connection.Peer.(SessionPeer).GID + " to Session: " +
		session.DevID, logs.InfoLevel, logs.JSONLogFormat)
	session.Lock()																// lock the session
	defer session.Unlock()														// after function completes unlock
	session.GetConnections()[connection.Peer.(SessionPeer).GID] = connection	// add connection to list
}

/*
"RemoveConnection" removes a connection object from the session
 */
func (session *Session) RemoveConnection(connection Connection) {
	session.Lock()         														// lock the session
	defer session.Unlock() 														// after the function completes unlock
	logs.NewLog("Removed peer: "+connection.Peer.(SessionPeer).GID, logs.InfoLevel, logs.JSONLogFormat)
	delete(session.ConnList, connection.Peer.(SessionPeer).GID) 				// delete the item from the map
}

/*
"GetConnection" returns the connection from the session by peer.GID
 */
func (session *Session) GetConnection(gid string) Connection {
	session.Lock()																// lock the session
	defer session.Unlock()														// after function completes unlock
	return session.GetConnections()[gid]										// get and return
}

/*
"NewConnections" adds the connections to the session from []SessionPeer and dials them
 */
// TODO dial each individual sessionPeer -> add connection to connList
// TODO if fail then donot add connections
// TODO allow flexible ports: startup flag -> in Node structure -> in SessionPeer structure -> called here
func (session *Session) NewConnections(sp []SessionPeer) {
	for _, sessionPeer := range sp {
		//connection.Dial("3333", sessionPeer.LocalIP)
		connection := Connection{Peer: sessionPeer}
		session.AddConnection(connection)
	}
}

/*
"ClearConnections" removes all connections from a session
 */
func (session *Session) ClearConnections() {
	session.Lock()																// lock the session
	defer session.Unlock()														// after function completes unlock
	for _, conn := range session.ConnList {										// for each connection
		conn.CloseConnection()													// close the conn
	}
	session.ConnList = make(map[string]Connection)								// clear the list
}
