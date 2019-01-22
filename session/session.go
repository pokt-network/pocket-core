// This package is network code relating to pocket 'sessions'
package session

import (
	"sync"

	"github.com/pokt-network/pocket-core/types"
)

type Session struct {
	DevID string   `json:"devid"`
	Peers PeerList `json:"peerlist"`
}

var one sync.Once

// "NewSession" returns an empty session object with the devID prefilled
func NewSession(dID string) Session {
	return Session{DevID: dID}
}

// "GetPeers" returns a map of Connection objects [GID]Connection
func (s *Session) GetPeers() PeerList {
	one.Do(func() {
		s.Peers = NewPeerList()
	})
	return s.Peers
}

// "Add" adds a peer to the session
func (s *Session) AddPeer(sPeer Peer) {
	(*types.List)(&s.Peers).Add(sPeer.GID, sPeer)
}

func (s *Session) AddPeers(peers []Peer) {
	for _, v := range peers {
		s.AddPeer(v)
	}
}

// "Remove" removes a connection object from the session
func (s *Session) RemovePeer(gid string) {
	(*types.List)(&s.Peers).Remove(gid)
}

// "GetPeer" returns the connection from the session by peer.GID
func (s *Session) GetPeer(gid string) Peer {
	return (*types.List)(&s.Peers).Get(gid).(Peer)
}
