package tests

import (
	"testing"

	"github.com/pokt-network/pocket-core/session"
)

func dummySession() session.Session {
	return session.Session{DevID: "test", PL: session.NewPeerList()}
}

func TestSession(t *testing.T) {
	s := dummySession()
	s.AddPeer(session.Peer{Role: session.SERVICER, Node: dummyNode()})
	if len(s.PL.M) != 1 {
		t.Fatalf("AddPeer(Peer) did not result in a peercount of 1")
	}
	s.RemovePeer(dummyNode().GID)
	if len(s.PL.M) != 0 {
		t.Fatalf("The peerlist size does not equal 0 after the RemovePeer()")
	}
}

func TestSessionList(t *testing.T) {
	sl := session.SessionList()
	sl.Add(dummySession())
	if sl.Count() != len(sl.M) {
		t.Fatalf("SessionList.Count() != len(SessionList)")
	}
	if sl.Count() != 1 {
		t.Fatalf("After SessionList.Add() the result did not = 1 session")
	}
	if !sl.Contains(dummySession().DevID) {
		t.Fatalf("SessionList.Contains() returned false when it is expected that the session is in the list")
	}
	sl.Remove(dummySession())
	if sl.Count() != 0 {
		t.Fatalf("After calling SessionList.Remove() the result did not = 0 sessions")
	}
}
