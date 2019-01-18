// This package is network code relating to pocket 'sessions'
package session

import (
	"sync"

	"github.com/pokt-network/pocket-core/logs"
)

type Session struct {
	DevID      string   `json:"devid"`
	Peers      PeerList `json:"sessionpeerlist"`
	sync.Mutex `json:"mutex"`
}

var o sync.Once

// "NewSession" returns an empty session object with the devID prefilled
func NewSession(dID string) Session {
	return Session{DevID: dID}
}

// "GetPeers" returns a map of Connection objects [GID]Connection
func (session *Session) GetPeers() map[string]SessionPeer {
	o.Do(func() {
		session.Peers.List = make(map[string]SessionPeer)
	})
	return session.Peers.List
}

// "Add" adds a connection object to the session
func (session *Session) AddPeer(sPeer SessionPeer) {
	logs.NewLog("Adding Connection: "+sPeer.GID+" to Session: "+
		session.DevID, logs.InfoLevel, logs.JSONLogFormat)
	session.Lock()
	defer session.Unlock()
	session.GetPeers()[sPeer.GID] = sPeer
}

// "Remove" removes a connection object from the session
func (session *Session) RemovePeer(sPeer SessionPeer) {
	session.Lock()
	defer session.Unlock()
	delete(session.Peers.List, sPeer.GID)
	logs.NewLog("Removed peer: "+sPeer.GID, logs.InfoLevel, logs.JSONLogFormat)
}

// "GetPeer" returns the connection from the session by peer.GID
func (session *Session) GetPeer(gid string) SessionPeer {
	session.Lock()
	defer session.Unlock()
	return session.GetPeers()[gid]
}

// "NewPeers" adds the connections to the session from []SessionPeer and dials them
func (session *Session) NewPeers(sp []SessionPeer) {
	for _, sessionPeer := range sp {
		session.AddPeer(sessionPeer)
	}
}

// "ClearPeers" removes all connections from a session
func (session *Session) ClearPeers() {
	session.Lock()
	defer session.Unlock()
	session.Peers.List = make(map[string]SessionPeer)
}
