// This package is network code relating to pocket 'sessions'
package session

import (
	"sync"

	"github.com/pokt-network/pocket-core/types"
)

type Session struct {
	DevID string   `json:"devid"`
	PL    PeerList `json:"peerlist"`
}

var one sync.Once

// "NewSession" returns an empty session object with the devID prefilled
func NewSession(dID string) Session {
	return Session{DevID: dID}
}

// "Peers" returns a map of sessionPeers [GID]Connection
func (s *Session) Peers() PeerList {
	one.Do(func() {
		s.PL = NewPeerList()
	})
	return s.PL
}

// "Add" adds a peer to the session
func (s *Session) AddPeer(sPeer Peer) {
	(*types.List)(&s.PL).Add(sPeer.GID, sPeer)
}

func (s *Session) AddPeers(peers []Peer) {
	for _, v := range peers {
		s.AddPeer(v)
	}
}

// "Remove" removes a sessionPeer from the session
func (s *Session) RemovePeer(gid string) {
	(*types.List)(&s.PL).Remove(gid)
}

// "GetPeer" returns the sessionPeer from the session by peer.GID
func (s *Session) GetPeer(gid string) Peer {
	return (*types.List)(&s.PL).Get(gid).(Peer)
}
