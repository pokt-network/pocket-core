package session

import (
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/net/session"
	"testing"
)

func TestConnection(t *testing.T) {
	const port = "3333"
	const host = "localhost"
	go session.Listen(port,host)
	peer1 := session.NewPeer()
	peer1.CreateConnection(port,host)
	// get peer2 from the peerlist
	for session.GetSessionPeerlist()[peer1.Conn.LocalAddr().String()]==(session.Peer{}){

	}
	peer2 := session.GetSessionPeerlist()[peer1.Conn.LocalAddr().String()]
	// Have both listen to demonstrate sending and receiving ability
	go peer1.Receive() // test if you can receive and send
	go peer2.Receive() // test if you can receive and send
	// Create a new message to send
	peer1Payload := message.Payload{0,"I am  peer 1: "+peer1.Conn.LocalAddr().String()}
	peer2Payload := message.Payload{0,"I am  peer 2: "+peer2.Conn.LocalAddr().String()}
	peer1.Send(message.NewMessage(peer1Payload))
	peer2.Send(message.NewMessage(peer2Payload))
	peer1.CloseConnection()
}
