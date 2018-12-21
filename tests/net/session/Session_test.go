package session

import (
	"github.com/pokt-network/pocket-core/net/sessio"
	"github.com/pokt-network/pocket-core/node"
	"testing"
	"time"
)

func TestSessionMessage(t *testing.T) {
	const LPORT = "3333"      // port for listener
	const LHOST = "localhost" // host for listener
	const SHOST = "localhost" // host for sender
	// STEP 1: CREATE DUMMY SESSION PEERS
	var spl []sessio.SessionPeer
	sNode1 := sessio.SessionPeer{Role: sessio.SERVICER, Node: node.Node{GID: "sNode1", RemoteIP: "localhost", LocalIP: "localhost"}}
	sNode2 := sessio.SessionPeer{Role: sessio.SERVICER, Node: node.Node{GID: "sNode2", RemoteIP: "localhost", LocalIP: "localhost"}}
	vNode := sessio.SessionPeer{Role: sessio.VALIDATOR, Node: node.Node{GID: "vNode", RemoteIP: "localhost", LocalIP: "localhost"}}
	spl = append(spl, sNode1,sNode2, vNode)
	// STEP 2: CREATE NEW SESSION MESSAGE
	nSPL := sessio.NewSessionPayload{DevID: "dummy-id", Peers: spl}
	message:= sessio.NewSessionMessage(nSPL)
	// STEP 3: CREATE DUMMY SESSION
	dSess:= sessio.Session{DevID:"senderSession"}
	// STEP 3: LISTEN FOR INCOMING MESSAGE
	sess := sessio.Session{DevID:"receiverSession"}
	go sess.Listen(LPORT,LHOST)
	// STEP 4: ESTABLISH CONN WITH SERVER
	dSess.Dial(LPORT,LHOST,sessio.Connection{SessionPeer: sessio.SessionPeer{sessio.DISPATCHER, node.Node{GID: "dNode", RemoteIP: "localhost", LocalIP: "localhost"}}})
	// STEP 5: SEND MESSAGE OVER THE WIRE
	for len(dSess.ConnList) ==0 {}
	dConn:=dSess.ConnList["dNode"]
	dConn.Send(message)
	time.Sleep(time.Second*2)
}
