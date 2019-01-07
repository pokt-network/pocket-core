// This package is network code relating to pocket 'sessions'
package session

import (
	"github.com/pokt-network/pocket-core/logs"
	"sync"
)

// "session.go" specifies the session structure, methods, and functions

/*
This is the session structure.
*/
type Session struct {
	DevID      string               `json:"devid"`
	Peers      SessionPeerList		`json:"sessionpeerlist"`
	sync.Mutex 						`json:"mutex"`
}

/***********************************************************************************************************************
Session Constructor
*/

/*
"NewEmptySession" returns an empty session object with the devID prefilled
 */
func NewEmptySession(dID string) Session {
	return Session{DevID: dID}											// prefill the devID and return
}

/***********************************************************************************************************************
Session Methods
*/

/*
"GetPeers" returns a map of Connection objects [GID]Connection
 */
func (session *Session) GetPeers() map[string]SessionPeer {
	var once sync.Once
	once.Do(func() {													// only do once
		if session.Peers.List == nil { 									// if nil connectionList
			session.Peers.List = make(map[string]SessionPeer) 			// make a new map
		}
	})
	return session.Peers.List 											// return the connectionlist
}
/*
"AddPeer" adds a connection object to the session
 */
func (session *Session) AddPeer(sPeer SessionPeer) {
	logs.NewLog("Adding Connection: " + sPeer.GID + " to Session: " +
		session.DevID, logs.InfoLevel, logs.JSONLogFormat)
	session.Lock()                        								// lock the session
	defer session.Unlock()                								// after function completes unlock
	session.GetPeers()[sPeer.GID] = sPeer 								// add sPeer to list
}

/*
"RemovePeer" removes a connection object from the session
 */
func (session *Session) RemovePeer(sPeer SessionPeer) {
	session.Lock()                        								// lock the session
	defer session.Unlock()                								// after the function completes unlock
	delete(session.Peers.List, sPeer.GID) 								// delete the item from the map
	logs.NewLog("Removed peer: "+sPeer.GID, logs.InfoLevel, logs.JSONLogFormat)
}

/*
"GetPeer" returns the connection from the session by peer.GID
 */
func (session *Session) GetPeer(gid string) SessionPeer {
	session.Lock()                 										// lock the session
	defer session.Unlock()         										// after function completes unlock
	return session.GetPeers()[gid] 										// get and return
}

/*
"NewPeers" adds the connections to the session from []SessionPeer and dials them
 */
func (session *Session) NewPeers(sp []SessionPeer) {
	for _, sessionPeer := range sp {
		session.AddPeer(sessionPeer)
	}
}

/*
"ClearPeers" removes all connections from a session
 */
func (session *Session) ClearPeers() {
	session.Lock()                                    					// lock the session
	defer session.Unlock()                            					// after function completes unlock
	session.Peers.List = make(map[string]SessionPeer) 					// clear the list
}
