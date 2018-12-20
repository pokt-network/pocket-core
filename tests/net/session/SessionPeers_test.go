package session

import (
	"encoding/json"
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/net/session"
	"testing"
)

func TestSessionPeers(t *testing.T) { //TODO can't send spl over the wire cause it has conn object. instead send array of ip's and their roles
	const host = "localhost"
	const port = "3333"
	// create a listener on localhost:lport
	go session.Listen(port,host)
	// create a peer to be the sender
	sender := session.NewPeer()
	// create a connection to localhost:lport
	sender.CreateConnection(port,host)
	// wait for registration
	for session.GetSessionPeerlist()[sender.Conn.LocalAddr().String()]==(session.Peer{}){}
	// get receiver from global peer list
	receiver := session.GetSessionPeerlist()[sender.Conn.LocalAddr().String()]
	if len(session.GetSessionPeerlist())==0 {
		t.Fatalf("Empty sessionPeerList")
	}
	t.Log(session.GetSessionPeerlist()[sender.Conn.LocalAddr().String()].Conn.RemoteAddr().String())
	// marshal sessionPeerList into JSON
	splJSON, err := json.Marshal(session.GetSessionPeerlist())
	if err != nil {
		t.Fatalf("Unable to marshal sessionList into JSON "+ err.Error())
	}
	// generate the SessionPeerMessage to send
	t.Log(string(splJSON))
	m :=message.NewSessionPeersMessage(splJSON)
	// clear session list
	session.ClearSessionPeerList()
	t.Log("cleared spl: ")
	t.Log(session.GetSessionPeerlist())
	t.Log("Sending spl over wire")
	// send spl over the wire
	sender.Send(m)
	//for len(session.GetSessionPeerlist())==0 {}
	//t.Log(session.GetSessionPeerlist())
	receiver.CloseConnection()
}
