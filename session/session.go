package session

import (
	"crypto"
)

var HashingAlgorithm = crypto.SHA256

type Session struct {
	Key       []byte       `json:"key"`
	DevID     []byte       `json:"devid"`
	BlockHash []byte       `json:"blockhash"`
	Nodes     SessionNodes `json:"node"`
	Chain     []byte       `json:"chain"`
}

type SessionNodes struct {
	ServiceNodes    []Node
	ValidatorNodes  []Node
	DelegatedMinter Node
}

// "NewSession" creates a new session from seed data.
func NewSession(s Seed) (*Session, error) {
	err := s.ErrorCheck()
	if err != nil {
		return nil, err
	}
	session := &Session{BlockHash: s.BlockHash, DevID: s.DevID, Chain: s.RequestedChain}
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

// "Hash" returns the session hashing algorithm result of input
func (s *Session) Hash(input []byte) []byte {
	hasher := HashingAlgorithm.New()
	hasher.Write(input)
	return hasher.Sum(nil)
}

// "GenKey" generates the session key = SessionHashingAlgo(devid+chain+blockhash)
func (s *Session) GenKey() error {
	if s.BlockHash == nil || s.DevID == nil || s.Chain == nil {
		return IncompleteSession
	}
	var seed = s.BlockHash
	seed = append(seed, s.DevID...)
	seed = append(seed, s.Chain...)
	s.Key = s.Hash(seed)
	return nil
}

// "GenNodes" generates the nodes of the session
func (s *Session) GenNodes(pool NodePool) error {
	n, err := pool.GetSessionNodes(*s)
	if err != nil {
		return err
	}
	s.Nodes = n
	return nil
}
