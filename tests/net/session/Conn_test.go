package session

import (
	"fmt"
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/net/session"
	"testing"
)

func TestConnection(t *testing.T) {
	const port = "3333"
	const host = "localhost"
	go session.Listen(port,host)                                   // start the server
	peer1 := session.NewPeer()
	peer1.CreateConnection(port,host)
	t.Log(peer1.Conn.RemoteAddr())
	fmt.Println(peer1.Conn.LocalAddr().String())
	// get peer2 from the peerlist
	for session.GetSessionPeerlist()[peer1.Conn.LocalAddr().String()]==(session.Peer{}){

	}
	peer2 := session.GetSessionPeerlist()[peer1.Conn.LocalAddr().String()]
	// Create a new message to send
	pl := message.Payload{0,"I am peer 2 : "+peer2.Conn.LocalAddr().String()}
	m:= message.NewMessage(pl)
	peer2.Send(m)
	peer1.Receive()
	peer1.CloseConnection()
}
