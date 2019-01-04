package session

import (
	"testing"
)

func TestSessionMessage(t *testing.T) { // deprecated
	//const PORT = "3333"      // port for listener
	//const HOST = "localhost" // host for listener
	//const DEVID = "dummy-id"
	//// STEP 1: CREATE DUMMY SESSION PEERS
	//var spl []session.SessionPeer
	//sNode1 := session.SessionPeer{Role: session.SERVICER, Node: node.Node{GID: "sNode1", RemoteIP: "localhost", LocalIP: "localhost"}}
	//sNode2 := session.SessionPeer{Role: session.SERVICER, Node: node.Node{GID: "sNode2", RemoteIP: "localhost", LocalIP: "localhost"}}
	//vNode := session.SessionPeer{Role: session.VALIDATOR, Node: node.Node{GID: "vNode", RemoteIP: "localhost", LocalIP: "localhost"}}
	//dNode := session.SessionPeer{Role: session.DISPATCHER, Node: node.Node{GID: "dNode", RemoteIP: "localhost", LocalIP: "localhost"}}
	//spl = append(spl, sNode1, sNode2, vNode, dNode) // add them to a list
	//// STEP 2: CREATE NEW SESSION MESSAGE
	//nSPL := session.NewSessionPayload{DevID: DEVID, Peers: spl}
	//message := session.NewSessionMessage(nSPL)
	//// STEP 3: SERVE ON PORT
	//server := session.Connection{}
	//go server.Listen(PORT, HOST)
	//// STEP 4: DIAL TO PORT
	//client := session.Connection{}
	//client.Dial(PORT, HOST)
	//// STEP 5: send message over the wire
	//time.Sleep(time.Second)
	//server.Send(message, session.NewSessionPayload{})
	//client.Send(message, session.NewSessionPayload{})
	//time.Sleep(time.Second)
	//// STEP 6: confirm added to session list
	//sessionList := session.GetSessionList().List
	//if len(sessionList) == 0 {
	//	t.Fatalf("Empty Session List")
	//}
	//if _, contains := sessionList[DEVID]; !contains {
	//	t.Fatal("Session not within list")
	//}
	//// STEP 7: confirm nodes are within peerlist
	//peerList := peers.GetPeerList()
	//if peers.GetPeerCount() == 0 {
	//	t.Fatalf("Empty Peer List")
	//}
	//if !peerList.Contains(sNode1.GID) {
	//	t.Fatalf("Peer: " + sNode1.GID + " is not within the peerList")
	//}
	//if !peerList.Contains(sNode2.GID) {
	//	t.Fatalf("Peer: " + sNode2.GID + " is not within the peerList")
	//}
	//if !peerList.Contains(vNode.GID) {
	//	t.Fatalf("Peer: " + vNode.GID + " is not within the peerList")
	//}
	//if !peerList.Contains(dNode.GID) {
	//	t.Fatalf("Peer: " + dNode.GID + " is not within the peerList")
	//}
	//// STEP 8: confirm that session contains the session peers
	//session := sessionList[DEVID]
	//if len(session.GetPeers()) == 0 {
	//	t.Fatalf("There are no peers within the session")
	//}
	//if session.GetPeer(sNode1.GID) == (session.Connection{}) {
	//	t.Fatalf("Peer: " + sNode1.GID + " is not within the sessionList")
	//}
	//if session.GetPeer(sNode2.GID) == (session.Connection{}) {
	//	t.Fatalf("Peer: " + sNode2.GID + " is not within the sessionList")
	//}
	//if session.GetPeer(vNode.GID) == (session.Connection{}) {
	//	t.Fatalf("Peer: " + vNode.GID + " is not within the sessionList")
	//}
	//if session.GetPeer(dNode.GID) == (session.Connection{}) {
	//	t.Fatalf("Peer: " + dNode.GID + " is not within the sessionList")
	//}
}
