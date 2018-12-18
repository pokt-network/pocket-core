// This package deals with all things networking related.
package session

import (
	"net"
	"sync"
)

type Peer struct {
	Conn net.Conn
	sync.Mutex
}

var peerList map[string]Peer
var once sync.Once

func GetSessionPeerlist() map[string]Peer{
	// TODO add concurrency protection
	once.Do(func(){
		peerList = make(map[string]Peer)
	})
	return peerList
}

func AddPeerToSessionPeersList(peer Peer){
	GetSessionPeerlist()[peer.Conn.RemoteAddr().String()]=peer // added by remote addr
}
