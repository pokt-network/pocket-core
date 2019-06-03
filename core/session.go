package core

import (
	"crypto"
)

var HashingAlgorithm = crypto.SHA256

type Session struct {
	Key       []byte       `json:"key"`
	DevID     []byte       `json:"devid"`
	BlockHash []byte       `json:"blockhash"`
	Chain     []byte       `json:"chain"`
	Capacity  int          `json:"capacity"`
	Nodes     SessionNodes `json:"node"`
}

type SessionNodes []Node

// "NewSession" creates a new session from seed data.
func NewSession(s SessionSeed) (*Session, error) {
	err := s.ErrorCheck()
	if err != nil {
		return nil, err
	}
	session := &Session{BlockHash: s.BlockHash, DevID: s.DevID, Chain: s.RequestedChain, Capacity: s.Capacity}
	err = session.GenKey()
	if err != nil {
		return nil, err
	}
	err = session.GenNodes(s.NodeList)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// TODO again optimizations can be made depending on how we store Node World State
// for now this is slow but passable
func (sn *SessionNodes) Contains(gid string) bool {
	for _, node := range *sn {
		if node.GID == gid {
			return true
		}
	}
	return false
}

func (s *Session) ValidityCheck(mygid string) error {
	if !s.Nodes.Contains(mygid) {
		return InvalidSessionError
	}
	return nil
}

// "Hash" returns the session hashing algorithm result of input
func (s *Session) Hash(input []byte) []byte {
	hasher := HashingAlgorithm.New()
	hasher.Write(input)
	return hasher.Sum(nil)
}

// "GenKey" generates the session key = SessionHashingAlgo(devid+chain+blockhash)
func (s *Session) GenKey() error {
	if s.BlockHash == nil || s.DevID == nil || s.Chain == nil {
		return IncompleteSessionError
	}
	var seed = s.BlockHash
	seed = append(seed, s.DevID...)
	seed = append(seed, s.Chain...)
	s.Key = s.Hash(seed)
	return nil
}

// "GenNodes" generates the nodes of the session
func (s *Session) GenNodes(pool NodePool) error {
	n, err := pool.GetNodes(*s)
	if err != nil {
		return err
	}
	s.Nodes = n
	return nil
}
