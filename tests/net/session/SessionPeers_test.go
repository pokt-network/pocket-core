package session

import (
	"encoding/json"
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/net/sessio"
	"testing"
)

func TestSessionPeers(t *testing.T) { //TODO can't send spl over the wire cause it has conn object. instead send array of ip's and their roles
	const host = "localhost" // broken test
	const port = "3333"
	// create a listener on localhost:lport
	go sessio.Listen(port,host)
	// create a peer to be the sender
	sender := sessio.NewConnection()
	// create a connection to localhost:lport
	sender.CreateConnection(port,host)
	// wait for registration
	for sessio.GetSessionConnList()[sender.Conn.LocalAddr().String()].GID==""{}
	// get receiver from global peer list
	receiver := sessio.GetSessionConnList()[sender.Conn.LocalAddr().String()]
	if len(sessio.GetSessionConnList())==0 {
		t.Fatalf("Empty sessionPeerList")
	}
	t.Log(sessio.GetSessionConnList()[sender.Conn.LocalAddr().String()].Conn.RemoteAddr().String())
	// marshal sessionPeerList into JSON
	splJSON, err := json.Marshal(sessio.GetSessionConnList())
	if err != nil {
		t.Fatalf("Unable to marshal sessionList into JSON "+ err.Error())
	}
	// generate the SessionPeerMessage to send
	t.Log(string(splJSON))
	m :=message.NewSessionPeersMessage(splJSON)
	// clear session list
	sessio.ClearSessionConnList()
	t.Log("cleared spl: ")
	t.Log(sessio.GetSessionConnList())
	t.Log("Sending spl over wire")
	// send spl over the wire
	sender.Send(m)
	//for len(session.GetSessionConnList())==0 {}
	//t.Log(session.GetSessionConnList())
	receiver.CloseConnection()
}
