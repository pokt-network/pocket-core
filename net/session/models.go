// This package deals with all things networking related.
package session

import (
	"fmt"
	"net"
	"sync"
)
// The peer structure represents a persistent connection between two nodes within a session
type Peer struct {
	Conn net.Conn		// the persistent connection between the two
	sync.Mutex			// the lock for sending messages and closing the connection
}

/***********************************************************************************************************************
Everything below is temporary. This session peerlist is a basic structure to register a new peer connection easily for
testing. The currently developed solution is decoupled from the sessionList under the global session package
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
	fmt.Println("REGISTERING CONN "+peer.Conn.RemoteAddr().String())
	GetSessionPeerlist()[peer.Conn.RemoteAddr().String()] = peer // added by remote addr
}

func ClearSessionPeerList(){
	peerList = make(map[string]Peer)
}

/**********************************************************************************************************************/
