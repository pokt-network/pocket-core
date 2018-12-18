package session

import (
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/net/session"
	"testing"
	"time"
)

func TestConnection(t *testing.T) {
	// establish a new listener
	const port = "3333"
	const host = "localhost"
	go session.Listen(port,host)
	// create a new peer
	peer1 := session.NewPeer()
	peer1.CreateConnection(port,host)
	// get peer2 from the peerlist
	for session.GetSessionPeerlist()[peer1.Conn.LocalAddr().String()]==(session.Peer{}){} // wait for peer registration
	// get peer2 from global peer list
	peer2 := session.GetSessionPeerlist()[peer1.Conn.LocalAddr().String()]
	// create a new message to send
	peer1Payload := message.Payload{0,"I am  peer 1: "+peer1.Conn.LocalAddr().String()}
	peer2Payload := message.Payload{0,"I am  peer 2: "+peer2.Conn.LocalAddr().String()}
	// send several messages back and forth
	peer1.Send(message.NewMessage(peer1Payload))
	peer2.Send(message.NewMessage(peer2Payload))
	peer1.Send(message.NewMessage(peer1Payload))
	peer2.Send(message.NewMessage(peer2Payload))
	// don't close the connection without waiting for a second to allow for threads to complete their process
	time.Sleep(1 * time.Second)
	// close the connections
	peer1.CloseConnection()
	peer2.CloseConnection()
}
