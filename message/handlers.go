package message

import (
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/session"
	"fmt"
	"net"
	"sync"
)

func HandleMessage(message *Message, addr *net.UDPAddr){
	switch message.Payload.ID {
	case 1: // sessionPeers message
		newSessionMessageHandler(message)
	case 2: //enter network
		enterNetworkMessage(message)
	case 3: //exit network
		exitNetworkMessage(message)
	}
}

func newSessionMessageHandler(message *Message) {
	sList := session.GetSessionList()
	nSPL := message.Payload.Data.(NewSessionPayload)				// extract the NewSessionPayload
	s := session.Session{DevID: nSPL.DevID, 						// create a session using developerID from payload
		Peers: session.SessionPeerList{List: make(map[string]session.SessionPeer)},
		Mutex: sync.Mutex{}}
	s.NewPeers(nSPL.Peers)                  						// create new connections with each peer
	sList.AddSession(s)                     						// register the session
	session.AddSessPeersToPL(nSPL.Peers) 							// add peers to peerList
}

//NOTE: this is for pocket core mvp centralized dispatcher
// may remove for production
func enterNetworkMessage(message *Message){
	// get node
	n := message.Payload.Data.(EnterNetworkPayload)
	// Cross reference white list
	if !node.GetDispatchWhitelist().Contains(n.GID){
		fmt.Println("Node: ", n.GID, " rejected because it is not within service_whitelist.json")
		return
	}
	// add to peerlist
	node.GetPeerList().AddPeer(n.Node)
	// add to dispatch peers
	node.NewDispatchPeer(n.Node)
}

//NOTE: this is for pocket core mvp centralized dispatcher
// may remove for production
func exitNetworkMessage(message *Message){
	// get node
	n := message.Payload.Data.(ExitNetworkPayload)
	// add to peerlist
	node.GetPeerList().RemovePeer(n.Node)
	// add to dispatch peers
	node.DeleteDispatchPeer(n.Node)
}
