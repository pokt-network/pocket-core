// This package deals with all things networking related.
package session

import (
	"net"
	"sync"
)
// TODO the name peer clashes too much, consider new name for persistent connection instance
// The peer structure represents a persistent connection between two nodes within a session
type Peer struct {
	Conn net.Conn		// the persistent connection between the two
	sync.Mutex			// the lock for sending messages and closing the connection
}

/***********************************************************************************************************************
Everything below is temporary. This session peerlist is a basic structure to register a new peer connection easily for
testing. The currently developed solution is decoupled from the sessionList under the global session package

TODO need to integrate the SessionPeerlist with the sessionList
 */
var peerList map[string]Peer
var once sync.Once

func GetSessionPeerlist() map[string]Peer {
	once.Do(func() {
		peerList = make(map[string]Peer)
	})
	return peerList
}

func RegisterSessionPeerConnection(peer Peer) {
	GetSessionPeerlist()[peer.Conn.RemoteAddr().String()] = peer // added by remote addr
}

/**********************************************************************************************************************/
